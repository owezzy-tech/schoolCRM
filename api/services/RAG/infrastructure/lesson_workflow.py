import json
from contextlib import AsyncExitStack
from dataclasses import asdict
from hashlib import sha256
from typing import TypedDict
from uuid import UUID

from langgraph.checkpoint.base import BaseCheckpointSaver
from langgraph.checkpoint.postgres.aio import AsyncPostgresSaver
from langgraph.graph import END, START, StateGraph
from langgraph.runtime import Runtime
from langsmith import tracing_context

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

    async def execute(self, request: LessonGenerationRequest, context: GenerationContext) -> dict:
        await self.writer.authorize(context.token, request.scope)
        if self.graph is None:
            raise CurriculumError(503, "Lesson workflow is not running")
        payload = request_data(request)
        fingerprint = sha256(json.dumps(payload, sort_keys=True).encode()).hexdigest()
        config = {"configurable": {"thread_id": f"{context.actor_id}:{request.request_id}"}}
        async with self.runs.command(
            context.actor_id, request.request_id, fingerprint, payload
        ) as (conn, row):
            # A concurrent request may have waited behind a long model call.
            await self.writer.authorize(context.token, request.scope)
            if row["result"] is not None:
                return row["result"]
            # Tokens live in Runtime.context only. Disable external tracing for
            # this private workflow rather than exporting student/source state.
            with tracing_context(enabled=False):
                checkpoint = await self.graph.aget_state(config)
                if checkpoint.values.get("request") not in (None, payload):
                    raise CurriculumError(409, "Generation checkpoint belongs to different input")
                if checkpoint.values.get("result") is not None:
                    result = checkpoint.values["result"]
                else:
                    state = await self.graph.ainvoke(
                        None if checkpoint.next else {"request": payload},
                        config,
                        context=context,
                        durability="sync",
                    )
                    result = state["result"]
            await self.writer.authorize(context.token, request.scope)
            await self.runs.finish(conn, context.actor_id, request.request_id, result)
            return result

    async def status(self, request_id: UUID, context: GenerationContext) -> dict:
        row = await self.runs.get(context.actor_id, request_id)
        request = parse_request(row["request"])
        await self.writer.authorize(context.token, request.scope)
        return {
            "status": "completed" if row["result"] is not None else "pending",
            "result": row["result"],
        }

    async def _retrieve(self, state: LessonState, runtime: Runtime[GenerationContext]) -> dict:
        request = parse_request(state["request"])
        await self.writer.authorize(runtime.context.token, request.scope)
        evidence = await self.curriculum.search(
            runtime.context.token, request.scope, request.topic, 5
        )
        if not evidence:
            raise CurriculumError(
                422, "No approved curriculum evidence supports this lesson request"
            )
        return {"evidence": json.loads(json.dumps([asdict(e) for e in evidence], default=str))}

    async def _current_evidence(
        self, state: LessonState, context: GenerationContext
    ) -> list[CurriculumEvidence]:
        request = parse_request(state["request"])
        await self.writer.authorize(context.token, request.scope)
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
        draft = await self.generator.generate(parse_request(state["request"]), evidence)
        return {"draft": asdict(draft)}

    async def _persist(self, state: LessonState, runtime: Runtime[GenerationContext]) -> dict:
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
