import asyncio
import json
from contextlib import asynccontextmanager
from uuid import uuid4

import pytest
from fastapi.testclient import TestClient
from langgraph.checkpoint.memory import InMemorySaver

from adapters.controllers.lesson_generation_controller import get_workflow
from domain.entities.curriculum import CurriculumError
from domain.entities.lesson_generation import GenerationContext
from infrastructure.app import build_app
from infrastructure.auth import AuthContext, get_auth_context
from infrastructure.lesson_workflow import LessonWorkflow
from tests.lesson_fixture import CurriculumFixture, GoFixture, ModelFixture, generation_fixture


class RunsFixture:
    """In-memory stand-in for the PostgreSQL actor/request binding."""

    def __init__(self):
        self.rows = {}
        self.latest_calls = 0
        self.get_calls = 0

    @asynccontextmanager
    async def command(self, actor, request_id, input_hash, request):
        row = self.rows.setdefault(
            (actor, request_id), {"request": request, "result": None, "input_hash": input_hash}
        )
        if row["input_hash"] != input_hash:
            raise CurriculumError(409, "request_id already used with different generation input")
        yield None, row

    async def finish(self, conn, actor, request_id, result):
        self.rows[(actor, request_id)]["result"] = result

    async def get(self, actor, request_id):
        self.get_calls += 1
        if (actor, request_id) not in self.rows:
            raise CurriculumError(404, "Generation request not found")
        return self.rows[(actor, request_id)]

    async def scope(self, actor, request_id):
        row = self.rows.get((actor, request_id))
        if row is None:
            raise CurriculumError(404, "Generation request not found")
        scope = row["request"]["scope"]
        return scope["school_id"], scope["department_id"]

    async def latest(self, actor, school_id, department_id):
        self.latest_calls += 1
        return [row for (owner, _), row in self.rows.items() if owner == actor]


async def started_workflow(curriculum, model, go):
    workflow = LessonWorkflow("unused", curriculum, model, go)
    workflow.runs = RunsFixture()
    await workflow.start(InMemorySaver())
    return workflow


async def collect(workflow, request, context):
    return [event async for event in await workflow.stream(request, context)]


def test_stream_sends_progress_and_trusted_citations_before_model_and_result_after_receipt():
    async def run():
        request, evidence, output = generation_fixture()
        actor = uuid4()
        go = GoFixture()
        observed = []
        saved_before_result = []

        class ObservedModel(ModelFixture):
            async def generate(self, request, evidence):
                self.seen = list(observed)
                return await super().generate(request, evidence)

        model = ObservedModel(output)
        workflow = await started_workflow(CurriculumFixture(evidence), model, go)
        context = GenerationContext(actor, "private-token")
        async for name, fields in await workflow.stream(request, context):
            observed.append((name, fields))
            if name == "structured-result":
                row = await workflow.runs.get(actor, request.request_id)
                saved_before_result.append(row["result"])

        assert [name for name, _ in observed] == [
            "progress",
            "citation",
            "progress",
            "progress",
            "structured-result",
        ]
        assert [f["stage"] for n, f in observed if n == "progress"] == [
            "retrieve",
            "generate",
            "persist",
        ]
        assert [name for name, _ in model.seen][:2] == ["progress", "citation"]
        [citation] = observed[1][1]["citations"]
        assert citation["id"] == f"{evidence.source.id}:{evidence.passage.ordinal}"
        assert citation["passage"] == evidence.passage.text
        assert citation["sourceSHA256"] == evidence.source.sha256
        assert citation["page"] == evidence.passage.page
        result = observed[-1][1]["result"]
        assert result == go.receipts[request.request_id]
        assert saved_before_result == [result]
        assert "private-token" not in json.dumps(observed)

        replay = await collect(workflow, request, GenerationContext(actor, "fresh-token"))
        assert [name for name, _ in replay] == ["citation", "structured-result"]
        assert replay[-1][1]["result"] == result
        assert model.calls == 1 and go.calls == 1

    asyncio.run(run())


def test_revoked_teacher_cannot_read_thread_or_history():
    async def run():
        request, evidence, output = generation_fixture()
        actor = uuid4()
        go = GoFixture()
        workflow = await started_workflow(CurriculumFixture(evidence), ModelFixture(output), go)
        context = GenerationContext(actor, "private-token")
        await workflow.execute(request, context)

        thread = await workflow.thread(request.request_id, context)
        assert thread["status"] == "completed"
        assert thread["request"] == {
            "request_id": str(request.request_id),
            "school_id": str(request.scope.school_id),
            "department_id": str(request.scope.department_id),
            "framework": "kenya-cbc",
            "stage": "grade-1",
            "subject": "english",
            "revision": "2024",
            "topic": request.topic,
            "duration_minutes": 30,
        }

        reads_before_revocation = workflow.runs.get_calls
        go.revoked = True
        for read in (
            workflow.thread(request.request_id, context),
            workflow.threads(request.scope.school_id, request.scope.department_id, context),
            workflow.stream(request, context),
        ):
            with pytest.raises(CurriculumError) as denied:
                await read
            assert denied.value.status_code == 403
        assert workflow.runs.latest_calls == 0
        assert workflow.runs.get_calls == reads_before_revocation

    asyncio.run(run())


class StreamingWorkflowFixture:
    def __init__(self, events=(), failure=None, denial=None):
        self.events = events
        self.failure = failure
        self.denial = denial
        self.iterated = False

    async def stream(self, request, context):
        if self.denial:
            raise self.denial
        return self._events()

    async def _events(self):
        self.iterated = True
        for event in self.events:
            yield event
        if self.failure:
            raise self.failure

    async def thread(self, request_id, context):
        return {"request": {"request_id": str(request_id)}, "status": "pending", "result": None}

    async def threads(self, school_id, department_id, context):
        return [await self.thread(uuid4(), context)]


