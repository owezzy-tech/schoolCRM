import httpx

from domain.entities.graph_context import GraphContext
from domain.ports.admissions_answer_provider import IAdmissionsAnswerProvider


class OllamaAnswerProvider(IAdmissionsAnswerProvider):
    def __init__(self, *, base_url: str, model: str, timeout_seconds: float) -> None:
        self._client = httpx.AsyncClient(base_url=base_url.rstrip("/"), timeout=timeout_seconds)
        self._model = model

    async def answer(self, *, question: str, context: list[GraphContext]) -> str:
        context_text = "\n".join(
            f"[{idx + 1}] {item.title}: {item.text}" for idx, item in enumerate(context)
        )
        prompt = (
            "Answer the admissions question using only the provided Neo4j knowledge graph "
            "context. If the context is insufficient, say so.\n\n"
            f"Context:\n{context_text}\n\nQuestion: {question}"
        )
        response = await self._client.post(
            "/api/generate",
            json={"model": self._model, "prompt": prompt, "stream": False},
        )
        response.raise_for_status()
        payload = response.json()
        return str(payload.get("response", "")).strip()

    async def close(self) -> None:
        await self._client.aclose()
