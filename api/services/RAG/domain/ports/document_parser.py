from abc import ABC, abstractmethod


class IDocumentParser(ABC):
    @abstractmethod
    async def parse(self, *, filename: str, content_type: str, payload: bytes) -> str:
        raise NotImplementedError
