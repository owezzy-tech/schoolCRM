import asyncio
from collections.abc import Callable
from hashlib import sha256
from uuid import UUID, uuid4

from domain.entities.curriculum import (
    REVIEW_TRANSITIONS,
    CurriculumError,
    CurriculumEvidence,
    CurriculumPassage,
    CurriculumScope,
    CurriculumSource,
)
from domain.ports.curriculum import CurriculumAccess, CurriculumEmbeddings, CurriculumRepository


class CurriculumService:
    def __init__(
        self,
        store: CurriculumRepository,
        embeddings: CurriculumEmbeddings,
        access: CurriculumAccess,
        parse_pdf: Callable[[bytes, int, int], list[CurriculumPassage]],
    ) -> None:
        self.store = store
        self.embeddings = embeddings
        self.access = access
        self.parse_pdf = parse_pdf

    async def ingest(
        self,
        token: str | None,
        actor: str,
        scope: CurriculumScope,
        title: str,
        topic: str,
        authority: str,
        source_url: str,
        filename: str,
        first_page: int,
        last_page: int,
        expected_sha256: str,
        content: bytes,
    ) -> CurriculumSource:
        await self.access.authorize(token, scope, manage=False)
        if sha256(content).hexdigest() != expected_sha256:
            raise CurriculumError(422, "Original checksum does not match the acquisition record")
        passages = await asyncio.to_thread(self.parse_pdf, content, first_page, last_page)
        source = CurriculumSource(
            uuid4(),
            scope,
            title,
            topic,
            authority,
            source_url,
            expected_sha256,
            filename,
            first_page,
            last_page,
            actor,
        )
        # Parsing can take time. Recheck current authority before committing.
        await self.access.authorize(token, scope, manage=False)
        await self.store.create(source, content, passages)
        return source

    async def source(self, token: str | None, source_id: UUID) -> CurriculumSource:
        source = await self.store.get(source_id)
        await self.access.authorize(token, source.scope, manage=False)
        return source

    async def review(
        self,
        token: str | None,
        actor: str,
        source_id: UUID,
        expected: str,
        decision: str,
        note: str,
    ) -> CurriculumSource:
        source = await self.store.get(source_id)
        await self.access.authorize(token, source.scope, manage=True)
        if source.status != expected or (expected, decision) not in REVIEW_TRANSITIONS:
            raise CurriculumError(409, "Source state changed or invalid review transition")
        identity, vectors = None, []
        if decision == "approved":
            passages = await self.store.passages(source_id)
            for start in range(0, len(passages), 16):
                batch_identity, batch = await self.embeddings.embed(
                    [p.text for p in passages[start : start + 16]]
                )
                if identity is not None and identity != batch_identity:
                    raise CurriculumError(503, "Embedding model changed during indexing")
                identity = batch_identity
                vectors.extend(batch)
        await self.access.authorize(token, source.scope, manage=True)
        return await self.store.review(
            source_id, actor, expected, decision, note, identity, vectors
        )

    async def search(
        self, token: str | None, scope: CurriculumScope, question: str, limit: int
    ) -> list[CurriculumEvidence]:
        await self.access.authorize(token, scope, manage=False)
        identity, vectors = await self.embeddings.embed([question])
        # Deliberately require relevant evidence rather than inventing an answer.
        evidence = await self.store.search(scope, identity, vectors[0], limit, minimum_score=0.55)
        await self.access.authorize(token, scope, manage=False)
        return evidence

    async def close(self) -> None:
        await self.embeddings.close()
        await self.access.close()
