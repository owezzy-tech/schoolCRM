import asyncio
import os

from adapters.graph.in_memory_graph_retriever import InMemoryGraphRetriever
from infrastructure.config import Settings
from infrastructure.dependencies import build_container
from infrastructure.tracing import (
    TracedAdmissionsQueryUseCase,
    TracedGraphRetriever,
    configure_langsmith,
)
from use_cases.admissions_query import AdmissionsQueryUseCase


def test_configure_langsmith_disables_tracing_by_default() -> None:
    os.environ["LANGSMITH_TRACING"] = "true"

    configure_langsmith(Settings(langsmith_enabled=False))

    assert os.environ["LANGSMITH_TRACING"] == "false"


def test_configure_langsmith_exports_enabled_environment() -> None:
    settings = Settings(
        langsmith_enabled=True,
        langsmith_api_key="test-key",
        langsmith_project="schoolcrm-rag-test",
        langsmith_endpoint="https://example.invalid",
    )

    configure_langsmith(settings)

    assert os.environ["LANGSMITH_TRACING"] == "true"
    assert os.environ["LANGSMITH_API_KEY"] == "test-key"
    assert os.environ["LANGSMITH_PROJECT"] == "schoolcrm-rag-test"
    assert os.environ["LANGSMITH_ENDPOINT"] == "https://example.invalid"


def test_langsmith_disabled_container_keeps_original_components() -> None:
    container = build_container(Settings(langsmith_enabled=False))

    assert isinstance(container.admissions_query, AdmissionsQueryUseCase)

    asyncio.run(container.close())


def test_langsmith_enabled_container_wraps_rag_components() -> None:
    settings = Settings(
        graph_retriever="memory",
        admissions_answer_provider="stub",
        langsmith_enabled=True,
        langsmith_project="schoolcrm-rag-test",
    )
    container = build_container(settings)

    assert isinstance(container.admissions_query, TracedAdmissionsQueryUseCase)
    assert isinstance(container.admissions_query._wrapped._graph_retriever, TracedGraphRetriever)
    assert isinstance(
        container.admissions_query._wrapped._graph_retriever._wrapped,
        InMemoryGraphRetriever,
    )

    asyncio.run(container.close())
