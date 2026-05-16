from dataclasses import dataclass

from fastapi import Header, HTTPException, status

from infrastructure.config import get_settings


@dataclass(slots=True)
class AuthContext:
    subject: str
    token: str | None


async def get_auth_context(authorization: str | None = Header(default=None)) -> AuthContext:
    settings = get_settings()
    if not authorization:
        if settings.allow_anonymous:
            return AuthContext(subject="anonymous", token=None)
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Missing Authorization header",
        )

    scheme, _, token = authorization.partition(" ")
    if scheme.lower() != "bearer" or not token:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid bearer token")

    return AuthContext(subject="authenticated-user", token=token)
