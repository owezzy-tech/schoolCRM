from dataclasses import dataclass
from enum import StrEnum


class AdmissionsScope(StrEnum):
    STAFF = "staff"
    APPLICANT_OWN = "applicant_own"
    REVIEWER_ASSIGNED = "reviewer_assigned"


@dataclass(frozen=True, slots=True)
class AdmissionsQueryContext:
    scope: AdmissionsScope
    collection_name: str
    subject: str
    application_id: str | None = None
