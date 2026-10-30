from dataclasses import dataclass
from uuid import UUID

REVIEW_TRANSITIONS = frozenset({("pending", "approved"), ("approved", "withdrawn")})


@dataclass(frozen=True, slots=True)
class CurriculumScope:
    school_id: UUID
    department_id: UUID
    framework: str
    stage: str
    subject: str
    revision: str


@dataclass(frozen=True, slots=True)
class CurriculumSource:
    id: UUID
    scope: CurriculumScope
    title: str
    topic: str
    authority: str
    source_url: str
    sha256: str
    filename: str
    first_page: int
    last_page: int
    uploaded_by: str
    status: str = "pending"
    reviewed_by: str | None = None
    review_note: str | None = None
    model_identity: str | None = None
    embedding_dimensions: int = 1024
    parser_revision: str = "llamaindex-sentence-400-50-v1"


@dataclass(frozen=True, slots=True)
class CurriculumPassage:
    page: int
    ordinal: int
    text: str


@dataclass(frozen=True, slots=True)
class CurriculumEvidence:
    source: CurriculumSource
    passage: CurriculumPassage
    score: float
    model_identity: str
    index_revision: UUID


class CurriculumError(Exception):
    def __init__(self, status_code: int, detail: str) -> None:
        self.status_code = status_code
        self.detail = detail
        super().__init__(detail)
