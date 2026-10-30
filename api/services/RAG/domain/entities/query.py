from dataclasses import dataclass


@dataclass(frozen=True, slots=True)
class CitationSource:
    document_id: str
    chunk_index: int
    excerpt: str
    collection: str
    score: float
    source_type: str = "document"
    title: str = ""
    path: list[str] | None = None


@dataclass(slots=True)
class Query:
    question: str
    requested_by: str
    top_k: int = 3


@dataclass(slots=True)
class QueryResult:
    answer: str
    document_ids: list[str]
    snippets: list[str]
    citations: list[CitationSource] | None = None
