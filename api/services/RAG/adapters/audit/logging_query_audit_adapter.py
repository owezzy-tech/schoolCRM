import logging

from domain.entities.audit import QueryAuditRecord
from domain.ports.query_audit import IQueryAuditPort

logger = logging.getLogger(__name__)


class LoggingQueryAuditAdapter(IQueryAuditPort):
    async def record(self, entry: QueryAuditRecord) -> None:
        logger.info(
            "rag_admissions_query",
            extra={
                "query_id": entry.query_id,
                "subject": entry.subject,
                "scope": entry.scope,
                "collection": entry.collection,
                "application_id": entry.application_id,
                "question_hash": entry.question_hash,
                "citation_count": entry.citation_count,
                "created_at": entry.created_at.isoformat(),
            },
        )
