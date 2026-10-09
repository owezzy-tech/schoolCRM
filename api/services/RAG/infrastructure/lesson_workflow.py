import json
from collections.abc import AsyncGenerator
from contextlib import AsyncExitStack, aclosing
from dataclasses import asdict
from hashlib import sha256
from typing import TypedDict
from uuid import UUID

from langgraph.checkpoint.base import BaseCheckpointSaver
from langgraph.checkpoint.postgres.aio import AsyncPostgresSaver
from langgraph.graph import END, START, StateGraph
from langgraph.runtime import Runtime
from langsmith import tracing_context

from adapters.llm.lesson_schema import evidence_citation
from adapters.repositories.lesson_runs import LessonRuns
from domain.entities.curriculum import (
    CurriculumError,
    CurriculumEvidence,
    CurriculumPassage,
    CurriculumScope,
    CurriculumSource,
)
from domain.entities.lesson_generation import (
    GeneratedLesson,
    GenerationContext,
    LessonGenerationRequest,
)
from domain.ports.lesson_generation import LessonGenerator, LessonWriter
from use_cases.curriculum import CurriculumService


class LessonState(TypedDict, total=False):
    request: dict
    evidence: list[dict]
    draft: dict
    result: dict


def request_data(request: LessonGenerationRequest) -> dict:
    return json.loads(json.dumps(asdict(request), default=str))


def parse_request(value: dict) -> LessonGenerationRequest:
    scope = value["scope"].copy()
    scope["school_id"] = UUID(scope["school_id"])
    scope["department_id"] = UUID(scope["department_id"])
    return LessonGenerationRequest(
        UUID(value["request_id"]),
        CurriculumScope(**scope),
        value["topic"],
        value["duration_minutes"],
    )


def flatten_request(value: dict) -> dict:
    return {
        "request_id": value["request_id"],
        **value["scope"],
        "topic": value["topic"],
        "duration_minutes": value["duration_minutes"],
    }


def thread_view(row: dict) -> dict:
    return {
        "request": flatten_request(row["request"]),
        "status": "completed" if row["result"] is not None else "pending",
        "result": row["result"],
    }


def progress(stage: str, detail: str) -> tuple[str, dict]:
    return "progress", {"stage": stage, "detail": detail}


def citations(evidence: list[dict]) -> tuple[str, dict]:
    return "citation", {"citations": [evidence_citation(parse_evidence(e)) for e in evidence]}


def parse_evidence(value: dict) -> CurriculumEvidence:
    source = value["source"].copy()
    scope = source.pop("scope").copy()
    scope["school_id"] = UUID(scope["school_id"])
    scope["department_id"] = UUID(scope["department_id"])
    source["id"] = UUID(source["id"])
    return CurriculumEvidence(
        CurriculumSource(**source, scope=CurriculumScope(**scope)),
        CurriculumPassage(**value["passage"]),
        value["score"],
        value["model_identity"],
        UUID(value["index_revision"]),
    )


