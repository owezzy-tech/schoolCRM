"""Opt-in HTTP journey against isolated Go + RAG services. PYTEST_DONT_REWRITE.

Disable assertion introspection so HTTP authentication headers cannot enter failure reports.
"""

import json
import os
from pathlib import Path
from uuid import uuid4

import httpx
import pytest

pytestmark = pytest.mark.skipif(
    not os.environ.get("CURRICULUM_LIVE_RAG_URL"),
    reason="Set CURRICULUM_LIVE_RAG_URL and isolated Go service URLs",
)


def test_real_go_authority_upload_review_search_revoke_and_withdraw() -> None:
    auth_url = os.environ["CURRICULUM_LIVE_AUTH_URL"]
    school_url = os.environ["CURRICULUM_LIVE_SCHOOL_URL"]
    rag_url = os.environ["CURRICULUM_LIVE_RAG_URL"]
    content = Path(os.environ["CURRICULUM_LIVE_PDF"]).read_bytes()
    with httpx.Client(timeout=120) as client:

        def login(email: str) -> dict:
            response = client.post(
                f"{auth_url}/v1/auth/login",
                json={
                    "email": email,
                    "password": "gophers",
                },
            )
            assert response.status_code == 200
            return response.json()["data"]["attributes"]

        manager = login("superadmin@example.com")
        teacher = login("teacher@example.com")
        manager_headers = {"Authorization": f"Bearer {manager['accessToken']}"}
        teacher_headers = {"Authorization": f"Bearer {teacher['accessToken']}"}

        def create(path: str, data: dict) -> dict:
            response = client.post(f"{school_url}{path}", json=data, headers=manager_headers)
            assert response.status_code in (200, 201), response.text
            resource = response.json()["data"]
            return {**resource["attributes"], "id": resource["id"]}

        school = create("/v1/schools", {"name": f"Curriculum proof {uuid4()}"})
        school_id = school["id"]
        department = create(f"/v1/schools/{school_id}/departments", {"name": "Languages"})
        membership = create(
            f"/v1/schools/{school_id}/memberships",
            {
                "userID": teacher["user"]["id"],
                "departmentID": department["id"],
                "capability": "teach",
            },
        )
        scope = {
            "school_id": school_id,
            "department_id": department["id"],
            "framework": "kenya-cbc",
            "stage": "grade-1",
            "subject": "english",
            "revision": "2024",
        }
        metadata = {
            **scope,
            "title": "English Grades 1-3 revised 2024",
            "topic": "School",
            "authority": "Kenya Institute of Curriculum Development",
            "source_url": "https://kicd.ac.ke/wp-content/uploads/2025/07/"
            "Grade-1-3-English-Activities-Revised-Sept.pdf",
            "sha256": "eeea001277e1d952f4a6ebd4e42a21072cbb1e261d2e7228e50fe699ef041d0d",
            "first_page": 17,
            "last_page": 17,
        }
        response = client.post(
            f"{rag_url}/v1/rag/curriculum/sources",
            data={"metadata": json.dumps(metadata)},
            headers=teacher_headers,
            files={"file": ("english.pdf", content, "application/pdf")},
        )
        assert response.status_code == 201, response.text
        assert response.headers["Content-Type"].startswith("application/vnd.api+json")
        source_id = response.json()["data"]["id"]
        query = {**scope, "question": "school classroom pronunciation vocabulary"}

        def search(headers: dict, request: dict = query) -> httpx.Response:
            return client.post(f"{rag_url}/v1/rag/curriculum/search", json=request, headers=headers)

        assert search(teacher_headers).json()["meta"]["abstained"] is True
        review_url = f"{rag_url}/v1/rag/curriculum/sources/{source_id}/review"
        approval = {
            "expected_status": "pending",
            "decision": "approved",
            "note": "Isolated test review: official source, revised 2024, Grade 1 page 17",
        }
        assert client.post(review_url, json=approval, headers=teacher_headers).status_code == 403
        response = client.post(review_url, json=approval, headers=manager_headers)
        assert response.status_code == 200, response.text
        result = search(teacher_headers)
        assert result.status_code == 200, result.text
        assert result.json()["data"]
        assert all(row["attributes"]["passage"]["page"] == 17 for row in result.json()["data"])
        assert search(teacher_headers, {**query, "revision": "2025"}).json()["data"] == []
        assert search(teacher_headers, {**query, "school_id": str(uuid4())}).status_code == 403
        assert (
            search(teacher_headers, {**query, "question": "quantum chromodynamics"}).json()["data"]
            == []
        )
        assert search({}).status_code == 401
        # Revoke the Go membership while keeping the same bearer token.
        response = client.delete(
            f"{school_url}/v1/schools/{school_id}/memberships/{membership['id']}",
            headers=manager_headers,
        )
        assert response.status_code == 204
        denied = search(teacher_headers)
        assert denied.status_code == 403, denied.text
        assert "capability" in denied.json()["errors"][0]["detail"]
        assert (
            client.post(
                review_url,
                json={
                    "expected_status": "approved",
                    "decision": "withdrawn",
                    "note": "End of isolated proof",
                },
                headers=manager_headers,
            ).status_code
            == 200
        )
        assert search(manager_headers).json()["data"] == []
        original = client.get(
            f"{rag_url}/v1/rag/curriculum/sources/{source_id}/original", headers=manager_headers
        )
        assert original.status_code == 200 and original.content == content
        print("Live Go + RAG journey passed: upload, approval, citations, revocation, withdrawal")
