from dataclasses import dataclass

from fastapi.testclient import TestClient

from domain.entities.query import CitationSource, QueryResult
from infrastructure.auth_client import AuthenticatedUser
from infrastructure.dependencies import (
    get_admissions_query_use_case,
    get_application_ownership_checker,
    get_auth_service_client,
)
from main import app
from use_cases.admissions_query import AdmissionsQueryCommand


class FakeAuthServiceClient:
    async def authenticate(self, bearer_token: str) -> AuthenticatedUser:
        users = {
            "staff-token": AuthenticatedUser(
                user_id="staff-1",
                subject="staff-1",
                roles=["ADMISSIONS_ADMIN"],
                claims={"sub": "staff-1", "roles": ["ADMISSIONS_ADMIN"]},
            ),
            "applicant-token": AuthenticatedUser(
                user_id="applicant-1",
                subject="applicant-1",
                roles=["APPLICANT"],
                claims={"sub": "applicant-1", "roles": ["APPLICANT"]},
            ),
            "reviewer-token": AuthenticatedUser(
                user_id="reviewer-1",
                subject="reviewer-1",
                roles=["APPLICATION_REVIEWER"],
                claims={"sub": "reviewer-1", "roles": ["APPLICATION_REVIEWER"]},
            ),
        }
        return users[bearer_token]


class FakeOwnershipChecker:
    async def is_owner(self, subject: str, application_id: str) -> bool:
        return subject == "applicant-1" and application_id == "app-1"

    async def is_assigned_reviewer(self, subject: str, application_id: str) -> bool:
        return subject == "reviewer-1" and application_id == "app-1"


@dataclass(slots=True)
class FakeAdmissionsQueryUseCase:
    command: AdmissionsQueryCommand | None = None

    async def execute(self, command: AdmissionsQueryCommand) -> QueryResult:
        self.command = command
        return QueryResult(
            answer="Use the admissions checklist before submitting.",
            document_ids=["policy-1"],
            snippets=["Checklist guidance"],
            citations=[
                CitationSource(
                    document_id="policy-1",
                    chunk_index=0,
                    excerpt="Checklist guidance",
                    collection=command.context.collection_name,
                    score=0.88,
                    source_type="Policy",
                    title="Checklist",
                    path=["AdmissionsPolicy", "Checklist"],
                )
            ],
        )


def test_staff_policy_query_returns_jsonapi_result() -> None:
    use_case = FakeAdmissionsQueryUseCase()
    app.dependency_overrides[get_auth_service_client] = lambda: FakeAuthServiceClient()
    app.dependency_overrides[get_admissions_query_use_case] = lambda: use_case
    client = TestClient(app)

    try:
        response = client.post(
            "/v1/rag/admissions/staff/query",
            json={"question": "What is the checklist?", "top_k": 2},
            headers={"Authorization": "Bearer staff-token"},
        )
    finally:
        app.dependency_overrides.clear()

    assert response.status_code == 200
    payload = response.json()
    assert payload["jsonapi"] == {"version": "1.1"}
    attributes = payload["data"]["attributes"]
    assert attributes["answer"] == "Use the admissions checklist before submitting."
    assert attributes["citations"] == [
        {
            "documentId": "policy-1",
            "chunkIndex": 0,
            "excerpt": "Checklist guidance",
            "collection": "admissions-policy",
            "score": 0.88,
            "sourceType": "Policy",
            "title": "Checklist",
            "path": ["AdmissionsPolicy", "Checklist"],
        }
    ]
    assert use_case.command is not None
    assert use_case.command.context.collection_name == "admissions-policy"
    assert use_case.command.top_k == 2


def test_applicant_query_is_limited_to_owned_application() -> None:
    use_case = FakeAdmissionsQueryUseCase()
    app.dependency_overrides[get_auth_service_client] = lambda: FakeAuthServiceClient()
    app.dependency_overrides[get_application_ownership_checker] = lambda: FakeOwnershipChecker()
    app.dependency_overrides[get_admissions_query_use_case] = lambda: use_case
    client = TestClient(app)

    try:
        response = client.post(
            "/v1/rag/admissions/applicant/applications/app-1/query",
            json={"question": "What remains?"},
            headers={"Authorization": "Bearer applicant-token"},
        )
    finally:
        app.dependency_overrides.clear()

    assert response.status_code == 200
    assert use_case.command is not None
    assert use_case.command.context.collection_name == "application-app-1"
    assert use_case.command.context.application_id == "app-1"


def test_applicant_query_rejects_unowned_application() -> None:
    app.dependency_overrides[get_auth_service_client] = lambda: FakeAuthServiceClient()
    app.dependency_overrides[get_application_ownership_checker] = lambda: FakeOwnershipChecker()
    app.dependency_overrides[get_admissions_query_use_case] = lambda: FakeAdmissionsQueryUseCase()
    client = TestClient(app)

    try:
        response = client.post(
            "/v1/rag/admissions/applicant/applications/app-2/query",
            json={"question": "What remains?"},
            headers={"Authorization": "Bearer applicant-token"},
        )
    finally:
        app.dependency_overrides.clear()

    assert response.status_code == 403
    assert response.json()["detail"] == "Not authorized to query this application"


def test_reviewer_query_is_limited_to_assigned_application() -> None:
    use_case = FakeAdmissionsQueryUseCase()
    app.dependency_overrides[get_auth_service_client] = lambda: FakeAuthServiceClient()
    app.dependency_overrides[get_application_ownership_checker] = lambda: FakeOwnershipChecker()
    app.dependency_overrides[get_admissions_query_use_case] = lambda: use_case
    client = TestClient(app)

    try:
        response = client.post(
            "/v1/rag/admissions/reviewer/applications/app-1/query",
            json={"question": "Summarize documents."},
            headers={"Authorization": "Bearer reviewer-token"},
        )
    finally:
        app.dependency_overrides.clear()

    assert response.status_code == 200
    assert use_case.command is not None
    assert use_case.command.context.collection_name == "application-app-1"
    assert use_case.command.context.application_id == "app-1"


def test_reviewer_query_rejects_unassigned_application() -> None:
    app.dependency_overrides[get_auth_service_client] = lambda: FakeAuthServiceClient()
    app.dependency_overrides[get_application_ownership_checker] = lambda: FakeOwnershipChecker()
    app.dependency_overrides[get_admissions_query_use_case] = lambda: FakeAdmissionsQueryUseCase()
    client = TestClient(app)

    try:
        response = client.post(
            "/v1/rag/admissions/reviewer/applications/app-2/query",
            json={"question": "Summarize documents."},
            headers={"Authorization": "Bearer reviewer-token"},
        )
    finally:
        app.dependency_overrides.clear()

    assert response.status_code == 403
    assert response.json()["detail"] == "Not assigned as reviewer for this application"
