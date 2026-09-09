import asyncio

from domain.entities.admissions_query_context import AdmissionsQueryContext, AdmissionsScope
from domain.entities.audit import QueryAuditRecord
from domain.entities.graph_context import GraphContext
from use_cases.admissions_query import AdmissionsQueryCommand, AdmissionsQueryUseCase


class FakeGraphRetriever:
    def __init__(self) -> None:
        self.question: str | None = None
        self.collection: str | None = None
        self.limit: int | None = None

    async def retrieve(self, *, question: str, collection: str, limit: int) -> list[GraphContext]:
        self.question = question
        self.collection = collection
        self.limit = limit
        return [
            GraphContext(
                source_id="policy-1",
                source_type="Policy",
                title="Admissions deadline",
                text="Applications close on 31 March for the regular admissions cycle.",
                collection=collection,
                score=0.91,
                path=["AdmissionsPolicy", "Deadline"],
            ),
            GraphContext(
                source_id="policy-1",
                source_type="Policy",
                title="Admissions deadline duplicate path",
                text="Late applications require admissions office approval.",
                collection=collection,
                score=0.72,
                path=["AdmissionsPolicy", "LateApplication"],
            ),
        ]


class FakeAnswerProvider:
    def __init__(self) -> None:
        self.question: str | None = None
        self.context: list[GraphContext] | None = None

    async def answer(self, *, question: str, context: list[GraphContext]) -> str:
        self.question = question
        self.context = context
        return "Applications close on 31 March."


class FakeAuditPort:
    def __init__(self) -> None:
        self.records: list[QueryAuditRecord] = []

    async def record(self, record: QueryAuditRecord) -> None:
        self.records.append(record)


def test_admissions_query_uses_graph_context_for_answer_citations_and_audit() -> None:
    graph_retriever = FakeGraphRetriever()
    answer_provider = FakeAnswerProvider()
    audit_port = FakeAuditPort()
    use_case = AdmissionsQueryUseCase(
        graph_retriever=graph_retriever,
        answer_provider=answer_provider,
        audit_port=audit_port,
    )

    result = asyncio.run(
        use_case.execute(
            AdmissionsQueryCommand(
                question="When do applications close?",
                top_k=2,
                context=AdmissionsQueryContext(
                    scope=AdmissionsScope.STAFF,
                    collection_name="admissions-policy",
                    subject="staff-1",
                ),
            )
        )
    )

    assert graph_retriever.question == "When do applications close?"
    assert graph_retriever.collection == "admissions-policy"
    assert graph_retriever.limit == 2
    assert answer_provider.question == "When do applications close?"
    assert answer_provider.context is not None
    assert result.answer == "Applications close on 31 March."
    assert result.document_ids == ["policy-1"]
    assert result.snippets == [
        "Applications close on 31 March for the regular admissions cycle.",
        "Late applications require admissions office approval.",
    ]
    assert result.citations is not None
    assert len(result.citations) == 2
    assert result.citations[0].document_id == "policy-1"
    assert result.citations[0].source_type == "Policy"
    assert result.citations[0].title == "Admissions deadline"
    assert result.citations[0].path == ["AdmissionsPolicy", "Deadline"]
    assert len(audit_port.records) == 1
    assert audit_port.records[0].subject == "staff-1"
    assert audit_port.records[0].scope == "staff"
    assert audit_port.records[0].collection == "admissions-policy"
    assert audit_port.records[0].citation_count == 2
    assert len(audit_port.records[0].question_hash) == 64
