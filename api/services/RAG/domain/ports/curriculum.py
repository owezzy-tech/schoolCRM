from typing import Protocol
from uuid import UUID

from domain.entities.curriculum import (
    CurriculumEvidence,
    CurriculumPassage,
    CurriculumScope,
    CurriculumSource,
)


class CurriculumRepository(Protocol):
    async def create(
        self, source: CurriculumSource, content: bytes, passages: list[CurriculumPassage]
    ) -> None: ...

    async def get(self, source_id: UUID) -> CurriculumSource: ...

    async def passages(self, source_id: UUID) -> list[CurriculumPassage]: ...

    async def review(
        self,
        source_id: UUID,
        actor: str,
        expected: str,
        decision: str,
        note: str,
        model_identity: str | None,
        vectors: list[list[float]],
    ) -> CurriculumSource: ...

    async def search(
        self,
        scope: CurriculumScope,
        model_identity: str,
        vector: list[float],
        limit: int,
        minimum_score: float,
    ) -> list[CurriculumEvidence]: ...

    async def original(self, source_id: UUID) -> bytes: ...


class CurriculumEmbeddings(Protocol):
    async def embed(self, texts: list[str]) -> tuple[str, list[list[float]]]: ...

    async def close(self) -> None: ...


class CurriculumAccess(Protocol):
    async def authorize(self, token: str | None, scope: CurriculumScope, manage: bool) -> None: ...

    async def close(self) -> None: ...
