import asyncio
import os
from dataclasses import replace
from uuid import uuid4

import psycopg
import pytest
from fastapi.testclient import TestClient
from langgraph.checkpoint.postgres.aio import AsyncPostgresSaver

from adapters.controllers.lesson_generation_controller import get_workflow
from domain.entities.curriculum import CurriculumError
from domain.entities.lesson_generation import GenerationContext
from infrastructure.app import build_app
from infrastructure.auth import AuthContext, get_auth_context
from infrastructure.lesson_workflow import LessonWorkflow
from tests.lesson_fixture import CurriculumFixture, GoFixture, ModelFixture
from tests.lesson_fixture import generation_fixture as fixture

DSN = os.environ.get("CURRICULUM_TEST_DATABASE_URL")
pytestmark = pytest.mark.skipif(
    not DSN, reason="Set CURRICULUM_TEST_DATABASE_URL for workflow proof"
)


async def prepare(workflow):
    await workflow.runs.migrate()
    async with AsyncPostgresSaver.from_conn_string(DSN) as saver:
        await saver.setup()
    await workflow.start()


def test_resume_after_ambiguous_write_uses_checkpoint_and_fresh_token_without_duplicate_plan():
    async def run():
        request, evidence, output = fixture()
        actor = uuid4()
        curriculum, model, go = CurriculumFixture(evidence), ModelFixture(output), GoFixture()
        workflow = LessonWorkflow(DSN, curriculum, model, go)
        await prepare(workflow)
        go.fail_once = True
        try:
            with pytest.raises(CurriculumError, match="Ambiguous"):
                await workflow.execute(request, GenerationContext(actor, "private-first-token"))
            assert model.calls == 1 and len(go.receipts) == 1
        finally:
            await workflow.close()
        restarted_model = ModelFixture(output)
        restarted = LessonWorkflow(DSN, curriculum, restarted_model, go)
        await restarted.start()
        try:
            result = await restarted.execute(
                request, GenerationContext(actor, "private-fresh-token")
            )
            assert restarted_model.calls == 0
            assert len(go.receipts) == 1 and go.calls == 2
            assert result["planID"] == go.receipts[request.request_id]["planID"]
            assert "private-fresh-token" in go.tokens
            status = await restarted.thread(
                request.request_id, GenerationContext(actor, "private-fresh-token")
            )
            assert status["status"] == "completed"
            with pytest.raises(CurriculumError) as unauthorized:
                await restarted.thread(
                    request.request_id, GenerationContext(uuid4(), "other-token")
                )
            assert unauthorized.value.status_code == 404
            # Check JSONB and binary checkpoint storage, not just the returned DTO.
            async with await psycopg.AsyncConnection.connect(DSN) as conn:
                thread = f"{actor}:{request.request_id}"
                cursor = await conn.execute(
                    "SELECT checkpoint::text FROM checkpoints WHERE thread_id=%s", (thread,)
                )
                stored = [row[0].encode() for row in await cursor.fetchall()]
                cursor = await conn.execute(
                    "SELECT blob FROM checkpoint_blobs WHERE thread_id=%s", (thread,)
                )
                stored += [bytes(row[0]) for row in await cursor.fetchall() if row[0] is not None]
                assert all(
                    b"private-first-token" not in data and b"private-fresh-token" not in data
                    for data in stored
                )
        finally:
            await restarted.close()

    asyncio.run(run())


def test_concurrent_retries_changed_input_and_revocation_are_enforced():
    async def run():
        request, evidence, output = fixture()
        actor = uuid4()
        context = GenerationContext(actor, "private-token")
        model, go = ModelFixture(output), GoFixture()
        workflow = LessonWorkflow(DSN, CurriculumFixture(evidence), model, go)
        await prepare(workflow)
        try:
            results = await asyncio.gather(*[workflow.execute(request, context) for _ in range(4)])
            assert len({result["planID"] for result in results}) == 1
            assert model.calls == 1 and go.calls == 1
            with pytest.raises(CurriculumError) as changed:
                await workflow.execute(replace(request, topic="Different topic"), context)
            assert changed.value.status_code == 409
            go.revoked = True
            with pytest.raises(CurriculumError) as revoked:
                await workflow.execute(request, context)
            assert revoked.value.status_code == 403
            assert go.calls == 1
        finally:
            await workflow.close()

    asyncio.run(run())


