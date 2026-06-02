from dataclasses import dataclass
from datetime import datetime


@dataclass(frozen=True, slots=True)
class QueryAuditRecord:
    query_id: str
    subject: str
    scope: str
    collection: str
    application_id: str | None
    question_hash: str
    citation_count: int
    created_at: datetime
