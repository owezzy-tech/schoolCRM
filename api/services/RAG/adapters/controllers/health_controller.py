from fastapi import APIRouter


router = APIRouter(tags=["health"])


@router.get("/v1/liveness")
async def liveness() -> dict[str, str]:
    return {"status": "ok"}


@router.get("/v1/readiness")
async def readiness() -> dict[str, str]:
    return {"status": "ok"}
