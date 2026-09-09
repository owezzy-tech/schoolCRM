from dataclasses import dataclass
from datetime import datetime, timezone
from hashlib import sha256
from uuid import uuid4

from domain.entities.admissions_query_context import AdmissionsQueryContext
from domain.entities.audit import QueryAuditRecord
from domain.entities.query import CitationSource, QueryResult
from domain.ports.admissions_answer_provider import IAdmissionsAnswerProvider
from domain.ports.graph_retriever import IGraphRetriever
from domain.ports.query_audit import IQueryAuditPort


@dataclass(frozen=True, slots=True)
class AdmissionsQueryCommand:
    question: str
    context: AdmissionsQueryContext
    top_k: int = 3


class AdmissionsQueryUseCase:
    def __init__(
        self,
        *,
        graph_retriever: IGraphRetriever,
        answer_provider: IAdmissionsAnswerProvider,
        audit_port: IQueryAuditPort,
    ) -> None:
        self._graph_retriever = graph_retriever
        self._answer_provider = answer_provider
        self._audit_port = audit_port

    async def execute(self, command: AdmissionsQueryCommand) -> QueryResult:
        graph_context = await self._graph_retriever.retrieve(
            question=command.question,
            collection=command.context.collection_name,
            limit=command.top_k,
        )
        answer = await self._answer_provider.answer(
            question=command.question,
            context=graph_context,
        )
        query_id = str(uuid4())
        citations = [
            CitationSource(
                document_id=item.source_id,
                chunk_index=idx,
                excerpt=item.text[:300],
                collection=item.collection,
                score=item.score,
                source_type=item.source_type,
                title=item.title,
                path=item.path,
            )
            for idx, item in enumerate(graph_context)
        ]
        document_ids: list[str] = []
        snippets: list[str] = []

        for item in graph_context:
            if item.source_id not in document_ids:
                document_ids.append(item.source_id)
            snippets.append(item.text[:200])

        result = QueryResult(
            answer=answer,
            document_ids=document_ids,
            snippets=snippets,
            citations=citations,
        )

        await self._audit_port.record(
            QueryAuditRecord(
                query_id=query_id,
                subject=command.context.subject,
                scope=command.context.scope.value,
                collection=command.context.collection_name,
                application_id=command.context.application_id,
                question_hash=sha256(command.question.encode("utf-8")).hexdigest(),
                citation_count=len(citations),
                created_at=datetime.now(timezone.utc),
            )
        )

        return result
