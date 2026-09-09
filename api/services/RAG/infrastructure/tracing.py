import os
from collections.abc import Awaitable, Callable, Sequence
from functools import wraps
from hashlib import sha256
from typing import Any, ParamSpec, TypeVar

from langsmith import traceable, tracing_context

from domain.entities.document import DocumentChunk
from domain.entities.graph_context import GraphContext
from domain.ports.admissions_answer_provider import IAdmissionsAnswerProvider
from domain.ports.document_parser import IDocumentParser
from domain.ports.document_repository import IDocumentRepository
from domain.ports.embedding_provider import IEmbeddingProvider
from domain.ports.file_store import IFileStore
from domain.ports.graph_retriever import IGraphRetriever
from domain.ports.llm_provider import ILLMProvider
from domain.ports.vector_store import IVectorStore, SearchMatch
from domain.types import Embedding
from infrastructure.config import Settings
from use_cases.admissions_query import AdmissionsQueryCommand, AdmissionsQueryUseCase
from use_cases.ingest_document import IngestDocumentCommand, IngestDocumentUseCase
from use_cases.query_documents import QueryDocumentsCommand, QueryDocumentsUseCase

P = ParamSpec("P")
T = TypeVar("T")


def configure_langsmith(settings: Settings) -> None:
    if not settings.langsmith_enabled:
        os.environ["LANGSMITH_TRACING"] = "false"
        return

    os.environ["LANGSMITH_TRACING"] = "true"
    os.environ["LANGSMITH_PROJECT"] = settings.langsmith_project

    if settings.langsmith_api_key:
        os.environ["LANGSMITH_API_KEY"] = settings.langsmith_api_key

    if settings.langsmith_endpoint:
        os.environ["LANGSMITH_ENDPOINT"] = settings.langsmith_endpoint


def trace_container(
    *,
    settings: Settings,
    ingest_document: IngestDocumentUseCase,
    query_documents: QueryDocumentsUseCase,
    admissions_query: AdmissionsQueryUseCase,
    graph_retriever: IGraphRetriever,
    answer_provider: IAdmissionsAnswerProvider,
    embedding_provider: IEmbeddingProvider,
    vector_store: IVectorStore,
    llm_provider: ILLMProvider,
    audit_port: Any,
    file_store: IFileStore,
    parser: IDocumentParser,
    document_repository: IDocumentRepository,
) -> tuple[
    IngestDocumentUseCase,
    QueryDocumentsUseCase,
    AdmissionsQueryUseCase,
    IGraphRetriever,
    IAdmissionsAnswerProvider,
    IEmbeddingProvider,
    IVectorStore,
    ILLMProvider,
]:
    if not settings.langsmith_enabled:
        return (
            ingest_document,
            query_documents,
            admissions_query,
            graph_retriever,
            answer_provider,
            embedding_provider,
            vector_store,
            llm_provider,
        )

    traced_graph_retriever = TracedGraphRetriever(graph_retriever, settings.langsmith_project)
    traced_answer_provider = TracedAdmissionsAnswerProvider(
        answer_provider,
        settings.langsmith_project,
    )
    traced_embedding_provider = TracedEmbeddingProvider(
        embedding_provider,
        settings.langsmith_project,
    )
    traced_vector_store = TracedVectorStore(vector_store, settings.langsmith_project)
    traced_llm_provider = TracedLLMProvider(llm_provider, settings.langsmith_project)

    return (
        TracedIngestDocumentUseCase(
            IngestDocumentUseCase(
                file_store=file_store,
                parser=parser,
                embedding_provider=traced_embedding_provider,
                vector_store=traced_vector_store,
                document_repository=document_repository,
            ),
            settings.langsmith_project,
        ),
        TracedQueryDocumentsUseCase(
            QueryDocumentsUseCase(
                embedding_provider=traced_embedding_provider,
                vector_store=traced_vector_store,
                llm_provider=traced_llm_provider,
            ),
            settings.langsmith_project,
        ),
        TracedAdmissionsQueryUseCase(
            AdmissionsQueryUseCase(
                graph_retriever=traced_graph_retriever,
                answer_provider=traced_answer_provider,
                audit_port=audit_port,
            ),
            settings.langsmith_project,
        ),
        traced_graph_retriever,
        traced_answer_provider,
        traced_embedding_provider,
        traced_vector_store,
        traced_llm_provider,
    )


