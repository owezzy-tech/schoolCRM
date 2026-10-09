import logging
from collections.abc import AsyncGenerator, AsyncIterable
from contextlib import aclosing
from typing import Annotated
from uuid import UUID

import psycopg
from fastapi import APIRouter, Depends, Query, Request
from fastapi.sse import EventSourceResponse, ServerSentEvent
from pydantic import Field, model_validator

from adapters.controllers.curriculum_controller import JSONAPIResponse, ScopeDTO
from domain.entities.curriculum import CurriculumError
from domain.entities.lesson_generation import GenerationContext, LessonGenerationRequest
from infrastructure.auth import AuthContext, get_auth_context
from infrastructure.lesson_workflow import LessonWorkflow

logger = logging.getLogger(__name__)
router = APIRouter(
    prefix="/v1/rag/lessons", tags=["lesson-generation"], default_response_class=JSONAPIResponse
)


class GenerateRequest(ScopeDTO):
    request_id: UUID
    topic: str = Field(min_length=1, max_length=1000)
    duration_minutes: int = Field(ge=1, le=180, strict=True)

    @model_validator(mode="after")
    def validate_identity(self) -> "GenerateRequest":
        if self.request_id.int == 0:
            raise ValueError("request_id must be a nonzero UUID")
        return self

    def generation(self) -> LessonGenerationRequest:
        return LessonGenerationRequest(
            self.request_id, self.scope(), self.topic, self.duration_minutes
        )


def get_workflow(request: Request) -> LessonWorkflow:
    workflow = request.app.state.container.lesson_workflow
    if workflow is None:
        raise CurriculumError(503, "Configure curriculum storage before lesson generation")
    return workflow


Actor = Annotated[AuthContext, Depends(get_auth_context)]
Workflow = Annotated[LessonWorkflow, Depends(get_workflow)]


def context(actor: AuthContext) -> GenerationContext:
    try:
        actor_id = UUID(actor.subject)
    except ValueError:
        raise CurriculumError(401, "A valid authenticated Go user is required") from None
    if not actor.token or actor_id.int == 0:
        raise CurriculumError(401, "An authenticated Go user is required")
    return GenerationContext(actor_id, actor.token)


async def authorized_events(
    request: GenerateRequest, actor: Actor, workflow: Workflow
) -> AsyncGenerator[tuple[str, dict], None]:
    # A dependency completes before streaming starts, so authentication and
    # authorization failures keep their ordinary JSON:API status and body.
    return await workflow.stream(request.generation(), context(actor))


def thread_resource(request_id: str, thread: dict) -> dict:
    return {"type": "lesson-thread", "id": request_id, "attributes": thread}


@router.post("/generate")
async def generate(request: GenerateRequest, actor: Actor, workflow: Workflow) -> dict:
    result = await workflow.execute(request.generation(), context(actor))
    return {
        "jsonapi": {"version": "1.1"},
        "data": {"type": "lesson-generation", "id": str(request.request_id), "attributes": result},
    }


@router.post("/generate/stream", response_class=EventSourceResponse)
async def generate_stream(
    request: GenerateRequest,
    events: Annotated[AsyncGenerator[tuple[str, dict], None], Depends(authorized_events)],
) -> AsyncIterable[ServerSentEvent]:
    def event(name: str, **fields) -> ServerSentEvent:
        data = {"protocolVersion": 1, "requestID": str(request.request_id), **fields}
        return ServerSentEvent(event=name, data=data)

    async with aclosing(events):
        try:
            async for name, fields in events:
                yield event(name, **fields)
                if name == "structured-result":
                    result = fields["result"]
                    yield event(
                        "approval-request",
                        planID=result["planID"],
                        version=result["version"],
                        action="submit-for-review",
                    )
        except CurriculumError as error:
            yield event("terminal-error", status=error.status_code, detail=error.detail)
            return
        except psycopg.Error:
            yield event("terminal-error", status=503, detail="Curriculum storage unavailable")
            return
        except Exception as error:
            logger.error(
                "Lesson generation stream failed (%s, request %s)",
                type(error).__name__,
                request.request_id,
            )
            yield event(
                "terminal-error",
                status=500,
                detail="Lesson generation failed; retry with the same request",
            )
            return
    yield event("completed")


@router.get("/generations/{request_id}")
async def generation_status(request_id: UUID, actor: Actor, workflow: Workflow) -> dict:
    thread = await workflow.thread(request_id, context(actor))
    return {
        "jsonapi": {"version": "1.1"},
        "data": {
            "type": "lesson-generation",
            "id": str(request_id),
            "attributes": {"status": thread["status"], "result": thread["result"]},
        },
    }


@router.get("/threads")
async def lesson_threads(
    school_id: Annotated[UUID, Query()],
    department_id: Annotated[UUID, Query()],
    actor: Actor,
    workflow: Workflow,
) -> dict:
    if school_id.int == 0 or department_id.int == 0:
        raise CurriculumError(422, "School and department must be nonzero UUIDs")
    threads = await workflow.threads(school_id, department_id, context(actor))
    return {
        "jsonapi": {"version": "1.1"},
        "data": [thread_resource(t["request"]["request_id"], t) for t in threads],
        "meta": {"limit": 50},
    }


@router.get("/threads/{request_id}")
async def lesson_thread(request_id: UUID, actor: Actor, workflow: Workflow) -> dict:
    thread = await workflow.thread(request_id, context(actor))
    return {"jsonapi": {"version": "1.1"}, "data": thread_resource(str(request_id), thread)}
