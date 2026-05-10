from abc import ABC, abstractmethod

from domain.types import Embedding


class IEmbeddingProvider(ABC):
    @abstractmethod
    async def embed_text(self, text: str) -> Embedding:
        raise NotImplementedError
