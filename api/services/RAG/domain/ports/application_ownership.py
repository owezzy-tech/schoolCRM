from abc import ABC, abstractmethod


class IApplicationOwnershipChecker(ABC):
    @abstractmethod
    async def is_owner(self, subject: str, application_id: str) -> bool:
        raise NotImplementedError

    @abstractmethod
    async def is_assigned_reviewer(self, subject: str, application_id: str) -> bool:
        raise NotImplementedError
