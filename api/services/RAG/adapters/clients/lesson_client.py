"""Delegated Go business tools. Python never writes lesson tables directly."""

from uuid import UUID

import httpx

from domain.entities.curriculum import CurriculumError, CurriculumScope
from domain.entities.lesson_generation import GeneratedLesson


class GoLessonClient:
    def __init__(self, base_url: str, transport: httpx.AsyncBaseTransport | None = None):
        self.client = httpx.AsyncClient(
            base_url=base_url.rstrip("/"), timeout=30, transport=transport
        )

    async def authorize(self, token: str, school_id: UUID, department_id: UUID) -> None:
        response = await self._request("GET", "/v1/me/school-memberships", token)
        try:
            rows = response.json()["data"]
            if not isinstance(rows, list):
                raise TypeError("Invalid membership collection")
            allowed = any(
                row["attributes"].get("active") is True
                and row["attributes"].get("schoolID") == str(school_id)
                and row["attributes"].get("departmentID") == str(department_id)
                and row["attributes"].get("capability") == "teach"
                for row in rows
            )
        except (ValueError, KeyError, TypeError, AttributeError):
            raise CurriculumError(503, "Go membership response is invalid") from None
        if not allowed:
            raise CurriculumError(
                403, "Lesson generation requires current department teaching access"
            )

    async def create(
        self, token: str, scope: CurriculumScope, request_id: UUID, draft: GeneratedLesson
    ) -> dict:
        response = await self._request(
            "POST",
            f"/v1/schools/{scope.school_id}/departments/{scope.department_id}/lessons",
            token,
            {
                "title": draft.title,
                "content": draft.content,
                "changeSummary": "Generated from reviewed curriculum evidence",
                "requestID": str(request_id),
            },
        )
        try:
            resource = response.json()["data"]
            plan_id = UUID(resource["id"])
            if plan_id.int == 0:
                raise ValueError("Invalid plan identity")
        except (ValueError, KeyError, TypeError):
            raise CurriculumError(503, "Go lesson creation response is invalid") from None
        return {
            "planID": str(plan_id),
            "version": 1,
            "title": draft.title,
            "creationStatus": "draft",
            "citations": draft.content["citations"],
            "generation": draft.content["generation"],
        }

    async def _request(
        self, method: str, path: str, token: str, body: dict | None = None
    ) -> httpx.Response:
        if not token:
            raise CurriculumError(401, "An authenticated user is required")
        try:
            response = await self.client.request(
                method, path, headers={"Authorization": f"Bearer {token}"}, json=body
            )
            if response.status_code in (400, 401, 403, 404, 409, 422):
                detail = (
                    response.json().get("errors", [{}])[0].get("detail", "Lesson command rejected")
                )
                raise CurriculumError(response.status_code, detail)
            response.raise_for_status()
            return response
        except (httpx.HTTPError, ValueError, KeyError, TypeError, IndexError):
            raise CurriculumError(503, "Go lesson service unavailable") from None

    async def close(self) -> None:
        await self.client.aclose()