def test_withdrawn_evidence_blocks_resume_before_model_or_write():
    async def run():
        request, evidence, output = fixture()
        actor = uuid4()
        curriculum, model, go = CurriculumFixture(evidence), ModelFixture(output), GoFixture()
        workflow = LessonWorkflow(DSN, curriculum, model, go)
        await prepare(workflow)
        go.fail_once = True
        try:
            with pytest.raises(CurriculumError):
                await workflow.execute(request, GenerationContext(actor, "private-token"))
            curriculum.withdrawn = True
            with pytest.raises(CurriculumError, match="withdrawn"):
                await workflow.execute(request, GenerationContext(actor, "private-token"))
            assert model.calls == 1 and go.calls == 1
        finally:
            await workflow.close()

    asyncio.run(run())


def test_cached_reply_rechecks_permission_after_waiting_for_request_lock():
    async def run():
        request, evidence, output = fixture()
        actor = uuid4()
        context = GenerationContext(actor, "private-token")
        model, go = ModelFixture(output), GoFixture()
        workflow = LessonWorkflow(DSN, CurriculumFixture(evidence), model, go)
        await prepare(workflow)
        try:
            await workflow.execute(request, context)
            checked = asyncio.Event()
            initial_authorize = go.authorize

            async def observed_authorize(token, school_id, department_id):
                await initial_authorize(token, school_id, department_id)
                checked.set()

            go.authorize = observed_authorize
            async with await psycopg.AsyncConnection.connect(DSN, autocommit=True) as conn:
                key = f"rag-lesson-generation:{actor}:{request.request_id}"
                await conn.execute("SELECT pg_advisory_lock(hashtextextended(%s,0))", (key,))
                replay = asyncio.create_task(workflow.execute(request, context))
                await checked.wait()
                go.revoked = True
                await conn.execute("SELECT pg_advisory_unlock(hashtextextended(%s,0))", (key,))
                with pytest.raises(CurriculumError) as denied:
                    await replay
                assert denied.value.status_code == 403 and go.calls == 1
        finally:
            await workflow.close()

    asyncio.run(run())


def test_cancelled_stream_releases_lock_and_reconnect_resumes_checkpoint():
    async def run():
        request, evidence, output = fixture()
        actor = uuid4()
        model_started, release_model = asyncio.Event(), asyncio.Event()

        class BlockingModel(ModelFixture):
            async def generate(self, request, evidence):
                model_started.set()
                await release_model.wait()
                return await super().generate(request, evidence)

        model, go = BlockingModel(output), GoFixture()
        workflow = LessonWorkflow(DSN, CurriculumFixture(evidence), model, go)
        await prepare(workflow)
        try:
            seen = []

            async def consume(token):
                async for event in await workflow.stream(request, GenerationContext(actor, token)):
                    seen.append(event)

            disconnected = asyncio.create_task(consume("private-first-token"))
            await model_started.wait()
            assert [name for name, _ in seen][:2] == ["progress", "citation"]
            disconnected.cancel()
            with pytest.raises(asyncio.CancelledError):
                await disconnected
            key = f"rag-lesson-generation:{actor}:{request.request_id}"
            async with await psycopg.AsyncConnection.connect(DSN, autocommit=True) as conn:
                for _ in range(50):
                    cursor = await conn.execute(
                        "SELECT pg_try_advisory_lock(hashtextextended(%s,0))", (key,)
                    )
                    if (await cursor.fetchone())[0]:
                        break
                    await asyncio.sleep(0.1)
                else:
                    pytest.fail("Cancelled stream kept the request lock")
                await conn.execute("SELECT pg_advisory_unlock(hashtextextended(%s,0))", (key,))
            pending = await workflow.thread(request.request_id, GenerationContext(actor, "t"))
            assert pending["status"] == "pending" and pending["result"] is None

            release_model.set()
            seen.clear()
            await consume("private-fresh-token")
            # Retrieval was checkpointed: resume starts with saved citations, not a new search.
            assert [name for name, _ in seen] == [
                "citation",
                "progress",
                "progress",
                "structured-result",
            ]
            assert model.calls == 1 and go.calls == 1
            seen.clear()
            await consume("private-replay-token")
            assert [name for name, _ in seen] == ["citation", "structured-result"]
            assert go.calls == 1
        finally:
            await workflow.close()

    asyncio.run(run())


