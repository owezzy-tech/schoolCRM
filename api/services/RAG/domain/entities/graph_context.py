from dataclasses import dataclass


@dataclass(frozen=True, slots=True)
class GraphContext:
    source_id: str
    source_type: str
    title: str
    text: str
    collection: str
    score: float
    path: list[str]
