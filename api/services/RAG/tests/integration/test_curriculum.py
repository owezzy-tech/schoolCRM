import asyncio
import os
from dataclasses import asdict, replace
from hashlib import sha256
from pathlib import Path
from uuid import uuid4

import httpx
import pytest
from fastapi.testclient import TestClient

from adapters.controllers.curriculum_controller import get_curriculum
from adapters.embeddings.bge_m3 import BGEM3Embeddings
from adapters.parsers.curriculum_pdf import parse_curriculum_pdf
from adapters.repositories.curriculum_postgres import CurriculumPostgres
from domain.entities.curriculum import CurriculumError, CurriculumScope
from infrastructure.app import build_app
from infrastructure.auth import AuthContext, get_auth_context
from infrastructure.school_access import SchoolAccessClient
from use_cases.curriculum import CurriculumService

DSN = os.environ.get("CURRICULUM_TEST_DATABASE_URL")
pytestmark = pytest.mark.skipif(
    not DSN, reason="Set CURRICULUM_TEST_DATABASE_URL for pgvector tests"
)


def pdf(
    text: str = "Grade one English listening and speaking pronunciation vocabulary school",
) -> bytes:
    # A complete one-page PDF, with explicit font and xref offsets.
    stream = f"BT /F1 12 Tf 40 700 Td ({text}) Tj ET".encode()
    objects = [
        b"<< /Type /Catalog /Pages 2 0 R >>",
        b"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
        b"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] "
        b"/Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>",
        b"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
        f"<< /Length {len(stream)} >>\nstream\n".encode() + stream + b"\nendstream",
    ]
    result = b"%PDF-1.4\n"
    offsets = []
    for index, obj in enumerate(objects, 1):
        offsets.append(len(result))
        result += f"{index} 0 obj\n".encode() + obj + b"\nendobj\n"
    xref = len(result)
    result += b"xref\n0 6\n0000000000 65535 f \n"
    result += b"".join(f"{offset:010} 00000 n \n".encode() for offset in offsets)
    return result + f"trailer\n<< /Size 6 /Root 1 0 R >>\nstartxref\n{xref}\n%%EOF".encode()


class TestEmbeddings:
    identity = "test:bge-m3:1024"

    async def embed(self, texts: list[str]) -> tuple[str, list[list[float]]]:
        return self.identity, [
            ([0.0, 1.0] if "unsupported" in text else [1.0, 0.0]) + [0.0] * 1022 for text in texts
        ]

    async def close(self) -> None:
        pass


def service(target: CurriculumScope, store: CurriculumPostgres | None = None) -> CurriculumService:
    def handler(request: httpx.Request) -> httpx.Response:
        manager = request.headers["Authorization"] == "Bearer manager"
        if request.url.path.endswith("/memberships"):
            return httpx.Response(200 if manager else 403, json={"data": []})
        if request.url.path.endswith("/departments"):
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
        return httpx.Response(
            200,
            json={
                "data": [
                    {
                        "attributes": {
                            "schoolID": str(target.school_id),
                            "departmentID": str(target.department_id),
                            "capability": "teach",
                            "active": True,
                        }
                    }
                ]
            },
        )

    return CurriculumService(
        store or CurriculumPostgres(DSN),
        TestEmbeddings(),
        SchoolAccessClient("http://school", httpx.MockTransport(handler)),
        parse_curriculum_pdf,
    )


def target_scope() -> CurriculumScope:
    return CurriculumScope(uuid4(), uuid4(), "kenya-cbc", "grade-1", "english", "2024")


async def upload(svc: CurriculumService, target: CurriculumScope, content: bytes):
    return await svc.ingest(
        "teacher",
        "teacher-user",
        target,
        "English",
        "School vocabulary",
        "KICD",
        "https://kicd.ac.ke/original.pdf",
        "original.pdf",
        1,
        1,
        sha256(content).hexdigest(),
        content,
    )


def test_review_scoped_search_restart_withdrawal_and_immutable_original() -> None:
    async def run() -> None:
        target = target_scope()
        svc = service(target)
        await svc.store.migrate()
        content = pdf()
        try:
            source = await upload(svc, target, content)
            assert await svc.search("teacher", target, "pronunciation", 5) == []
            with pytest.raises(CurriculumError, match="management"):
                await svc.review("teacher", "teacher-user", source.id, "pending", "approved", "ok")
            with pytest.raises(CurriculumError, match="already exist"):
                await upload(svc, target, content)
            approved = await svc.review(
                "manager",
                "manager-user",
                source.id,
                "pending",
                "approved",
                "Original authority, revision and page verified",
            )
            assert approved.reviewed_by == "manager-user"
            evidence = await svc.search("teacher", target, "pronunciation", 5)
            assert evidence[0].source.sha256 == sha256(content).hexdigest()
            assert evidence[0].passage.page == 1
            assert evidence[0].passage.text.startswith("Grade one English")
            assert await svc.search("teacher", target, "unsupported", 5) == []
            for changed in [
                replace(target, framework="cambridge"),
                replace(target, stage="grade-2"),
                replace(target, revision="2025"),
                replace(target, subject="mathematics"),
            ]:
                assert await svc.search("teacher", changed, "pronunciation", 5) == []
            with pytest.raises(CurriculumError):
                await svc.search("teacher", replace(target, school_id=uuid4()), "pronunciation", 5)
            # A fresh service/store instance proves no in-memory index is involved.
            restarted = service(target)
            try:
                assert await restarted.search("teacher", target, "pronunciation", 5)
                assert await restarted.store.original(source.id) == content
                restarted.embeddings.identity = "test:other-provider"
                assert await restarted.search("teacher", target, "pronunciation", 5) == []
            finally:
                await restarted.close()
            await svc.review(
                "manager", "manager-user", source.id, "approved", "withdrawn", "replaced"
            )
            assert await svc.search("teacher", target, "pronunciation", 5) == []
            assert await svc.store.original(source.id) == content
            assert (await svc.source("teacher", source.id)).status == "withdrawn"
        finally:
            await svc.close()

    asyncio.run(run())


