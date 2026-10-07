import asyncio
import json
import os
from dataclasses import replace
from uuid import uuid4

import psycopg
import pytest
from langgraph.checkpoint.postgres.aio import AsyncPostgresSaver

from adapters.llm.lesson_schema import assemble_lesson
from domain.entities.curriculum import CurriculumError
from domain.entities.lesson_generation import GenerationContext
from infrastructure.lesson_workflow import LessonWorkflow
from tests.lesson_fixture import generation_fixture as fixture

DSN = os.environ.get("CURRICULUM_TEST_DATABASE_URL")
pytestmark = pytest.mark.skipif(
    not DSN, reason="Set CURRICULUM_TEST_DATABASE_URL for workflow proof"
)


class CurriculumFixture:
    def __init__(self, evidence):
        self.evidence = evidence
        self.withdrawn = False

    async def search(self, token, scope, question, limit):
        return [self.evidence] if scope == self.evidence.source.scope else []

    async def source(self, token, source_id):
        return (
            replace(self.evidence.source, status="withdrawn")
            if self.withdrawn
            else self.evidence.source
        )


class ModelFixture:
    def __init__(self, output):
        self.output = output
        self.calls = 0

    async def generate(self, request, evidence):
        self.calls += 1
        return assemble_lesson(json.dumps(self.output), request, evidence, "fixture-completion")

    async def close(self):
        pass


class GoFixture:
    def __init__(self):
        self.receipts = {}
        self.calls = 0
        self.fail_once = False
        self.revoked = False
        self.tokens = []

    async def authorize(self, token, scope):
        self.tokens.append(token)
        if self.revoked:
            raise CurriculumError(403, "Teaching permission revoked")

    async def create(self, token, scope, request_id, draft):
        await self.authorize(token, scope)
        self.calls += 1
        if request_id not in self.receipts:
            self.receipts[request_id] = {"planID": str(uuid4()), "version": 1, "title": draft.title}
        if self.fail_once:
            self.fail_once = False
            raise CurriculumError(503, "Ambiguous Go response after commit")
        return self.receipts[request_id]

    async def close(self):
        pass


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
            status = await restarted.status(
                request.request_id, GenerationContext(actor, "private-fresh-token")
            )
            assert status["status"] == "completed"
            with pytest.raises(CurriculumError) as unauthorized:
                await restarted.status(
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

            async def observed_authorize(token, scope):
                await initial_authorize(token, scope)
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
