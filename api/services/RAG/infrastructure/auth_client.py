from dataclasses import dataclass
from typing import Any

import httpx


class AuthServiceError(Exception):
    """Raised when the auth service rejects or cannot validate a token."""


@dataclass(slots=True)
class AuthenticatedUser:
    user_id: str
    subject: str
    roles: list[str]
    claims: dict[str, Any]


class AuthServiceClient:
    def __init__(self, base_url: str, timeout_seconds: float) -> None:
        self._client = httpx.AsyncClient(base_url=base_url.rstrip("/"), timeout=timeout_seconds)

    async def authenticate(self, bearer_token: str) -> AuthenticatedUser:
        response = await self._client.get(
            "/v1/auth/authenticate",
            headers={"Authorization": f"Bearer {bearer_token}"},
        )
        if response.status_code != 200:
            raise AuthServiceError("auth service rejected bearer token")

        payload = response.json()
        claims = payload.get("Claims") or payload.get("claims") or {}
        user_id = str(
            payload.get("UserID") or payload.get("userID") or payload.get("user_id") or ""
        )
        subject = str(claims.get("sub") or claims.get("subject") or user_id)
        roles = claims.get("roles") or payload.get("roles") or []

        if not user_id or not isinstance(roles, list):
            raise AuthServiceError("auth service returned an invalid authenticate response")

        return AuthenticatedUser(
            user_id=user_id,
            subject=subject,
            roles=[str(role) for role in roles],
            claims=claims,
        )

    async def close(self) -> None:
        await self._client.aclose()
