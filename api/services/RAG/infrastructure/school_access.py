import httpx

from domain.entities.curriculum import CurriculumError, CurriculumScope


class SchoolAccessClient:
    def __init__(self, base_url: str, transport: httpx.AsyncBaseTransport | None = None) -> None:
        self.client = httpx.AsyncClient(
            base_url=base_url.rstrip("/"), timeout=5, transport=transport
        )

    async def authorize(self, token: str | None, scope: CurriculumScope, manage: bool) -> None:
        """Use Go's fresh membership/role checks, never token role labels or model text."""
        if not token:
            raise CurriculumError(401, "Curriculum access requires an authenticated user")
        headers = {"Authorization": f"Bearer {token}"}
        try:
            response = await self.client.get(
                f"/v1/schools/{scope.school_id}/memberships", headers=headers
            )
            if response.status_code == 200:
                # Go authorizes current SUPER_ADMIN or delegated school management here.
                departments = await self.client.get(
                    f"/v1/schools/{scope.school_id}/departments", headers=headers
                )
                departments.raise_for_status()
                if any(
                    row["id"] == str(scope.department_id)
                    for row in departments.json()["data"]
                ):
                    return
                raise CurriculumError(403, "Department does not belong to this school")
            if response.status_code not in (401, 403, 404):
                raise CurriculumError(503, "School access service unavailable")
            if manage:
                raise CurriculumError(403, "Curriculum review requires school management access")
            response = await self.client.get("/v1/me/school-memberships", headers=headers)
            if response.status_code in (401, 403):
                raise CurriculumError(403, "Current school access is required")
            response.raise_for_status()
            if any(
                row["attributes"].get("active") is True
                and row["attributes"].get("schoolID") == str(scope.school_id)
                and row["attributes"].get("departmentID") == str(scope.department_id)
                and row["attributes"].get("capability")
                in {"teach", "review_lessons", "approve_lessons"}
                for row in response.json()["data"]
            ):
                return
            raise CurriculumError(403, "Current department capability is required")
        except (httpx.HTTPError, KeyError, TypeError, ValueError) as exc:
            raise CurriculumError(503, "School access service unavailable") from exc

    async def close(self) -> None:
        await self.client.aclose()
