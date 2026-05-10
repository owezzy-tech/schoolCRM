from domain.ports.llm_provider import ILLMProvider
from domain.ports.vector_store import SearchMatch


class EchoLLMProvider(ILLMProvider):
    async def answer(self, *, question: str, context: list[SearchMatch]) -> str:
        if not context:
            return f"No indexed context was found for: {question}"

        snippets = "; ".join(match.chunk.text[:120] for match in context[:2])
        return f"Scaffold answer for '{question}'. Context snippets: {snippets}"