class LessonWorkflow:
    def __init__(
        self,
        dsn: str,
        curriculum: CurriculumService,
        generator: LessonGenerator,
        writer: LessonWriter,
    ):
        self.dsn = dsn
        self.curriculum = curriculum
        self.generator = generator
        self.writer = writer
        self.runs = LessonRuns(dsn)
        self.stack = AsyncExitStack()
        self.graph = None

    async def start(self, checkpointer: BaseCheckpointSaver | None = None) -> None:
        if checkpointer is None:
            checkpointer = await self.stack.enter_async_context(
                AsyncPostgresSaver.from_conn_string(self.dsn)
            )
        builder = StateGraph(LessonState, context_schema=GenerationContext)
        builder.add_node("retrieve", self._retrieve)
        builder.add_node("generate", self._generate)
        builder.add_node("persist", self._persist)
        builder.add_edge(START, "retrieve")
        builder.add_edge("retrieve", "generate")
        builder.add_edge("generate", "persist")
        builder.add_edge("persist", END)
        self.graph = builder.compile(checkpointer=checkpointer)

    async def stream(
        self, request: LessonGenerationRequest, context: GenerationContext
    ) -> AsyncGenerator[tuple[str, dict], None]:
        """Authorize now, then return the run's events; the caller must consume them."""
        if self.graph is None:
            raise CurriculumError(503, "Lesson workflow is not running")
        await self._authorize(context, request.scope)
        return self._events(request, context)

    async def execute(self, request: LessonGenerationRequest, context: GenerationContext) -> dict:
        async for event, data in await self.stream(request, context):
            if event == "structured-result":
                result = data["result"]
        return result

    async def _events(
        self, request: LessonGenerationRequest, context: GenerationContext
    ) -> AsyncGenerator[tuple[str, dict], None]:
        payload = request_data(request)
        fingerprint = sha256(json.dumps(payload, sort_keys=True).encode()).hexdigest()
        config = {"configurable": {"thread_id": f"{context.actor_id}:{request.request_id}"}}
        async with self.runs.command(
            context.actor_id, request.request_id, fingerprint, payload
        ) as (conn, row):
            # A concurrent request may have waited behind a long model call.
            await self._authorize(context, request.scope)
            # Tokens live in Runtime.context only. Disable external tracing for
            # this private workflow rather than exporting student/source state.
            with tracing_context(enabled=False):
                checkpoint = await self.graph.aget_state(config)
                if checkpoint.values.get("request") not in (None, payload):
                    raise CurriculumError(409, "Generation checkpoint belongs to different input")
                if checkpoint.values.get("evidence"):
                    yield citations(checkpoint.values["evidence"])
                result = row["result"] or checkpoint.values.get("result")
                if result is None:
                    async with aclosing(
                        self.graph.astream(
                            None if checkpoint.next else {"request": payload},
                            config,
                            context=context,
                            durability="sync",
                            stream_mode="custom",
                        )
                    ) as events:
                        async for event in events:
                            yield event
                    result = (await self.graph.aget_state(config)).values["result"]
            if row["result"] is None:
                await self._authorize(context, request.scope)
                await self.runs.finish(conn, context.actor_id, request.request_id, result)
            yield "structured-result", {"result": result}

    async def thread(self, request_id: UUID, context: GenerationContext) -> dict:
        school_id, department_id = await self.runs.scope(context.actor_id, request_id)
        await self.writer.authorize(context.token, school_id, department_id)
        row = await self.runs.get(context.actor_id, request_id)
        return thread_view(row)

    async def threads(
        self, school_id: UUID, department_id: UUID, context: GenerationContext
    ) -> list[dict]:
        await self.writer.authorize(context.token, school_id, department_id)
        rows = await self.runs.latest(context.actor_id, school_id, department_id)
        return [thread_view(row) for row in rows]

    async def _authorize(self, context: GenerationContext, scope: CurriculumScope) -> None:
        await self.writer.authorize(context.token, scope.school_id, scope.department_id)

    async def _retrieve(self, state: LessonState, runtime: Runtime[GenerationContext]) -> dict:
        runtime.stream_writer(progress("retrieve", "Retrieving approved curriculum evidence"))
        request = parse_request(state["request"])
        await self._authorize(runtime.context, request.scope)
        evidence = await self.curriculum.search(
            runtime.context.token, request.scope, request.topic, 5
        )
        if not evidence:
            raise CurriculumError(
                422, "No approved curriculum evidence supports this lesson request"
            )
        stored = json.loads(json.dumps([asdict(e) for e in evidence], default=str))
        runtime.stream_writer(citations(stored))
        return {"evidence": stored}

    async def _current_evidence(
        self, state: LessonState, context: GenerationContext
    ) -> list[CurriculumEvidence]:
        request = parse_request(state["request"])
        await self._authorize(context, request.scope)
        evidence = [parse_evidence(e) for e in state["evidence"]]
        for item in evidence:
            current = await self.curriculum.source(context.token, item.source.id)
            if (
                current.status != "approved"
                or current.scope != request.scope
                or current.sha256 != item.source.sha256
            ):
                raise CurriculumError(
                    409, "Curriculum evidence was withdrawn or changed; start a new request"
                )
        return evidence

    async def _generate(self, state: LessonState, runtime: Runtime[GenerationContext]) -> dict:
        evidence = await self._current_evidence(state, runtime.context)
        runtime.stream_writer(progress("generate", "Drafting the lesson from cited evidence"))
        draft = await self.generator.generate(parse_request(state["request"]), evidence)
        return {"draft": asdict(draft)}

    async def _persist(self, state: LessonState, runtime: Runtime[GenerationContext]) -> dict:
        runtime.stream_writer(progress("persist", "Saving the draft through Go"))
        await self._current_evidence(state, runtime.context)
        request = parse_request(state["request"])
        result = await self.writer.create(
            runtime.context.token,
            request.scope,
            request.request_id,
            GeneratedLesson(**state["draft"]),
        )
        return {"result": result}

    async def close(self) -> None:
        await self.stack.aclose()
        await self.generator.close()
        await self.writer.close()
