from abc import ABC, abstractmethod

from domain.ports.vector_store import SearchMatch


class ILLMProvider(ABC):
    @abstractmethod
    async def answer(self, *, question: str, context: list[SearchMatch]) -> str:
        raise NotImplementedError
