from dataclasses import dataclass
from typing import Any

from fastapi import Request

from adapters.application_ownership.stub_application_ownership_checker import (
    StubApplicationOwnershipChecker,
)
from adapters.audit.logging_query_audit_adapter import LoggingQueryAuditAdapter
from adapters.embeddings.noop_embedding_provider import NoopEmbeddingProvider
from adapters.graph.in_memory_graph_retriever import InMemoryGraphRetriever
from adapters.llm.echo_llm_provider import EchoLLMProvider
from adapters.llm.graph_context_answer_provider import GraphContextAnswerProvider
from adapters.llm.ollama_answer_provider import OllamaAnswerProvider
from adapters.parsers.noop_document_parser import NoopDocumentParser
from adapters.repositories.in_memory_document_repository import InMemoryDocumentRepository
from adapters.storage.local_file_store import LocalFileStore
from adapters.vector_stores.in_memory_vector_store import InMemoryVectorStore
from domain.ports.application_ownership import IApplicationOwnershipChecker
from infrastructure.auth_client import AuthServiceClient
from infrastructure.config import Settings
from use_cases.admissions_query import AdmissionsQueryUseCase
from use_cases.delete_document import DeleteDocumentUseCase
from use_cases.ingest_document import IngestDocumentUseCase
from use_cases.query_documents import QueryDocumentsUseCase


@dataclass(slots=True)
class Container:
    ingest_document: IngestDocumentUseCase
    query_documents: QueryDocumentsUseCase
    admissions_query: AdmissionsQueryUseCase
    application_ownership_checker: IApplicationOwnershipChecker
    delete_document: DeleteDocumentUseCase
    auth_service: AuthServiceClient
    neo4j_driver: Any | None = None
    ollama_answer_provider: OllamaAnswerProvider | None = None

    async def close(self) -> None:
        await self.auth_service.close()
        if self.neo4j_driver is not None:
            await self.neo4j_driver.close()
        if self.ollama_answer_provider is not None:
            await self.ollama_answer_provider.close()


def build_container(settings: Settings) -> Container:
    parser = NoopDocumentParser()
    embedding_provider = NoopEmbeddingProvider()
    vector_store = InMemoryVectorStore()
    llm_provider = EchoLLMProvider()
    document_repository = InMemoryDocumentRepository()
    file_store = LocalFileStore(settings.file_storage_dir)
    neo4j_driver: Any | None = None
    if settings.graph_retriever == "neo4j":
        from neo4j import AsyncGraphDatabase

        from adapters.graph.neo4j_graph_retriever import Neo4jGraphRetriever

        neo4j_driver = AsyncGraphDatabase.driver(
            settings.neo4j_uri,
            auth=(settings.neo4j_username, settings.neo4j_password),
        )
        graph_retriever = Neo4jGraphRetriever(neo4j_driver)
    else:
        graph_retriever = InMemoryGraphRetriever()

    ollama_answer_provider: OllamaAnswerProvider | None = None
    if settings.admissions_answer_provider == "ollama":
        ollama_answer_provider = OllamaAnswerProvider(
            base_url=settings.ollama_base_url,
            model=settings.ollama_model,
            timeout_seconds=settings.ollama_timeout_seconds,
        )
        answer_provider = ollama_answer_provider
    else:
        answer_provider = GraphContextAnswerProvider()

    application_ownership_checker = StubApplicationOwnershipChecker()
    audit_port = LoggingQueryAuditAdapter()

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
        admissions_query=AdmissionsQueryUseCase(
            graph_retriever=graph_retriever,
            answer_provider=answer_provider,
            audit_port=audit_port,
        ),
        application_ownership_checker=application_ownership_checker,
        delete_document=DeleteDocumentUseCase(
            vector_store=vector_store,
            document_repository=document_repository,
            file_store=file_store,
        ),
        auth_service=AuthServiceClient(
            base_url=settings.auth_service_url,
            timeout_seconds=settings.auth_request_timeout_seconds,
        ),
        neo4j_driver=neo4j_driver,
        ollama_answer_provider=ollama_answer_provider,
    )


def get_ingest_document_use_case(request: Request) -> IngestDocumentUseCase:
    return request.app.state.container.ingest_document


def get_query_documents_use_case(request: Request) -> QueryDocumentsUseCase:
    return request.app.state.container.query_documents


def get_admissions_query_use_case(request: Request) -> AdmissionsQueryUseCase:
    return request.app.state.container.admissions_query


def get_application_ownership_checker(request: Request) -> IApplicationOwnershipChecker:
    return request.app.state.container.application_ownership_checker


def get_delete_document_use_case(request: Request) -> DeleteDocumentUseCase:
    return request.app.state.container.delete_document


def get_auth_service_client(request: Request) -> AuthServiceClient:
    return request.app.state.container.auth_service
