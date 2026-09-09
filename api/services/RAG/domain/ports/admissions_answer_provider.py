from abc import ABC, abstractmethod

from domain.entities.graph_context import GraphContext


class IAdmissionsAnswerProvider(ABC):
    @abstractmethod
    async def answer(self, *, question: str, context: list[GraphContext]) -> str:
        raise NotImplementedError
