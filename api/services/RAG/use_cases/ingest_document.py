from dataclasses import dataclass
from uuid import uuid4

from domain.entities.document import DocumentChunk, DocumentRecord
from domain.ports.document_parser import IDocumentParser
from domain.ports.document_repository import IDocumentRepository
from domain.ports.embedding_provider import IEmbeddingProvider
from domain.ports.file_store import IFileStore
from domain.ports.vector_store import IVectorStore
from domain.value_objects.chunk import ChunkMetadata


@dataclass(slots=True)
class IngestDocumentCommand:
    title: str
    source: str
    uploaded_by: str
    filename: str
    content_type: str
    payload: bytes


@dataclass(slots=True)
class IngestDocumentResult:
    document_id: str
    status: str
    chunk_count: int


class IngestDocumentUseCase:
    def __init__(
        self,
        *,
        file_store: IFileStore,
        parser: IDocumentParser,
        embedding_provider: IEmbeddingProvider,
        vector_store: IVectorStore,
        document_repository: IDocumentRepository,
        chunk_size: int = 800,
        chunk_overlap: int = 100,
    ) -> None:
        self._file_store = file_store
        self._parser = parser
        self._embedding_provider = embedding_provider
        self._vector_store = vector_store
        self._document_repository = document_repository
        self._chunk_size = chunk_size
        self._chunk_overlap = chunk_overlap

    async def execute(self, command: IngestDocumentCommand) -> IngestDocumentResult:
        document_id = str(uuid4())
        file_path = await self._file_store.save(
            document_id=document_id,
            filename=command.filename,
            payload=command.payload,
        )
        parsed_text = await self._parser.parse(
            filename=command.filename,
            content_type=command.content_type,
            payload=command.payload,
        )
        chunks = self._split(document_id=document_id, source=command.source, text=parsed_text)
        embeddings = [await self._embedding_provider.embed_text(chunk.text) for chunk in chunks]
        await self._vector_store.upsert_chunks(chunks, embeddings)
        await self._document_repository.save(
            DocumentRecord(
                id=document_id,
                title=command.title,
                source=command.source,
                content_type=command.content_type,
                uploaded_by=command.uploaded_by,
                status="indexed",
                file_path=file_path,
                chunk_count=len(chunks),
            )
        )
        return IngestDocumentResult(
            document_id=document_id,
            status="indexed",
            chunk_count=len(chunks),
        )

    def _split(self, *, document_id: str, source: str, text: str) -> list[DocumentChunk]:
        if not text.strip():
            return []

        chunks: list[DocumentChunk] = []
        step = max(1, self._chunk_size - self._chunk_overlap)

        for idx, start in enumerate(range(0, len(text), step)):
            end = min(len(text), start + self._chunk_size)
            chunk_text = text[start:end].strip()
            if not chunk_text:
                continue
            chunks.append(
                DocumentChunk(
                    id=str(uuid4()),
                    document_id=document_id,
                    text=chunk_text,
                    metadata=ChunkMetadata(
                        source=source,
                        chunk_index=idx,
                        start_offset=start,
                        end_offset=end,
                    ),
                )
            )

        return chunks
