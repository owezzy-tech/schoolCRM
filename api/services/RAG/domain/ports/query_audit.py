from abc import ABC, abstractmethod

from domain.entities.audit import QueryAuditRecord


class IQueryAuditPort(ABC):
    @abstractmethod
    async def record(self, entry: QueryAuditRecord) -> None:
        raise NotImplementedError