def client_for(workflow, actor=None):
    app = build_app()
    actor = actor or AuthContext(str(uuid4()), "private-token", [])
    app.dependency_overrides[get_auth_context] = lambda: actor
    app.dependency_overrides[get_workflow] = lambda: workflow
    return TestClient(app)


def body():
    return {
        "request_id": str(uuid4()),
        "school_id": str(uuid4()),
        "department_id": str(uuid4()),
        "framework": "kenya-cbc",
        "stage": "grade-1",
        "subject": "english",
        "revision": "2024",
        "topic": "School vocabulary",
        "duration_minutes": 30,
    }


def sse(text):
    events = []
    for block in text.strip().split("\n\n"):
        lines = dict(line.split(": ", 1) for line in block.splitlines() if ": " in line)
        if "event" in lines:
            events.append((lines["event"], json.loads(lines["data"])))
    return events


def test_stream_endpoint_emits_versioned_events_and_one_terminal_completion():
    receipt = {"planID": str(uuid4()), "version": 1, "title": "Classroom objects"}
    workflow = StreamingWorkflowFixture(
        [
            ("progress", {"stage": "retrieve", "detail": "Retrieving"}),
            ("citation", {"citations": [{"id": "source:0", "passage": "text"}]}),
            ("structured-result", {"result": receipt}),
        ]
    )
    request = body()
    response = client_for(workflow).post("/v1/rag/lessons/generate/stream", json=request)
    assert response.status_code == 200
    assert response.headers["content-type"].startswith("text/event-stream")
    events = sse(response.text)
    assert [name for name, _ in events] == [
        "progress",
        "citation",
        "structured-result",
        "approval-request",
        "completed",
    ]
    assert all(
        data["protocolVersion"] == 1 and data["requestID"] == request["request_id"]
        for _, data in events
    )
    assert events[2][1]["result"] == receipt
    assert events[3][1] == {
        "protocolVersion": 1,
        "requestID": request["request_id"],
        "planID": receipt["planID"],
        "version": 1,
        "action": "submit-for-review",
    }
    assert "private-token" not in response.text


@pytest.mark.parametrize(
    ("failure", "status"),
    [
        (CurriculumError(409, "Generation checkpoint belongs to different input"), 409),
        (RuntimeError("postgresql://rag:secret-password@db/rag"), 500),
    ],
)
def test_stream_failures_after_headers_end_with_one_client_safe_terminal_error(failure, status):
    workflow = StreamingWorkflowFixture(
        [("progress", {"stage": "retrieve", "detail": "Retrieving"})], failure=failure
    )
    response = client_for(workflow).post("/v1/rag/lessons/generate/stream", json=body())
    events = sse(response.text)
    assert [name for name, _ in events] == ["progress", "terminal-error"]
    assert events[-1][1]["status"] == status
    if isinstance(failure, CurriculumError):
        assert events[-1][1]["detail"] == failure.detail
    assert "secret-password" not in response.text and "postgresql" not in response.text


def test_stream_authorization_failures_return_json_api_errors_before_streaming():
    workflow = StreamingWorkflowFixture(
        denial=CurriculumError(403, "Lesson generation requires current department teaching access")
    )
    denied = client_for(workflow).post("/v1/rag/lessons/generate/stream", json=body())
    assert denied.status_code == 403
    assert denied.headers["content-type"].startswith("application/vnd.api+json")
    assert denied.json()["errors"][0]["detail"].startswith("Lesson generation requires")
    anonymous = client_for(workflow, AuthContext("anonymous", None, [])).post(
        "/v1/rag/lessons/generate/stream", json=body()
    )
    assert anonymous.status_code == 401 and "errors" in anonymous.json()
    invalid = client_for(workflow).post(
        "/v1/rag/lessons/generate/stream", json={**body(), "duration_minutes": 0}
    )
    assert invalid.status_code == 422 and "errors" in invalid.json()
    assert not workflow.iterated


def test_thread_endpoints_return_json_api_lesson_threads():
    workflow = StreamingWorkflowFixture()
    client = client_for(workflow)
    request_id, school, department = uuid4(), uuid4(), uuid4()
    one = client.get(f"/v1/rag/lessons/threads/{request_id}")
    assert one.status_code == 200
    assert one.headers["content-type"].startswith("application/vnd.api+json")
    assert one.json()["data"] == {
        "type": "lesson-thread",
        "id": str(request_id),
        "attributes": {
            "request": {"request_id": str(request_id)},
            "status": "pending",
            "result": None,
        },
    }
    listed = client.get(
        "/v1/rag/lessons/threads", params={"school_id": school, "department_id": department}
    )
    assert listed.status_code == 200 and listed.json()["data"][0]["type"] == "lesson-thread"
    assert client.get("/v1/rag/lessons/threads").status_code == 422
    status = client.get(f"/v1/rag/lessons/generations/{request_id}")
    assert status.json()["data"]["attributes"] == {"status": "pending", "result": None}


def test_closing_response_at_yield_closes_the_underlying_workflow_stream():
    from adapters.controllers.lesson_generation_controller import GenerateRequest, generate_stream

    async def run():
        closed = asyncio.Event()

        async def events():
            try:
                yield "progress", {"stage": "retrieve", "detail": "Retrieving"}
                await asyncio.Event().wait()
            finally:
                closed.set()

        response = generate_stream(GenerateRequest(**body()), events())
        first = await anext(response)
        assert first.event == "progress"
        await response.aclose()
        assert closed.is_set()

    asyncio.run(run())
