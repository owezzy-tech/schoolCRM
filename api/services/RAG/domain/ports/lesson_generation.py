from typing import Protocol
from uuid import UUID

from domain.entities.curriculum import CurriculumEvidence, CurriculumScope
from domain.entities.lesson_generation import GeneratedLesson, LessonGenerationRequest


class LessonGenerator(Protocol):
    async def generate(
        self, request: LessonGenerationRequest, evidence: list[CurriculumEvidence]
    ) -> GeneratedLesson: ...

    async def close(self) -> None: ...


class LessonWriter(Protocol):
    async def authorize(self, token: str, school_id: UUID, department_id: UUID) -> None: ...

    async def create(
        self, token: str, scope: CurriculumScope, request_id: UUID, draft: GeneratedLesson
    ) -> dict: ...

    async def close(self) -> None: ...