def test_concurrent_review_commits_once_and_audit_failure_rolls_back_index() -> None:
    async def run() -> None:
        target = target_scope()
        svc = service(target)
        await svc.store.migrate()
        try:
            source = await upload(svc, target, pdf())
            results = await asyncio.gather(
                *[
                    svc.review("manager", "manager", source.id, "pending", "approved", "verified")
                    for _ in range(4)
                ],
                return_exceptions=True,
            )
            assert sum(not isinstance(r, Exception) for r in results) == 1
            assert all(
                not isinstance(r, Exception) or isinstance(r, CurriculumError) for r in results
            )
            second = await upload(svc, replace(target, revision="2025"), pdf())
            original_audit = svc.store._audit

            async def failed_audit(*args) -> None:
                raise CurriculumError(503, "Injected audit failure")

            svc.store._audit = failed_audit
            with pytest.raises(CurriculumError, match="audit failure"):
                await svc.review("manager", "manager", second.id, "pending", "approved", "verified")
            svc.store._audit = original_audit
            assert (await svc.store.get(second.id)).status == "pending"
            assert (
                await svc.store.search(second.scope, "test:bge-m3:1024", [1.0] * 1024, 5, 0) == []
            )
        finally:
            await svc.close()

    asyncio.run(run())


def test_http_upload_search_and_error_contract() -> None:
    target = target_scope()
    svc = service(target)
    asyncio.run(svc.store.migrate())
    app = build_app()
    app.dependency_overrides[get_curriculum] = lambda: svc
    app.dependency_overrides[get_auth_context] = lambda: AuthContext("teacher", "teacher", [])
    client = TestClient(app)
    metadata = {
        **asdict(target),
        "title": "English",
        "topic": "School",
        "authority": "KICD",
        "source_url": "https://kicd.ac.ke/original.pdf",
        "sha256": sha256(pdf()).hexdigest(),
        "first_page": 1,
        "last_page": 1,
    }
    import json

    response = client.post(
        "/v1/rag/curriculum/sources",
        data={"metadata": json.dumps(metadata, default=str)},
        files={"file": ("original.pdf", pdf(), "application/pdf")},
    )
    assert response.status_code == 201, response.text
    assert response.json()["data"]["attributes"]["status"] == "pending"
    assert "id" not in response.json()["data"]["attributes"]
    response = client.post(
        "/v1/rag/curriculum/search",
        json={
            **json.loads(json.dumps(asdict(target), default=str)),
            "question": "pronunciation",
        },
    )
    assert response.status_code == 200
    assert response.json()["meta"]["abstained"] is True
    metadata["sha256"] = "0" * 64
    response = client.post(
        "/v1/rag/curriculum/sources",
        data={"metadata": json.dumps(metadata, default=str)},
        files={"file": ("original.pdf", pdf(), "application/pdf")},
    )
    assert response.status_code == 422
    assert "checksum" in response.json()["errors"][0]["detail"]
    response = client.post(
        "/v1/rag/curriculum/sources",
        data={"metadata": json.dumps(metadata, default=str)},
        files={"file": ("oversized.pdf", b"x" * (10 * 1024 * 1024 + 1), "application/pdf")},
    )
    assert response.status_code == 413
    invalid = client.post("/v1/rag/curriculum/search", json={"question": "missing scope"})
    assert invalid.status_code == 422 and "errors" in invalid.json()
    asyncio.run(svc.close())


@pytest.mark.skipif(not os.environ.get("CURRICULUM_LIVE_PDF"), reason="Set CURRICULUM_LIVE_PDF")
def test_real_kicd_pdf_and_local_bge_embeddings() -> None:
    async def run() -> None:
        target = target_scope()
        svc = service(target)
        await svc.embeddings.close()
        svc.embeddings = BGEM3Embeddings("ollama", "http://localhost:11434")
        await svc.store.migrate()
        content = Path(os.environ["CURRICULUM_LIVE_PDF"]).read_bytes()
        try:
            source = await svc.ingest(
                "teacher",
                "teacher",
                target,
                "English Grades 1-3 revised 2024",
                "School vocabulary",
                "Kenya Institute of Curriculum Development",
                "https://kicd.ac.ke/wp-content/uploads/"
                "2025/07/Grade-1-3-English-Activities-Revised-Sept.pdf",
                "english.pdf",
                17,
                17,
                "eeea001277e1d952f4a6ebd4e42a21072cbb1e261d2e7228e50fe699ef041d0d",
                content,
            )
            await svc.review(
                "manager",
                "test-reviewer",
                source.id,
                "pending",
                "approved",
                "Test fixture approval only; Grade 1 school vocabulary on PDF viewer page 17",
            )
            evidence = await svc.search(
                "teacher", target, "school classroom pronunciation vocabulary", 5
            )
            assert evidence and all(e.passage.page == 17 for e in evidence)
            assert all(e.model_identity.startswith("ollama:bge-m3:") for e in evidence)
            print(
                f"Real KICD retrieval: {len(evidence)} citations, page 17, "
                f"score={evidence[0].score:.3f}"
            )
            assert (
                await svc.search("teacher", target, "quantum chromodynamics gluon scattering", 5)
                == []
            )
        finally:
            await svc.close()

    asyncio.run(run())
