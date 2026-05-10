from abc import ABC, abstractmethod


class IFileStore(ABC):
    @abstractmethod
    async def save(self, *, document_id: str, filename: str, payload: bytes) -> str:
        raise NotImplementedError

    @abstractmethod
    async def delete(self, path: str | None) -> None:
        raise NotImplementedError
