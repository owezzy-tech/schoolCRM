from abc import ABC, abstractmethod
from dataclasses import dataclass

from domain.entities.document import DocumentChunk
from domain.types import Embedding


@dataclass(slots=True)
class SearchMatch:
    chunk: DocumentChunk
    score: float


class IVectorStore(ABC):
    @abstractmethod
    async def upsert_chunks(self, chunks: list[DocumentChunk], embeddings: list[Embedding]) -> None:
        raise NotImplementedError

    @abstractmethod
    async def search(self, embedding: Embedding, limit: int) -> list[SearchMatch]:
        raise NotImplementedError

    @abstractmethod
    async def delete_document(self, document_id: str) -> None:
        raise NotImplementedError
