from dataclasses import dataclass

from domain.entities.query import CitationSource, QueryResult
from domain.ports.embedding_provider import IEmbeddingProvider
from domain.ports.llm_provider import ILLMProvider
from domain.ports.vector_store import IVectorStore


@dataclass(slots=True)
class QueryDocumentsCommand:
    question: str
    requested_by: str
    top_k: int = 3
    collection: str | None = None


class QueryDocumentsUseCase:
    def __init__(
        self,
        *,
        embedding_provider: IEmbeddingProvider,
        vector_store: IVectorStore,
        llm_provider: ILLMProvider,
    ) -> None:
        self._embedding_provider = embedding_provider
        self._vector_store = vector_store
        self._llm_provider = llm_provider

    async def execute(self, command: QueryDocumentsCommand) -> QueryResult:
        question_embedding = await self._embedding_provider.embed_text(command.question)
        matches = await self._vector_store.search(
            question_embedding,
            command.top_k,
            command.collection,
        )
        answer = await self._llm_provider.answer(question=command.question, context=matches)

        document_ids: list[str] = []
        snippets: list[str] = []
        citations: list[CitationSource] = []
        for match in matches:
            if match.chunk.document_id not in document_ids:
                document_ids.append(match.chunk.document_id)
            snippets.append(match.chunk.text[:200])

            citations.append(
                CitationSource(
                    document_id=match.chunk.document_id,
                    chunk_index=match.chunk.metadata.chunk_index,
                    excerpt=match.chunk.text[:300],
                    collection=match.chunk.metadata.collection,
                    score=match.score,
                )
            )

        return QueryResult(
            answer=answer,
            document_ids=document_ids,
            snippets=snippets,
            citations=citations,
        )
