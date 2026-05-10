from dataclasses import dataclass

from domain.entities.document import DocumentChunk
from domain.ports.vector_store import IVectorStore, SearchMatch
from domain.types import Embedding


@dataclass(slots=True)
class _StoredChunk:
    chunk: DocumentChunk
    embedding: Embedding


class InMemoryVectorStore(IVectorStore):
    def __init__(self) -> None:
        self._chunks: list[_StoredChunk] = []

    async def upsert_chunks(self, chunks: list[DocumentChunk], embeddings: list[Embedding]) -> None:
        for chunk, embedding in zip(chunks, embeddings, strict=False):
            self._chunks.append(_StoredChunk(chunk=chunk, embedding=embedding))

    async def search(self, embedding: Embedding, limit: int) -> list[SearchMatch]:
        ranked = [
            SearchMatch(chunk=stored.chunk, score=self._dot(embedding, stored.embedding))
            for stored in self._chunks
        ]
        ranked.sort(key=lambda item: item.score, reverse=True)
        return ranked[:limit]

    async def delete_document(self, document_id: str) -> None:
        self._chunks = [stored for stored in self._chunks if stored.chunk.document_id != document_id]

    @staticmethod
    def _dot(left: Embedding, right: Embedding) -> float:
        size = min(len(left), len(right))
        return sum(left[idx] * right[idx] for idx in range(size))
