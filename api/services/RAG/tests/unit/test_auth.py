import asyncio

from fastapi import HTTPException

from infrastructure.auth import get_auth_context
from infrastructure.auth_client import AuthenticatedUser, AuthServiceError


class FakeAuthServiceClient:
    def __init__(
        self,
        user: AuthenticatedUser | None = None,
        error: Exception | None = None,
    ) -> None:
        self.user = user
        self.error = error
        self.token: str | None = None

    async def authenticate(self, bearer_token: str) -> AuthenticatedUser:
        self.token = bearer_token
        if self.error:
            raise self.error
        if self.user is None:
            raise AuthServiceError("missing fake user")
        return self.user


def test_get_auth_context_rejects_missing_authorization_header() -> None:
    try:
        asyncio.run(get_auth_context(authorization=None, auth_service=FakeAuthServiceClient()))
    except HTTPException as exc:
        assert exc.status_code == 401
        assert exc.detail == "Missing Authorization header"
    else:
        raise AssertionError("expected HTTPException")


def test_get_auth_context_rejects_invalid_bearer_format() -> None:
    try:
        asyncio.run(
            get_auth_context(
                authorization="Basic abc",
                auth_service=FakeAuthServiceClient(),
            )
        )
    except HTTPException as exc:
        assert exc.status_code == 401
        assert exc.detail == "Invalid bearer token"
    else:
        raise AssertionError("expected HTTPException")


def test_get_auth_context_returns_authenticated_user() -> None:
    auth_service = FakeAuthServiceClient(
        AuthenticatedUser(
            user_id="user-123",
            subject="user-123",
            roles=["SCHOOL_ADMIN"],
            claims={"sub": "user-123", "roles": ["SCHOOL_ADMIN"]},
        )
    )

    context = asyncio.run(
        get_auth_context(
            authorization="Bearer token-123",
            auth_service=auth_service,
        )
    )

    assert context.subject == "user-123"
    assert context.roles == ["SCHOOL_ADMIN"]
    assert context.token == "token-123"
    assert auth_service.token == "token-123"


def test_get_auth_context_rejects_auth_service_failures() -> None:
    auth_service = FakeAuthServiceClient(error=AuthServiceError("nope"))

    try:
        asyncio.run(get_auth_context(authorization="Bearer token-123", auth_service=auth_service))
    except HTTPException as exc:
        assert exc.status_code == 401
        assert exc.detail == "Invalid bearer token"
    else:
        raise AssertionError("expected HTTPException")
