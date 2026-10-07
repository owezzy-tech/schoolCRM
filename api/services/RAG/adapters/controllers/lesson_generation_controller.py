from typing import Annotated
from uuid import UUID

from fastapi import APIRouter, Depends, Request
from pydantic import Field, model_validator

from adapters.controllers.curriculum_controller import JSONAPIResponse, ScopeDTO
from domain.entities.curriculum import CurriculumError
from domain.entities.lesson_generation import GenerationContext, LessonGenerationRequest
from infrastructure.auth import AuthContext, get_auth_context
from infrastructure.lesson_workflow import LessonWorkflow

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


@router.post("/generate")
async def generate(request: GenerateRequest, actor: Actor, workflow: Workflow) -> dict:
    result = await workflow.execute(
        LessonGenerationRequest(
            request.request_id, request.scope(), request.topic, request.duration_minutes
        ),
        context(actor),
    )
    return {
        "jsonapi": {"version": "1.1"},
        "data": {"type": "lesson-generation", "id": str(request.request_id), "attributes": result},
    }


@router.get("/generations/{request_id}")
async def generation_status(request_id: UUID, actor: Actor, workflow: Workflow) -> dict:
    result = await workflow.status(request_id, context(actor))
    return {
        "jsonapi": {"version": "1.1"},
        "data": {"type": "lesson-generation", "id": str(request_id), "attributes": result},
    }
