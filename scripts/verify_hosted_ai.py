"""Opt-in hosted staging proof using synthetic identities and an official source.

Credentials are environment-only. Output contains resource IDs and verdicts,
never tokens, credentials, model prompts or generated classroom content.
"""

import hashlib
import json
import os
from pathlib import Path
from uuid import UUID

import httpx

BASE = "https://web-admin-staging-c78a.up.railway.app"
SOURCE_SHA = "eeea001277e1d952f4a6ebd4e42a21072cbb1e261d2e7228e50fe699ef041d0d"
TEACHER_EMAIL = "schoolai-teacher@staging.invalid"
REQUEST_ID = "b99268af-c194-47c1-80aa-07b731d2b815"


def resource(value):
    return {**value["attributes"], "id": value["id"]}


def verify():
    if os.environ.get("SCHOOLCRM_HOSTED_PROOF") != "staging":
        raise RuntimeError("Explicit staging proof opt-in required")
    pdf = Path(os.environ["SCHOOLCRM_HOSTED_PDF"]).read_bytes()
    if hashlib.sha256(pdf).hexdigest() != SOURCE_SHA:
        raise RuntimeError("Official reviewed source checksum mismatch")
    password = os.environ["SCHOOLCRM_STAGING_QA_PASSWORD"]
    with httpx.Client(base_url=BASE, timeout=180) as client:

        def call(method, path, token=None, expected=(200, 201, 204), **kwargs):
            headers = {"Authorization": f"Bearer {token}"} if token else {}
            response = client.request(method, path, headers=headers, **kwargs)
            if response.status_code not in expected:
                # Do not render request/response bodies, which could contain credentials.
                raise RuntimeError(
                    f"Hosted proof failed at {method} {path}: HTTP {response.status_code}"
                )
            if response.status_code == 204:
                return None
            return response.json()

        def login(email, secret):
            return call(
                "POST", "/v1/auth/login", json={"email": email, "password": secret}
            )["data"]["attributes"]

        manager = login(
            "owezzy@owenadirah.com", os.environ["SCHOOLCRM_STAGING_ADMIN_PASSWORD"]
        )
        admin_token = manager["accessToken"]
        users = call("GET", "/v1/users?rows=100", admin_token)["data"]
        teacher = next(
            (resource(u) for u in users if u["attributes"]["email"] == TEACHER_EMAIL),
            None,
        )
        if teacher is None:
            teacher = resource(
                call(
                    "POST",
                    "/v1/users",
                    admin_token,
                    json={
                        "name": "Staging AI Teacher",
                        "email": TEACHER_EMAIL,
                        "roles": ["TEACHER"],
                        "department": "Languages",
                        "password": password,
                        "passwordConfirm": password,
                    },
                )["data"]
            )
        teacher_session = login(TEACHER_EMAIL, password)
        teacher_token = teacher_session["accessToken"]
        schools = call("GET", "/v1/schools", admin_token)["data"]
        school = next(
            (
                resource(s)
                for s in schools
                if s["attributes"]["name"] == "School AI staging validation"
            ),
            None,
        )
        if school is None:
            school = resource(
                call(
                    "POST",
                    "/v1/schools",
                    admin_token,
                    json={"name": "School AI staging validation"},
                )["data"]
            )
        school_path = f"/v1/schools/{school['id']}"
        departments = call("GET", school_path + "/departments", admin_token)["data"]
        department = next(
            (
                resource(d)
                for d in departments
                if d["attributes"]["name"] == "Languages"
            ),
            None,
        )
        if department is None:
            department = resource(
                call(
                    "POST",
                    school_path + "/departments",
                    admin_token,
                    json={"name": "Languages"},
                )["data"]
            )
        memberships = call("GET", school_path + "/memberships", admin_token)["data"]
        if not any(
            m["attributes"]["userID"] == teacher["id"]
            and m["attributes"].get("departmentID") == department["id"]
            and m["attributes"]["capability"] == "teach"
            and m["attributes"]["active"]
            for m in memberships
        ):
            call(
                "POST",
                school_path + "/memberships",
                admin_token,
                json={
                    "userID": teacher["id"],
                    "departmentID": department["id"],
                    "capability": "teach",
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
        print(
            json.dumps(
                {
                    "phase": "scope",
                    "school": school["id"],
                    "department": department["id"],
                    "teacher": teacher["id"],
                }
            ),
            flush=True,
        )
        source_id = os.environ.get("SCHOOLCRM_HOSTED_SOURCE_ID")
        if not source_id:
            metadata = {
                **scope,
                "title": "English Grades 1-3",
                "topic": "School",
                "authority": "KICD",
                "source_url": "https://kicd.ac.ke/wp-content/uploads/2025/07/Grade-1-3-English-Activities-Revised-Sept.pdf",
                "sha256": SOURCE_SHA,
                "first_page": 17,
                "last_page": 17,
            }
            source_id = call(
                "POST",
                "/v1/rag/curriculum/sources",
                teacher_token,
                data={"metadata": json.dumps(metadata)},
                files={"file": ("english.pdf", pdf, "application/pdf")},
            )["data"]["id"]
        UUID(source_id)
        print(json.dumps({"phase": "source", "source": source_id}), flush=True)
        source = call("GET", f"/v1/rag/curriculum/sources/{source_id}", teacher_token)[
            "data"
        ]["attributes"]
        if source["status"] == "pending":
            call(
                "POST",
                f"/v1/rag/curriculum/sources/{source_id}/review",
                admin_token,
                json={
                    "expected_status": "pending",
                    "decision": "approved",
                    "note": "Staging proof: reviewed official KICD English page 17",
                },
            )
        source = call("GET", f"/v1/rag/curriculum/sources/{source_id}", teacher_token)[
            "data"
        ]["attributes"]
        if (
            source["status"] != "approved"
            or source["model_identity"] != "cloudflare:@cf/baai/bge-m3:1024"
        ):
            raise RuntimeError("Approved Cloudflare index identity mismatch")
        evidence = call(
            "POST",
            "/v1/rag/curriculum/search",
            teacher_token,
            json={
                **scope,
                "question": "school classroom pronunciation vocabulary",
                "limit": 5,
            },
        )
        if evidence["meta"]["abstained"] or not evidence["data"]:
            raise RuntimeError("Approved source retrieval abstained")
        print(
            json.dumps({"phase": "retrieval", "evidence_count": len(evidence["data"])}),
            flush=True,
        )
        request = {
            **scope,
            "request_id": REQUEST_ID,
            "topic": "school classroom pronunciation vocabulary",
            "duration_minutes": 30,
        }
        result = call("POST", "/v1/rag/lessons/generate", teacher_token, json=request)[
            "data"
        ]["attributes"]
        plan_id = result["planID"]
        before_retry = call("GET", f"/v1/lessons/{plan_id}/versions", teacher_token)["data"]
        repeat = call("POST", "/v1/rag/lessons/generate", teacher_token, json=request)[
            "data"
        ]["attributes"]
        if repeat["planID"] != plan_id:
            raise RuntimeError("Generation retry changed persisted plan")
        versions = call("GET", f"/v1/lessons/{plan_id}/versions", teacher_token)["data"]
        if [v["id"] for v in before_retry] != [v["id"] for v in versions]:
            raise RuntimeError("Generation retry changed immutable version history")
        generated = next((v for v in versions if v["attributes"]["version"] == 1), None)
        if generated is None:
            raise RuntimeError("Original generated version missing")
        saved = generated["attributes"]["content"]
        if saved["generation"]["modelVersion"] != "DeepSeek-V4.1-Flash":
            raise RuntimeError(
                "Persisted generation identity or version count mismatch"
            )
        if (
            saved["durationMinutes"] != 30
            or sum(a["minutes"] for a in saved["activities"]) != 30
        ):
            raise RuntimeError("Generated lesson duration invariant failed")
        if not saved["citations"] or not all(
            c["sourceID"] == source_id
            and c["embeddingModel"] == "cloudflare:@cf/baai/bge-m3:1024"
            for c in saved["citations"]
        ):
            raise RuntimeError("Persisted citations do not identify approved evidence")
        status = call(
            "GET", f"/v1/rag/lessons/generations/{REQUEST_ID}", teacher_token
        )["data"]["attributes"]
        if status["status"] != "completed":
            raise RuntimeError("Durable generation receipt not completed")
        print(
            json.dumps(
                {
                    "phase": "complete",
                    "plan": plan_id,
                    "source": source_id,
                    "request": REQUEST_ID,
                    "model": saved["generation"]["modelVersion"],
                    "citations": len(saved["citations"]),
                    "versions": len(versions),
                    "replay": "same plan",
                }
            ),
            flush=True,
        )


if __name__ == "__main__":
    try:
        verify()
    except (RuntimeError, KeyError, ValueError, httpx.HTTPError, OSError) as exc:
        # Keep the public diagnostic bounded and never display transport/request objects.
        print(
            str(exc)
            if isinstance(exc, RuntimeError)
            else f"Hosted proof failed: {type(exc).__name__}",
            flush=True,
        )
        raise SystemExit(1) from None
