import asyncio

import httpx

from infrastructure.auth_client import AuthServiceClient


def test_auth_service_client_accepts_jsonapi_authenticate_response() -> None:
    payload = {
        "jsonapi": {"version": "1.1"},
        "data": {
            "type": "authenticateresp",
            "attributes": {
                "userID": "user-1",
                "claims": {
                    "sub": "user-1",
                    "roles": ["SCHOOL_ADMIN", "ADMISSIONS_ADMIN"],
                },
            },
        },
    }

    user = asyncio.run(_authenticate(payload))

    assert user.user_id == "user-1"
    assert user.subject == "user-1"
    assert user.roles == ["SCHOOL_ADMIN", "ADMISSIONS_ADMIN"]


def test_auth_service_client_accepts_legacy_authenticate_response() -> None:
    payload = {
        "UserID": "user-1",
        "Claims": {
            "sub": "user-1",
            "roles": ["SCHOOL_ADMIN"],
        },
    }

    user = asyncio.run(_authenticate(payload))

    assert user.user_id == "user-1"
    assert user.subject == "user-1"
    assert user.roles == ["SCHOOL_ADMIN"]


async def _authenticate(payload: dict[str, object]):
    async def handler(request: httpx.Request) -> httpx.Response:
        assert request.url.path == "/v1/auth/authenticate"
        assert request.headers["Authorization"] == "Bearer token"
        return httpx.Response(200, json=payload)

    transport = httpx.MockTransport(handler)
    client = AuthServiceClient(
        base_url="http://auth-service:6000",
        timeout_seconds=5,
        transport=transport,
    )

    try:
        return await client.authenticate("token")
    finally:
        await client.close()