class TracedAdmissionsQueryUseCase:
    def __init__(self, wrapped: AdmissionsQueryUseCase, project_name: str) -> None:
        self._wrapped = wrapped
        self._project_name = project_name

    async def execute(self, command: AdmissionsQueryCommand):
        traced = traceable(
            name="admissions_rag_query",
            run_type="chain",
            process_inputs=lambda inputs: {
                "question_hash": _hash_text(inputs["command"].question),
                "scope": inputs["command"].context.scope.value,
                "collection": inputs["command"].context.collection_name,
                "application_id": inputs["command"].context.application_id,
                "top_k": inputs["command"].top_k,
            },
            process_outputs=_query_result_outputs,
        )(self._wrapped.execute)

        with tracing_context(project_name=self._project_name):
            return await traced(command)


class TracedQueryDocumentsUseCase:
    def __init__(self, wrapped: QueryDocumentsUseCase, project_name: str) -> None:
        self._wrapped = wrapped
        self._project_name = project_name

    async def execute(self, command: QueryDocumentsCommand):
        traced = traceable(
            name="document_rag_query",
            run_type="chain",
            process_inputs=lambda inputs: {
                "question_hash": _hash_text(inputs["command"].question),
                "requested_by": inputs["command"].requested_by,
                "collection": inputs["command"].collection,
                "top_k": inputs["command"].top_k,
            },
            process_outputs=_query_result_outputs,
        )(self._wrapped.execute)

        with tracing_context(project_name=self._project_name):
            return await traced(command)


class TracedIngestDocumentUseCase:
    def __init__(self, wrapped: IngestDocumentUseCase, project_name: str) -> None:
        self._wrapped = wrapped
        self._project_name = project_name

    async def execute(self, command: IngestDocumentCommand):
        traced = traceable(
            name="ingest_document",
            run_type="chain",
            process_inputs=lambda inputs: {
                "title": inputs["command"].title,
                "source": inputs["command"].source,
                "uploaded_by": inputs["command"].uploaded_by,
                "filename": inputs["command"].filename,
                "content_type": inputs["command"].content_type,
                "payload_size": len(inputs["command"].payload),
            },
            process_outputs=lambda outputs: {
                "document_id": outputs.document_id,
                "status": outputs.status,
                "chunk_count": outputs.chunk_count,
            },
        )(self._wrapped.execute)

        with tracing_context(project_name=self._project_name):
            return await traced(command)


class TracedGraphRetriever(IGraphRetriever):
    def __init__(self, wrapped: IGraphRetriever, project_name: str) -> None:
        self._wrapped = wrapped
        self._project_name = project_name

    async def retrieve(self, *, question: str, collection: str, limit: int) -> list[GraphContext]:
        traced = traceable(
            name="graph_context_retrieve",
            run_type="retriever",
            process_inputs=lambda inputs: {
                "question_hash": _hash_text(inputs["question"]),
                "collection": inputs["collection"],
                "limit": inputs["limit"],
            },
            process_outputs=lambda outputs: {
                "context_count": len(outputs),
                "source_ids": [item.source_id for item in outputs],
            },
        )(self._wrapped.retrieve)

        with tracing_context(project_name=self._project_name):
            return await traced(question=question, collection=collection, limit=limit)


class TracedAdmissionsAnswerProvider(IAdmissionsAnswerProvider):
    def __init__(self, wrapped: IAdmissionsAnswerProvider, project_name: str) -> None:
        self._wrapped = wrapped
        self._project_name = project_name

    async def answer(self, *, question: str, context: list[GraphContext]) -> str:
        traced = traceable(
            name="admissions_answer_generate",
            run_type="llm",
            process_inputs=lambda inputs: {
                "question_hash": _hash_text(inputs["question"]),
                "context_count": len(inputs["context"]),
                "source_ids": [item.source_id for item in inputs["context"]],
            },
            process_outputs=lambda output: {"answer_length": len(output)},
        )(self._wrapped.answer)

        with tracing_context(project_name=self._project_name):
            return await traced(question=question, context=context)