def test_thread_history_is_actor_private_department_scoped_and_rechecks_authority():
    async def run():
        request, evidence, output = fixture()
        actor = uuid4()
        go = GoFixture()
        workflow = LessonWorkflow(DSN, CurriculumFixture(evidence), ModelFixture(output), go)
        await prepare(workflow)
        try:
            context = GenerationContext(actor, "private-token")
            await workflow.execute(request, context)
            scope = request.scope
            threads = await workflow.threads(scope.school_id, scope.department_id, context)
            assert [t["request"]["request_id"] for t in threads] == [str(request.request_id)]
            assert threads[0]["status"] == "completed"
            assert threads[0]["request"]["department_id"] == str(scope.department_id)
            other = GenerationContext(uuid4(), "other-token")
            assert await workflow.threads(scope.school_id, scope.department_id, other) == []
            with pytest.raises(CurriculumError) as hidden:
                await workflow.thread(request.request_id, other)
            assert hidden.value.status_code == 404
            assert await workflow.threads(scope.school_id, uuid4(), context) == []
            go.revoked = True
            with pytest.raises(CurriculumError) as revoked:
                await workflow.threads(scope.school_id, scope.department_id, context)
            assert revoked.value.status_code == 403
        finally:
            await workflow.close()

    asyncio.run(run())


def test_http_stream_runs_workflow_to_completion_and_replays_receipt():
    request, evidence, output = fixture()
    actor = uuid4()
    model, go = ModelFixture(output), GoFixture()
    workflow = LessonWorkflow(DSN, CurriculumFixture(evidence), model, go)
    app = build_app()
    app.dependency_overrides[get_auth_context] = lambda: AuthContext(str(actor), "token", [])
    app.dependency_overrides[get_workflow] = lambda: workflow
    scope = request.scope
    body = {
        "request_id": str(request.request_id),
        "school_id": str(scope.school_id),
        "department_id": str(scope.department_id),
        "framework": scope.framework,
        "stage": scope.stage,
        "subject": scope.subject,
        "revision": scope.revision,
        "topic": request.topic,
        "duration_minutes": request.duration_minutes,
    }
    with TestClient(app) as client:
        client.portal.call(prepare, workflow)
        try:
            names = []
            for _ in range(2):
                response = client.post("/v1/rag/lessons/generate/stream", json=body)
                assert response.status_code == 200
                names.append(
                    [line[7:] for line in response.text.splitlines() if line.startswith("event:")]
                )
            assert names[0] == [
                "progress",
                "citation",
                "progress",
                "progress",
                "structured-result",
                "approval-request",
                "completed",
            ]
            assert names[1] == ["citation", "structured-result", "approval-request", "completed"]
            assert model.calls == 1 and go.calls == 1
            thread = client.get(f"/v1/rag/lessons/threads/{request.request_id}").json()
            assert thread["data"]["attributes"]["status"] == "completed"
        finally:
            client.portal.call(workflow.close)
