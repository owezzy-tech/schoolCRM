from dataclasses import asdict
from typing import Annotated, Literal
from uuid import UUID

import psycopg
from fastapi import APIRouter, Depends, Form, Request, Response, UploadFile
from fastapi.encoders import jsonable_encoder
from fastapi.exception_handlers import http_exception_handler, request_validation_exception_handler
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse
from pydantic import BaseModel, ConfigDict, Field, HttpUrl, ValidationError, model_validator
from starlette.exceptions import HTTPException

from domain.entities.curriculum import CurriculumError, CurriculumScope, CurriculumSource
from infrastructure.auth import AuthContext, get_auth_context
from use_cases.curriculum import CurriculumService


class JSONAPIResponse(JSONResponse):
    media_type = "application/vnd.api+json"


router = APIRouter(
    prefix="/v1/rag/curriculum", tags=["curriculum"], default_response_class=JSONAPIResponse
)
MAX_UPLOAD_BYTES = 10 * 1024 * 1024


class ScopeDTO(BaseModel):
    model_config = ConfigDict(extra="forbid", str_strip_whitespace=True)
    school_id: UUID
    department_id: UUID
    framework: Literal["kenya-cbc", "kenya-cbe", "kenya-8-4-4", "cambridge"]
    stage: str = Field(min_length=1, max_length=50)
    subject: str = Field(min_length=1, max_length=100)
    revision: str = Field(min_length=1, max_length=100)

    @model_validator(mode="after")
    def validate_scope(self) -> "ScopeDTO":
        if self.school_id.int == 0 or self.department_id.int == 0:
            raise ValueError("School and department must be nonzero UUIDs")
        self.subject = self.subject.lower()
        self.stage = self.stage.lower()
        if self.framework == "cambridge":
            allowed = {"early-years"} | {f"year-{n}" for n in range(1, 10)}
        elif self.framework == "kenya-8-4-4":
            allowed = {f"standard-{n}" for n in range(1, 9)} | {f"form-{n}" for n in range(1, 5)}
        else:
            allowed = {"playgroup", "pp1", "pp2"} | {f"grade-{n}" for n in range(1, 13)}
        if self.stage not in allowed:
            raise ValueError("Stage is outside the confirmed framework scope")
        return self

    def scope(self) -> CurriculumScope:
        return CurriculumScope(**{key: getattr(self, key) for key in ScopeDTO.model_fields})


class SourceDTO(ScopeDTO):
    title: str = Field(min_length=1, max_length=200)
    topic: str = Field(min_length=1, max_length=200)
    authority: str = Field(min_length=1, max_length=200)
    source_url: HttpUrl
    sha256: str = Field(pattern="^[a-f0-9]{64}$")
    first_page: int = Field(ge=1)
    last_page: int = Field(ge=1)


class ReviewDTO(BaseModel):
    model_config = ConfigDict(extra="forbid", str_strip_whitespace=True)
    expected_status: Literal["pending", "approved"]
    decision: Literal["approved", "withdrawn"]
    note: str = Field(min_length=1, max_length=2000)


class SearchDTO(ScopeDTO):
    question: str = Field(min_length=1, max_length=2000)
    limit: int = Field(default=5, ge=1, le=20)


def get_curriculum(request: Request) -> CurriculumService:
    service = request.app.state.container.curriculum
    if service is None:
        raise CurriculumError(503, "Curriculum storage has not been configured")
    return service


Service = Annotated[CurriculumService, Depends(get_curriculum)]
Actor = Annotated[AuthContext, Depends(get_auth_context)]


def document(source: CurriculumSource) -> dict:
    attributes = asdict(source)
    attributes.pop("id")
    return {
        "jsonapi": {"version": "1.1"},
        "data": {
            "type": "curriculum-source",
            "id": str(source.id),
            "attributes": jsonable_encoder(attributes),
        },
    }


@router.post("/sources", status_code=201)
async def upload_source(
    metadata: Annotated[str, Form(max_length=4000)],
    file: UploadFile,
    actor: Actor,
    service: Service,
) -> dict:
    try:
        request = SourceDTO.model_validate_json(metadata)
    except ValidationError as exc:
        raise CurriculumError(422, "Invalid curriculum metadata") from exc
    await service.access.authorize(actor.token, request.scope(), manage=False)
    try:
        content = await file.read(MAX_UPLOAD_BYTES + 1)
        if len(content) > MAX_UPLOAD_BYTES:
            raise CurriculumError(413, "Curriculum originals must not exceed 10 MiB")
        source = await service.ingest(
            actor.token,
            actor.subject,
            request.scope(),
            request.title,
            request.topic,
            request.authority,
            str(request.source_url),
            f"{request.sha256}.pdf",
            request.first_page,
            request.last_page,
            request.sha256,
            content,
        )
        return document(source)
    finally:
        await file.close()


@router.get("/sources/{source_id}")
async def source_details(source_id: UUID, actor: Actor, service: Service) -> dict:
    return document(await service.source(actor.token, source_id))


@router.get("/sources/{source_id}/original")
async def source_original(source_id: UUID, actor: Actor, service: Service) -> Response:
    await service.source(actor.token, source_id)
    content = await service.store.original(source_id)
    await service.source(actor.token, source_id)
    return Response(
        content,
        media_type="application/pdf",
        headers={
            "Content-Disposition": f'attachment; filename="{source_id}.pdf"',
            "Cache-Control": "private, no-store",
        },
    )


@router.post("/sources/{source_id}/review")
async def review_source(
    source_id: UUID, request: ReviewDTO, actor: Actor, service: Service
) -> dict:
    return document(
        await service.review(
            actor.token,
            actor.subject,
            source_id,
            request.expected_status,
            request.decision,
            request.note,
        )
    )


@router.post("/search")
async def search_curriculum(request: SearchDTO, actor: Actor, service: Service) -> dict:
    evidence = await service.search(actor.token, request.scope(), request.question, request.limit)
    return {
        "jsonapi": {"version": "1.1"},
        "data": [
            {
                "type": "curriculum-evidence",
                "id": f"{e.source.id}:{e.passage.ordinal}",
                "attributes": jsonable_encoder(asdict(e)),
            }
            for e in evidence
        ],
        "meta": {
            "abstained": not evidence,
            "reason": None if evidence else "No approved relevant curriculum evidence",
        },
    }


async def curriculum_error_handler(request: Request, exc: CurriculumError) -> Response:
    return JSONResponse(
        {
            "jsonapi": {"version": "1.1"},
            "errors": [{"status": str(exc.status_code), "detail": exc.detail}],
        },
        status_code=exc.status_code,
        media_type="application/vnd.api+json",
    )


async def curriculum_storage_error_handler(request: Request, exc: psycopg.Error) -> Response:
    return await curriculum_error_handler(
        request, CurriculumError(503, "Curriculum storage unavailable")
    )


async def curriculum_http_error_handler(request: Request, exc: HTTPException) -> Response:
    if request.url.path.startswith("/v1/rag/curriculum/"):
        return await curriculum_error_handler(
            request, CurriculumError(exc.status_code, str(exc.detail))
        )
    return await http_exception_handler(request, exc)


async def curriculum_validation_error_handler(
    request: Request, exc: RequestValidationError
) -> Response:
    if request.url.path.startswith("/v1/rag/curriculum/"):
        return await curriculum_error_handler(
            request, CurriculumError(422, "Invalid curriculum request")
        )
    return await request_validation_exception_handler(request, exc)
