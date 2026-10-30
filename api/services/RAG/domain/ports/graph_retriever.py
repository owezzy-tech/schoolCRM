from abc import ABC, abstractmethod

from domain.entities.graph_context import GraphContext


class IGraphRetriever(ABC):
    @abstractmethod
    async def retrieve(self, *, question: str, collection: str, limit: int) -> list[GraphContext]:
        raise NotImplementedError
