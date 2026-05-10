from hashlib import sha256

from domain.ports.embedding_provider import IEmbeddingProvider


class NoopEmbeddingProvider(IEmbeddingProvider):
    async def embed_text(self, text: str) -> list[float]:
        digest = sha256(text.encode("utf-8")).digest()
        return [byte / 255 for byte in digest[:16]]
