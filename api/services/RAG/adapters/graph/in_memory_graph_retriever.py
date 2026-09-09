from domain.entities.graph_context import GraphContext
from domain.ports.graph_retriever import IGraphRetriever


class InMemoryGraphRetriever(IGraphRetriever):
    def __init__(self, contexts: list[GraphContext] | None = None) -> None:
        self._contexts = contexts or []

    def add_context(self, context: GraphContext) -> None:
        self._contexts.append(context)

    async def retrieve(self, *, question: str, collection: str, limit: int) -> list[GraphContext]:
        terms = {term.lower() for term in question.split() if term.strip()}
        matches = [context for context in self._contexts if context.collection == collection]

        def rank(context: GraphContext) -> float:
            text = f"{context.title} {context.text}".lower()
            lexical_score = sum(1 for term in terms if term in text)
            return context.score + lexical_score

        matches.sort(key=rank, reverse=True)
        return matches[:limit]