class TracedEmbeddingProvider(IEmbeddingProvider):
    def __init__(self, wrapped: IEmbeddingProvider, project_name: str) -> None:
        self._wrapped = wrapped
        self._project_name = project_name

    async def embed_text(self, text: str):
        traced = traceable(
            name="embed_text",
            run_type="embedding",
            process_inputs=lambda inputs: {"text_hash": _hash_text(inputs["text"])},
            process_outputs=lambda output: {"dimensions": len(output)},
        )(self._wrapped.embed_text)

        with tracing_context(project_name=self._project_name):
            return await traced(text)


class TracedVectorStore(IVectorStore):
    def __init__(self, wrapped: IVectorStore, project_name: str) -> None:
        self._wrapped = wrapped
        self._project_name = project_name

    async def upsert_chunks(self, chunks: list[DocumentChunk], embeddings: list[Embedding]) -> None:
        traced = traceable(
            name="vector_upsert_chunks",
            run_type="tool",
            process_inputs=lambda inputs: {
                "chunk_count": len(inputs["chunks"]),
                "document_ids": _document_ids(inputs["chunks"]),
                "embedding_count": len(inputs["embeddings"]),
            },
            process_outputs=lambda _: {"status": "indexed"},
        )(self._wrapped.upsert_chunks)

        with tracing_context(project_name=self._project_name):
            return await traced(chunks, embeddings)

    async def search(
        self,
        embedding: Embedding,
        limit: int,
        collection: str | None = None,
    ) -> list[SearchMatch]:
        traced = traceable(
            name="vector_search",
            run_type="retriever",
            process_inputs=lambda inputs: {
                "embedding_dimensions": len(inputs["embedding"]),
                "limit": inputs["limit"],
                "collection": inputs["collection"],
            },
            process_outputs=lambda outputs: {
                "match_count": len(outputs),
                "document_ids": [match.chunk.document_id for match in outputs],
            },
        )(self._wrapped.search)

        with tracing_context(project_name=self._project_name):
            return await traced(embedding, limit, collection)

    async def delete_document(self, document_id: str) -> None:
        return await self._wrapped.delete_document(document_id)


class TracedLLMProvider(ILLMProvider):
    def __init__(self, wrapped: ILLMProvider, project_name: str) -> None:
        self._wrapped = wrapped
        self._project_name = project_name

    async def answer(self, *, question: str, context: list[SearchMatch]) -> str:
        traced = traceable(
            name="document_answer_generate",
            run_type="llm",
            process_inputs=lambda inputs: {
                "question_hash": _hash_text(inputs["question"]),
                "context_count": len(inputs["context"]),
                "document_ids": [match.chunk.document_id for match in inputs["context"]],
            },
            process_outputs=lambda output: {"answer_length": len(output)},
        )(self._wrapped.answer)

        with tracing_context(project_name=self._project_name):
            return await traced(question=question, context=context)


def trace_async(
    fn: Callable[P, Awaitable[T]],
    *,
    enabled: bool,
    project_name: str,
    name: str,
    run_type: str,
) -> Callable[P, Awaitable[T]]:
    if not enabled:
        return fn

    traced = traceable(name=name, run_type=run_type)(fn)

    @wraps(fn)
    async def wrapper(*args: P.args, **kwargs: P.kwargs) -> T:
        with tracing_context(project_name=project_name):
            return await traced(*args, **kwargs)

    return wrapper


def _hash_text(text: str) -> str:
    return sha256(text.encode("utf-8")).hexdigest()


def _query_result_outputs(output: Any) -> dict[str, Any]:
    return {
        "citation_count": len(output.citations),
        "document_ids": output.document_ids,
        "answer_length": len(output.answer),
    }


def _document_ids(chunks: Sequence[DocumentChunk]) -> list[str]:
    document_ids: list[str] = []
    for chunk in chunks:
        if chunk.document_id not in document_ids:
            document_ids.append(chunk.document_id)

    return document_ids
