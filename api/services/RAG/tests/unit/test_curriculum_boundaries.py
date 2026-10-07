import asyncio
import json
from dataclasses import replace
from uuid import uuid4

import httpx
import pytest
from pydantic import ValidationError

from adapters.controllers.curriculum_controller import ScopeDTO
from adapters.embeddings.bge_m3 import BGEM3Embeddings
from domain.entities.curriculum import CurriculumError, CurriculumScope
from infrastructure.school_access import SchoolAccessClient


def scope() -> CurriculumScope:
    return CurriculumScope(uuid4(), uuid4(), "kenya-cbc", "grade-1", "english", "2024")


def test_stage_cannot_expand_cambridge_scope_or_infer_equivalence() -> None:
    with pytest.raises(ValidationError):
        ScopeDTO(**vars_scope(replace(scope(), framework="cambridge", stage="grade-1")))
    with pytest.raises(ValidationError):
        ScopeDTO(**vars_scope(replace(scope(), framework="cambridge", stage="year-10")))


def vars_scope(value: CurriculumScope) -> dict:
    from dataclasses import asdict

    return asdict(value)


def test_school_access_rechecks_go_memberships_and_denies_revocation() -> None:
    async def run() -> None:
        target = scope()
        active = True

        def handler(request: httpx.Request) -> httpx.Response:
            assert request.headers["Authorization"] == "Bearer teacher"
            if request.url.path.endswith("/memberships"):
                return httpx.Response(403)
            return httpx.Response(
                200,
                json={
                    "data": [
                        {
                            "attributes": {
                                "schoolID": str(target.school_id),
                                "departmentID": str(target.department_id),
                                "capability": "teach",
                                "active": active,
                            }
                        }
                    ]
                },
            )

        client = SchoolAccessClient("http://school", httpx.MockTransport(handler))
        try:
            await client.authorize("teacher", target, False)
            with pytest.raises(CurriculumError, match="management"):
                await client.authorize("teacher", target, True)
            with pytest.raises(CurriculumError, match="capability"):
                await client.authorize("teacher", replace(target, school_id=uuid4()), False)
            active = False
            with pytest.raises(CurriculumError, match="capability"):
                await client.authorize("teacher", target, False)
            with pytest.raises(CurriculumError, match="authenticated"):
                await client.authorize(None, target, False)
        finally:
            await client.close()

    asyncio.run(run())


def test_manager_must_use_a_department_in_the_authorized_school() -> None:
    async def run() -> None:
        target = scope()

        def handler(request: httpx.Request) -> httpx.Response:
            return httpx.Response(
                200,
                json={
                    "data": [
                        {
                            "id": str(target.department_id),
                            "attributes": {},
                        }
                    ]
                },
            )

        client = SchoolAccessClient("http://school", httpx.MockTransport(handler))
        try:
            await client.authorize("manager", target, True)
            with pytest.raises(CurriculumError, match="Department"):
                await client.authorize("manager", replace(target, department_id=uuid4()), True)
        finally:
            await client.close()

    asyncio.run(run())


@pytest.mark.parametrize("mutation", ["dimensions", "nonfinite", "changed_model", "valid"])
def test_ollama_rejects_invalid_vectors_and_changed_weights(mutation: str) -> None:
    async def run() -> None:
        tag_calls = 0

        def handler(request: httpx.Request) -> httpx.Response:
            nonlocal tag_calls
            if request.url.path == "/api/tags":
                tag_calls += 1
                digest = "changed" if mutation == "changed_model" and tag_calls > 1 else "pinned"
                return httpx.Response(
                    200, json={"models": [{"name": "bge-m3:latest", "digest": digest}]}
                )
            vector = [1.0] * (16 if mutation == "dimensions" else 1024)
            if mutation == "nonfinite":
                vector[0] = "NaN"
            return httpx.Response(200, json={"embeddings": [vector]})

        provider = BGEM3Embeddings(
            "ollama", "http://ollama", transport=httpx.MockTransport(handler)
        )
        try:
            if mutation == "valid":
                identity, vectors = await provider.embed(["curriculum"])
                assert identity == "ollama:bge-m3:pinned"
                assert len(vectors[0]) == 1024
            else:
                with pytest.raises(CurriculumError):
                    await provider.embed(["curriculum"])
        finally:
            await provider.close()

    asyncio.run(run())


def test_cloudflare_uses_separate_identity_and_does_not_send_token_in_payload() -> None:
    async def run() -> None:
        def handler(request: httpx.Request) -> httpx.Response:
            assert request.url.path == "/client/v4/accounts/account/ai/v1/embeddings"
            assert json.loads(request.content) == {
                "model": "@cf/baai/bge-m3",
                "input": ["curriculum"],
            }
            assert request.headers["Authorization"] == "Bearer private-token"
            assert b"private-token" not in request.content
            return httpx.Response(
                200,
                json={"data": [{"index": 0, "embedding": [1.0] * 1024}]},
            )

        provider = BGEM3Embeddings(
            "cloudflare", "http://unused", "account", "private-token", httpx.MockTransport(handler)
        )
        try:
            identity, _ = await provider.embed(["curriculum"])
            assert identity == "cloudflare:@cf/baai/bge-m3:1024"
        finally:
            await provider.close()

    asyncio.run(run())


@pytest.mark.parametrize("duplicate", [False, True])
def test_cloudflare_response_order_and_duplicate_indexes(duplicate: bool) -> None:
    async def run() -> None:
        def handler(request: httpx.Request) -> httpx.Response:
            return httpx.Response(
                200,
                json={
                    "data": [
                        {"index": 0 if duplicate else 1, "embedding": [2.0] * 1024},
                        {"index": 0, "embedding": [1.0] * 1024},
                    ]
                },
            )

        provider = BGEM3Embeddings(
            "cloudflare", "http://unused", "account", "private-token", httpx.MockTransport(handler)
        )
        try:
            if duplicate:
                with pytest.raises(CurriculumError, match="indexes"):
                    await provider.embed(["first", "second"])
            else:
                _, vectors = await provider.embed(["first", "second"])
                assert vectors[0][0] == 1.0 and vectors[1][0] == 2.0
        finally:
            await provider.close()

    asyncio.run(run())
