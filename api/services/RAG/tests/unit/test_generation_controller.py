from uuid import uuid4

from fastapi.testclient import TestClient

from adapters.controllers.lesson_generation_controller import get_workflow
from infrastructure.app import build_app
from infrastructure.auth import AuthContext, get_auth_context


class WorkflowFixture:
    def __init__(self):
        self.calls = []

    async def execute(self, request, context):
        self.calls.append((request, context))
        return {"planID": str(uuid4()), "version": 1, "creationStatus": "draft"}


def test_generation_uses_authenticated_actor_and_strict_scope_contract():
    app = build_app()
    actor, school, department, request_id = uuid4(), uuid4(), uuid4(), uuid4()
    workflow = WorkflowFixture()
    app.dependency_overrides[get_auth_context] = lambda: AuthContext(
        str(actor), "private-token", []
    )
    app.dependency_overrides[get_workflow] = lambda: workflow
    client = TestClient(app)
    body = {
        "request_id": str(request_id),
        "school_id": str(school),
        "department_id": str(department),
        "framework": "kenya-cbc",
        "stage": "grade-1",
        "subject": "English",
        "revision": "2024",
        "topic": "School vocabulary",
        "duration_minutes": 30,
    }
    response = client.post("/v1/rag/lessons/generate", json=body)
    assert response.status_code == 200
    assert response.headers["content-type"].startswith("application/vnd.api+json")
    assert response.json()["data"]["id"] == str(request_id)
    request, context = workflow.calls[0]
    assert context.actor_id == actor and context.token == "private-token"
    assert "private-token" not in repr(context)
    assert request.scope.subject == "english" and request.duration_minutes == 30
    invalid = client.post("/v1/rag/lessons/generate", json={**body, "duration_minutes": 30.5})
    assert invalid.status_code == 422 and "errors" in invalid.json()
    invalid = client.post("/v1/rag/lessons/generate", json={**body, "actor_id": str(uuid4())})
    assert invalid.status_code == 422 and len(workflow.calls) == 1
    app.dependency_overrides[get_auth_context] = lambda: AuthContext("anonymous", None, [])
    denied = client.post("/v1/rag/lessons/generate", json=body)
    assert denied.status_code == 401 and len(workflow.calls) == 1
