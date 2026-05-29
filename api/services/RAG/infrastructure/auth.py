from dataclasses import dataclass

import httpx
from fastapi import Depends, Header, HTTPException, status

from infrastructure.auth_client import AuthServiceClient, AuthServiceError
from infrastructure.config import get_settings
from infrastructure.dependencies import get_auth_service_client


@dataclass(slots=True)
class AuthContext:
    subject: str
    token: str | None
    roles: list[str]


async def get_auth_context(
    authorization: str | None = Header(default=None),
    auth_service: AuthServiceClient = Depends(get_auth_service_client),
) -> AuthContext:
    settings = get_settings()
    if not authorization:
        if settings.allow_anonymous:
            return AuthContext(subject="anonymous", token=None, roles=[])
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Missing Authorization header",
        )

    scheme, _, token = authorization.partition(" ")
    if scheme.lower() != "bearer" or not token:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid bearer token")

    try:
        user = await auth_service.authenticate(token)
    except (AuthServiceError, httpx.HTTPError):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid bearer token",
        ) from None

    return AuthContext(subject=user.subject, token=token, roles=user.roles)
