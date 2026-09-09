from domain.entities.graph_context import GraphContext
from domain.ports.admissions_answer_provider import IAdmissionsAnswerProvider


class GraphContextAnswerProvider(IAdmissionsAnswerProvider):
    async def answer(self, *, question: str, context: list[GraphContext]) -> str:
        if not context:
            return f"No admissions knowledge graph context was found for: {question}"

        snippets = "; ".join(item.text[:120] for item in context[:2])
        return f"Admissions graph answer for '{question}'. Context snippets: {snippets}"

    async def close(self) -> None:
        pass
