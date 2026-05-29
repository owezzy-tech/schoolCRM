from dataclasses import dataclass

from domain.ports.document_repository import IDocumentRepository
from domain.ports.file_store import IFileStore
from domain.ports.vector_store import IVectorStore


@dataclass(slots=True)
class DeleteDocumentResult:
    document_id: str
    deleted: bool


class DeleteDocumentUseCase:
    def __init__(
        self,
        *,
        vector_store: IVectorStore,
        document_repository: IDocumentRepository,
        file_store: IFileStore,
    ) -> None:
        self._vector_store = vector_store
        self._document_repository = document_repository
        self._file_store = file_store

    async def execute(self, document_id: str) -> DeleteDocumentResult:
        document = await self._document_repository.delete(document_id)
        if document is None:
            return DeleteDocumentResult(document_id=document_id, deleted=False)

        await self._vector_store.delete_document(document_id)
        await self._file_store.delete(document.file_path)
        return DeleteDocumentResult(document_id=document_id, deleted=True)
