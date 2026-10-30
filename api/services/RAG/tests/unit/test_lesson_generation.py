import asyncio
import json
from dataclasses import asdict
from uuid import uuid4

import httpx
import pytest
from pydantic import SecretStr

from adapters.llm.deepseek_lesson_generator import DeepSeekLessonGenerator
from adapters.llm.lesson_schema import assemble_lesson
from domain.entities.curriculum import CurriculumError
from tests.lesson_fixture import generation_fixture as fixture


@pytest.mark.parametrize("failure", ["duration", "citation", "commands", "empty", "type"])
def test_untrusted_model_output_cannot_become_a_storable_lesson(failure: str) -> None:
    request, evidence, output = fixture()
    if failure == "duration":
        output["activities"][0]["minutes"] = 29
    if failure == "citation":
        output["citation_ids"] = [str(uuid4())]
    if failure == "commands":
        output["publish"] = True
    if failure == "empty":
        output["citation_ids"] = []
    if failure == "type":
        output["activities"][0]["minutes"] = "30"
    with pytest.raises(CurriculumError) as failure_info:
        assemble_lesson(json.dumps(output), request, [evidence], "completion-1")
    assert failure_info.value.status_code == 422


def test_citations_and_model_metadata_come_from_trusted_inputs() -> None:
    request, evidence, output = fixture()
    draft = assemble_lesson(json.dumps(output), request, [evidence], "completion-1")
    citation = draft.content["citations"][0]
    assert citation["sourceID"] == str(evidence.source.id)
    assert citation["page"] == 17 and citation["sourceSHA256"] == "a" * 64
    assert draft.content["durationMinutes"] == 30
    assert draft.content["generation"]["modelVersion"] == "DeepSeek-V4.1-Flash"
    assert draft.content["generation"]["completionID"] == "completion-1"


def test_deepseek_uses_verified_model_json_mode_and_never_prompts_with_credentials() -> None:
    async def run():
        request, evidence, output = fixture()
        requests = []

        def handler(req: httpx.Request):
            requests.append(req.url.path)
            assert req.headers["Authorization"] == "Bearer private-key"
            if req.url.path == "/models":
                return httpx.Response(
                    200, json={"data": [{"id": "deepseek-flash", "name": "DeepSeek-V4.1-Flash"}]}
                )
            body = json.loads(req.content)
            assert body["model"] == "deepseek-flash"
            assert body["response_format"] == {"type": "json_object"}
            assert body["thinking"] == {"type": "disabled"}
            assert b"private-key" not in req.content
            assert b"Bearer" not in req.content
            return httpx.Response(
                200,
                json={
                    "id": "completion-1",
                    "model": "deepseek-flash",
                    "choices": [
                        {"finish_reason": "stop", "message": {"content": json.dumps(output)}}
                    ],
                },
            )

        generator = DeepSeekLessonGenerator(SecretStr("private-key"), httpx.MockTransport(handler))
        try:
            result = await generator.generate(request, [evidence])
            assert result.title == output["title"]
            assert requests == ["/models", "/chat/completions"]
            assert "private-key" not in json.dumps(asdict(result))
        finally:
            await generator.close()

    asyncio.run(run())


def test_deepseek_rejects_a_changed_managed_model_version_without_generating() -> None:
    async def run():
        request, evidence, _ = fixture()
        calls = []

        def handler(req: httpx.Request):
            calls.append(req.url.path)
            return httpx.Response(
                200, json={"data": [{"id": "deepseek-flash", "name": "different-version"}]}
            )

        generator = DeepSeekLessonGenerator(SecretStr("private-key"), httpx.MockTransport(handler))
        try:
            with pytest.raises(CurriculumError):
                await generator.generate(request, [evidence])
            assert calls == ["/models"]
        finally:
            await generator.close()

    asyncio.run(run())
