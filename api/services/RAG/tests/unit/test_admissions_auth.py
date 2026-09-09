import asyncio

from fastapi import HTTPException

from domain.entities.admissions_query_context import AdmissionsScope
from infrastructure.admissions_auth import (
    guard_applicant_query,
    guard_reviewer_query,
    guard_staff_query,
)
from infrastructure.auth import AuthContext


class FakeOwnershipChecker:
    def __init__(self, *, owns: bool = False, assigned: bool = False) -> None:
        self.owns = owns
        self.assigned = assigned

    async def is_owner(self, subject: str, application_id: str) -> bool:
        return self.owns

    async def is_assigned_reviewer(self, subject: str, application_id: str) -> bool:
        return self.assigned


def test_guard_staff_query_accepts_admissions_staff() -> None:
    context = guard_staff_query(
        AuthContext(subject="staff-1", token="token", roles=["ADMISSIONS_ADMIN"])
    )

    assert context.scope == AdmissionsScope.STAFF
    assert context.collection_name == "admissions-policy"
    assert context.subject == "staff-1"


def test_guard_staff_query_rejects_applicant() -> None:
    try:
        guard_staff_query(AuthContext(subject="applicant-1", token="token", roles=["APPLICANT"]))
    except HTTPException as exc:
        assert exc.status_code == 403
        assert exc.detail == "Insufficient role for staff RAG query"
    else:
        raise AssertionError("expected HTTPException")


def test_guard_applicant_query_requires_own_application() -> None:
    context = asyncio.run(
        guard_applicant_query(
            AuthContext(subject="applicant-1", token="token", roles=["APPLICANT"]),
            "app-1",
            FakeOwnershipChecker(owns=True),
        )
    )

    assert context.scope == AdmissionsScope.APPLICANT_OWN
    assert context.collection_name == "application-app-1"
    assert context.application_id == "app-1"


def test_guard_applicant_query_rejects_other_application() -> None:
    try:
        asyncio.run(
            guard_applicant_query(
                AuthContext(subject="applicant-1", token="token", roles=["APPLICANT"]),
                "app-2",
                FakeOwnershipChecker(owns=False),
            )
        )
    except HTTPException as exc:
        assert exc.status_code == 403
        assert exc.detail == "Not authorized to query this application"
    else:
        raise AssertionError("expected HTTPException")


def test_guard_reviewer_query_requires_assignment() -> None:
    context = asyncio.run(
        guard_reviewer_query(
            AuthContext(subject="reviewer-1", token="token", roles=["APPLICATION_REVIEWER"]),
            "app-1",
            FakeOwnershipChecker(assigned=True),
        )
    )

    assert context.scope == AdmissionsScope.REVIEWER_ASSIGNED
    assert context.collection_name == "application-app-1"
    assert context.application_id == "app-1"


def test_guard_reviewer_query_rejects_unassigned_reviewer() -> None:
    try:
        asyncio.run(
            guard_reviewer_query(
                AuthContext(subject="reviewer-1", token="token", roles=["APPLICATION_REVIEWER"]),
                "app-2",
                FakeOwnershipChecker(assigned=False),
            )
        )
    except HTTPException as exc:
        assert exc.status_code == 403
        assert exc.detail == "Not assigned as reviewer for this application"
    else:
        raise AssertionError("expected HTTPException")


def test_guard_reviewer_query_allows_admin_without_assignment() -> None:
    context = asyncio.run(
        guard_reviewer_query(
            AuthContext(subject="admin-1", token="token", roles=["ADMISSIONS_ADMIN"]),
            "app-1",
            FakeOwnershipChecker(assigned=False),
        )
    )

    assert context.scope == AdmissionsScope.REVIEWER_ASSIGNED
    assert context.subject == "admin-1"
