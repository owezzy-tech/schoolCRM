from dataclasses import dataclass

from fastapi import Request

from adapters.embeddings.noop_embedding_provider import NoopEmbeddingProvider
from adapters.llm.echo_llm_provider import EchoLLMProvider
from adapters.parsers.noop_document_parser import NoopDocumentParser
from adapters.repositories.in_memory_document_repository import InMemoryDocumentRepository
from adapters.storage.local_file_store import LocalFileStore
from adapters.vector_stores.in_memory_vector_store import InMemoryVectorStore
from infrastructure.config import Settings
from use_cases.delete_document import DeleteDocumentUseCase
from use_cases.ingest_document import IngestDocumentUseCase
from use_cases.query_documents import QueryDocumentsUseCase


@dataclass(slots=True)
class Container:
    ingest_document: IngestDocumentUseCase
    query_documents: QueryDocumentsUseCase
    delete_document: DeleteDocumentUseCase

    async def close(self) -> None:
        return None


def build_container(settings: Settings) -> Container:
    parser = NoopDocumentParser()
    embedding_provider = NoopEmbeddingProvider()
    vector_store = InMemoryVectorStore()
    llm_provider = EchoLLMProvider()
    document_repository = InMemoryDocumentRepository()
    file_store = LocalFileStore(settings.file_storage_dir)

    return Container(
        ingest_document=IngestDocumentUseCase(
            file_store=file_store,
            parser=parser,
            embedding_provider=embedding_provider,
            vector_store=vector_store,
            document_repository=document_repository,
        ),
        query_documents=QueryDocumentsUseCase(
            embedding_provider=embedding_provider,
            vector_store=vector_store,
            llm_provider=llm_provider,
        ),
        delete_document=DeleteDocumentUseCase(
            vector_store=vector_store,
            document_repository=document_repository,
            file_store=file_store,
        ),
    )


def get_ingest_document_use_case(request: Request) -> IngestDocumentUseCase:
    return request.app.state.container.ingest_document


def get_query_documents_use_case(request: Request) -> QueryDocumentsUseCase:
    return request.app.state.container.query_documents


def get_delete_document_use_case(request: Request) -> DeleteDocumentUseCase:
    return request.app.state.container.delete_document
