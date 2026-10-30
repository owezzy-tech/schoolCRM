from fastapi import HTTPException, status

from domain.entities.admissions_query_context import AdmissionsQueryContext, AdmissionsScope
from domain.ports.application_ownership import IApplicationOwnershipChecker
from infrastructure.auth import AuthContext

ADMISSIONS_STAFF_ROLES = frozenset(
    {
        "ADMISSIONS_ADMIN",
        "RECRUITER",
        "APPLICATION_REVIEWER",
        "MARKETING_MANAGER",
        "EVENT_MANAGER",
        "REPORT_VIEWER",
    }
)
APPLICANT_ROLE = "APPLICANT"
REVIEWER_ROLE = "APPLICATION_REVIEWER"
ADMIN_ROLE = "ADMISSIONS_ADMIN"


def guard_staff_query(auth: AuthContext) -> AdmissionsQueryContext:
    if not ADMISSIONS_STAFF_ROLES.intersection(auth.roles):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Insufficient role for staff RAG query",
        )

    return AdmissionsQueryContext(
        scope=AdmissionsScope.STAFF,
        collection_name="admissions-policy",
        subject=auth.subject,
        application_id=None,
    )


async def guard_applicant_query(
    auth: AuthContext,
    application_id: str,
    ownership_checker: IApplicationOwnershipChecker,
) -> AdmissionsQueryContext:
    if auth.roles != [APPLICANT_ROLE] and APPLICANT_ROLE not in auth.roles:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Applicant RAG queries require applicant scope",
        )

    if not await ownership_checker.is_owner(auth.subject, application_id):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to query this application",
        )

    return AdmissionsQueryContext(
        scope=AdmissionsScope.APPLICANT_OWN,
        collection_name=f"application-{application_id}",
        subject=auth.subject,
        application_id=application_id,
    )


async def guard_reviewer_query(
    auth: AuthContext,
    application_id: str,
    ownership_checker: IApplicationOwnershipChecker,
) -> AdmissionsQueryContext:
    if ADMIN_ROLE not in auth.roles and REVIEWER_ROLE not in auth.roles:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Reviewer RAG queries require reviewer scope",
        )

    if ADMIN_ROLE not in auth.roles and not await ownership_checker.is_assigned_reviewer(
        auth.subject,
        application_id,
    ):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not assigned as reviewer for this application",
        )

    return AdmissionsQueryContext(
        scope=AdmissionsScope.REVIEWER_ASSIGNED,
        collection_name=f"application-{application_id}",
        subject=auth.subject,
        application_id=application_id,
    )
