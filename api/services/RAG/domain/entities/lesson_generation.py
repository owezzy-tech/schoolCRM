from dataclasses import dataclass, field
from uuid import UUID

from domain.entities.curriculum import CurriculumScope


@dataclass(frozen=True, slots=True)
class LessonGenerationRequest:
    request_id: UUID
    scope: CurriculumScope
    topic: str
    duration_minutes: int


@dataclass(frozen=True, slots=True)
class GeneratedLesson:
    title: str
    content: dict


@dataclass(frozen=True, slots=True)
class GenerationContext:
    actor_id: UUID
    token: str = field(repr=False)
