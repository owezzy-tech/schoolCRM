from domain.ports.application_ownership import IApplicationOwnershipChecker


class StubApplicationOwnershipChecker(IApplicationOwnershipChecker):
    def __init__(self, *, owner_allowed: bool = True, reviewer_allowed: bool = True) -> None:
        self.owner_allowed = owner_allowed
        self.reviewer_allowed = reviewer_allowed
        self.owner_checks: list[tuple[str, str]] = []
        self.reviewer_checks: list[tuple[str, str]] = []

    async def is_owner(self, subject: str, application_id: str) -> bool:
        self.owner_checks.append((subject, application_id))
        return self.owner_allowed

    async def is_assigned_reviewer(self, subject: str, application_id: str) -> bool:
        self.reviewer_checks.append((subject, application_id))
        return self.reviewer_allowed
