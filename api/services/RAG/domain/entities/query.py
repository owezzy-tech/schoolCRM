from dataclasses import dataclass


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
