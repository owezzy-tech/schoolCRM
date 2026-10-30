"""Opt-in real Go/PostgreSQL/Ollama journey. PYTEST_DONT_REWRITE.

Credentials stay out of assertion introspection. The default DeepSeek transport
is a contract fixture; set DEEPSEEK_LIVE_TEST=1 to call the configured real API.
"""

import json
import os
from pathlib import Path
from uuid import uuid4

import httpx
import pytest
from fastapi.testclient import TestClient
from pydantic import SecretStr

from adapters.llm.deepseek_lesson_generator import DeepSeekLessonGenerator
from infrastructure.app import build_app
from infrastructure.config import get_settings

pytestmark = pytest.mark.skipif(
    not os.environ.get("GENERATION_LIVE_SERVICES"),
    reason="Set GENERATION_LIVE_SERVICES for isolated Go services",
)


def test_generation_persistence_reuse_and_private_drafts_through_real_go(monkeypatch):
    dsn = os.environ["CURRICULUM_TEST_DATABASE_URL"]
    pdf = Path(os.environ["CURRICULUM_LIVE_PDF"]).read_bytes()
    auth_url = os.environ.get("GENERATION_AUTH_URL", "http://127.0.0.1:6100")
    school_url = os.environ.get("GENERATION_SCHOOL_URL", "http://127.0.0.1:3100")
    monkeypatch.setenv("RAG_CURRICULUM_DATABASE_URL", dsn)
    monkeypatch.setenv("RAG_AUTH_SERVICE_URL", auth_url)
    monkeypatch.setenv("RAG_SCHOOL_SERVICE_URL", school_url)
    monkeypatch.setenv("RAG_OLLAMA_BASE_URL", "http://localhost:11434")
    get_settings.cache_clear()
    with httpx.Client(timeout=120) as go:

        def login(email):
            response = go.post(
                auth_url + "/v1/auth/login", json={"email": email, "password": "gophers"}
            )
            if response.status_code != 200:
                raise RuntimeError("Isolated fixture login failed")
            return response.json()["data"]["attributes"]

        manager, teacher, hod, dean = [
            login(email)
            for email in (
                "superadmin@example.com",
                "teacher@example.com",
                "admin@example.com",
                "user@example.com",
            )
        ]

        def command(actor, method, path, body=None):
            response = go.request(
                method,
                school_url + path,
                headers={"Authorization": f"Bearer {actor['accessToken']}"},
                json=body,
            )
            if response.status_code not in (200, 201, 204):
                raise RuntimeError(
                    f"Go proof command failed: {path}, status {response.status_code}"
                )
            if response.status_code == 204:
                return None
            resource = response.json()["data"]
            return {**resource["attributes"], "id": resource["id"]}

        school = command(manager, "POST", "/v1/schools", {"name": f"Generation proof {uuid4()}"})
        department = command(
            manager, "POST", f"/v1/schools/{school['id']}/departments", {"name": "Languages"}
        )
        for actor, capability in [
            (teacher, "teach"),
            (hod, "review_lessons"),
            (hod, "teach"),
            (dean, "approve_lessons"),
        ]:
            command(
                manager,
                "POST",
                f"/v1/schools/{school['id']}/memberships",
                {
                    "userID": actor["user"]["id"],
                    "departmentID": department["id"],
                    "capability": capability,
                },
            )
        scope = {
            "school_id": school["id"],
            "department_id": department["id"],
            "framework": "kenya-cbc",
            "stage": "grade-1",
            "subject": "english",
            "revision": "2024",
        }
        app = build_app()
        model_calls = []

        def model_contract(request):
            if request.url.path == "/models":
                return httpx.Response(
                    200, json={"data": [{"id": "deepseek-flash", "name": "DeepSeek-V4.1-Flash"}]}
                )
            model_calls.append(request.url.path)
            body = json.loads(request.content)
            prompt = json.loads(body["messages"][1]["content"])
            output = {
                "title": "Generated classroom vocabulary",
                "objectives": ["Name classroom objects"],
                "prerequisites": ["Listen to familiar words"],
                "materials": ["Picture cards"],
                "activities": [
                    {
                        "minutes": prompt["duration_minutes"],
                        "title": "Name and practise",
                        "detail": "Name a desk, chair and book",
                    }
                ],
                "assessment": "Name three classroom objects",
                "differentiation": "Use visual prompts",
                "citation_ids": [prompt["evidence"][0]["id"]],
            }
            return httpx.Response(
                200,
                json={
                    "id": "fixture-completion",
                    "model": "deepseek-flash",
                    "choices": [
                        {"finish_reason": "stop", "message": {"content": json.dumps(output)}}
                    ],
                },
            )

        with TestClient(app) as client:
            workflow = app.state.container.lesson_workflow
            if not os.environ.get("DEEPSEEK_LIVE_TEST"):
                client.portal.call(workflow.generator.close)
                workflow.generator = DeepSeekLessonGenerator(
                    SecretStr("fixture-key"), httpx.MockTransport(model_contract)
                )
            metadata = {
                **scope,
                "title": "English Grades 1-3",
                "topic": "School",
                "authority": "KICD",
                "source_url": "https://kicd.ac.ke/wp-content/uploads/2025/07/Grade-1-3-English-Activities-Revised-Sept.pdf",
                "sha256": "eeea001277e1d952f4a6ebd4e42a21072cbb1e261d2e7228e50fe699ef041d0d",
                "first_page": 17,
                "last_page": 17,
            }
            teacher_headers = {"Authorization": f"Bearer {teacher['accessToken']}"}
            response = client.post(
                "/v1/rag/curriculum/sources",
                headers=teacher_headers,
                data={"metadata": json.dumps(metadata)},
                files={"file": ("english.pdf", pdf, "application/pdf")},
            )
            assert response.status_code == 201
            source_id = response.json()["data"]["id"]
            response = client.post(
                f"/v1/rag/curriculum/sources/{source_id}/review",
                headers={"Authorization": f"Bearer {manager['accessToken']}"},
                json={
                    "expected_status": "pending",
                    "decision": "approved",
                    "note": "Isolated proof: reviewed official English page 17",
                },
            )
            assert response.status_code == 200
            request = {
                **scope,
                "request_id": str(uuid4()),
                "topic": "school classroom pronunciation vocabulary",
                "duration_minutes": 30,
            }
            first = client.post("/v1/rag/lessons/generate", headers=teacher_headers, json=request)
            if first.status_code != 200:
                details = [error.get("detail", "") for error in first.json().get("errors", [])]
                raise RuntimeError(f"Generation returned HTTP {first.status_code}: {details}")
            result = first.json()["data"]["attributes"]
            plan_id = result["planID"]
            repeat = client.post("/v1/rag/lessons/generate", headers=teacher_headers, json=request)
            assert (
                repeat.status_code == 200
                and repeat.json()["data"]["attributes"]["planID"] == plan_id
            )
            if not os.environ.get("DEEPSEEK_LIVE_TEST"):
                assert len(model_calls) == 1
            versions = go.get(
                f"{school_url}/v1/lessons/{plan_id}/versions", headers=teacher_headers
            ).json()["data"]
            assert len(versions) == 1
            saved = versions[0]["attributes"]["content"]
            assert saved["schemaVersion"] == 1 and saved["durationMinutes"] == 30
            assert sum(a["minutes"] for a in saved["activities"]) == 30
            assert saved["citations"][0]["sourceID"] == source_id
            assert saved["generation"]["modelVersion"] == "DeepSeek-V4.1-Flash"
            hidden = go.post(
                f"{school_url}/v1/lessons/{plan_id}/reuse",
                json={"version": 1, "requestID": str(uuid4())},
                headers={"Authorization": f"Bearer {hod['accessToken']}"},
            )
            assert hidden.status_code == 404
            version_path = f"/v1/lessons/{plan_id}/versions/1"
            command(teacher, "POST", version_path + "/submit")
            command(
                hod,
                "POST",
                version_path + "/review",
                {"decision": "approve", "feedback": "Evidence reviewed"},
            )
            command(
                dean,
                "POST",
                version_path + "/approval",
                {"decision": "approve", "feedback": "Approved"},
            )
            command(teacher, "POST", version_path + "/publish")
            reused = command(
                hod,
                "POST",
                f"/v1/lessons/{plan_id}/reuse",
                {"version": 1, "requestID": str(uuid4())},
            )
            assert reused["status"] == "draft" and reused["publishedVersion"] is None
            updated = {**saved, "objectives": ["Teacher-edited objective"]}
            command(
                teacher,
                "POST",
                f"/v1/lessons/{plan_id}/versions",
                {
                    "baseVersion": 1,
                    "title": "Edited generated plan",
                    "content": updated,
                    "changeSummary": "Human edit",
                },
            )
            plan = command(teacher, "GET", f"/v1/lessons/{plan_id}")
            assert plan["currentVersion"] == 2 and plan["publishedVersion"] == 1
            status = client.get(
                f"/v1/rag/lessons/generations/{request['request_id']}", headers=teacher_headers
            )
            assert (
                status.status_code == 200
                and status.json()["data"]["attributes"]["status"] == "completed"
            )
            print(
                json.dumps(
                    {
                        "generation_plan": plan_id,
                        "reused_plan": reused["id"],
                        "source": source_id,
                        "model_transport": "real DeepSeek"
                        if os.environ.get("DEEPSEEK_LIVE_TEST")
                        else "DeepSeek contract fixture",
                    }
                )
            )
    get_settings.cache_clear()
