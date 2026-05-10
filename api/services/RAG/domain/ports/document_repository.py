from abc import ABC, abstractmethod

from domain.entities.document import DocumentRecord


class IDocumentRepository(ABC):
    @abstractmethod
    async def save(self, document: DocumentRecord) -> DocumentRecord:
        raise NotImplementedError

    @abstractmethod
    async def get(self, document_id: str) -> DocumentRecord | None:
        raise NotImplementedError

    @abstractmethod
    async def delete(self, document_id: str) -> DocumentRecord | None:
        raise NotImplementedError
