from domain.entities.document import DocumentRecord
from domain.ports.document_repository import IDocumentRepository


class InMemoryDocumentRepository(IDocumentRepository):
    def __init__(self) -> None:
        self._documents: dict[str, DocumentRecord] = {}

    async def save(self, document: DocumentRecord) -> DocumentRecord:
        self._documents[document.id] = document
        return document

    async def get(self, document_id: str) -> DocumentRecord | None:
        return self._documents.get(document_id)

    async def delete(self, document_id: str) -> DocumentRecord | None:
        return self._documents.pop(document_id, None)
