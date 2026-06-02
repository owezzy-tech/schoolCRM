from typing import Any

from fastapi import APIRouter, Depends
from pydantic import BaseModel, Field

from domain.entities.query import CitationSource, QueryResult
from domain.ports.application_ownership import IApplicationOwnershipChecker
from infrastructure.admissions_auth import (
    guard_applicant_query,
    guard_reviewer_query,
    guard_staff_query,
)
from infrastructure.auth import AuthContext, get_auth_context
from infrastructure.dependencies import (
    get_admissions_query_use_case,
    get_application_ownership_checker,
)
from use_cases.admissions_query import AdmissionsQueryCommand, AdmissionsQueryUseCase


class AdmissionsQueryRequestDTO(BaseModel):
    question: str = Field(min_length=1)
    top_k: int = Field(default=3, ge=1, le=10)


router = APIRouter(prefix="/v1/rag/admissions", tags=["admissions-rag"])


@router.post("/staff/query")
async def query_staff_policy(
    dto: AdmissionsQueryRequestDTO,
    auth: AuthContext = Depends(get_auth_context),
    use_case: AdmissionsQueryUseCase = Depends(get_admissions_query_use_case),
) -> dict[str, Any]:
    context = guard_staff_query(auth)
    result = await use_case.execute(
        AdmissionsQueryCommand(question=dto.question, top_k=dto.top_k, context=context)
    )
    return _jsonapi_result(result)


@router.post("/applicant/applications/{application_id}/query")
async def query_applicant_application(
    application_id: str,
    dto: AdmissionsQueryRequestDTO,
    auth: AuthContext = Depends(get_auth_context),
    ownership_checker: IApplicationOwnershipChecker = Depends(get_application_ownership_checker),
    use_case: AdmissionsQueryUseCase = Depends(get_admissions_query_use_case),
) -> dict[str, Any]:
    context = await guard_applicant_query(auth, application_id, ownership_checker)
    result = await use_case.execute(
        AdmissionsQueryCommand(question=dto.question, top_k=dto.top_k, context=context)
    )
    return _jsonapi_result(result)


@router.post("/reviewer/applications/{application_id}/query")
async def query_reviewer_application(
    application_id: str,
    dto: AdmissionsQueryRequestDTO,
    auth: AuthContext = Depends(get_auth_context),
    ownership_checker: IApplicationOwnershipChecker = Depends(get_application_ownership_checker),
    use_case: AdmissionsQueryUseCase = Depends(get_admissions_query_use_case),
) -> dict[str, Any]:
    context = await guard_reviewer_query(auth, application_id, ownership_checker)
    result = await use_case.execute(
        AdmissionsQueryCommand(question=dto.question, top_k=dto.top_k, context=context)
    )
    return _jsonapi_result(result)


def _jsonapi_result(result: QueryResult) -> dict[str, Any]:
    return {
        "jsonapi": {"version": "1.1"},
        "data": {
            "type": "rag-result",
            "attributes": {
                "answer": result.answer,
                "documentIds": result.document_ids,
                "snippets": result.snippets,
                "citations": [
                    _citation_attributes(citation) for citation in result.citations or []
                ],
            },
        },
    }


def _citation_attributes(citation: CitationSource) -> dict[str, Any]:
    return {
        "documentId": citation.document_id,
        "chunkIndex": citation.chunk_index,
        "excerpt": citation.excerpt,
        "collection": citation.collection,
        "score": citation.score,
        "sourceType": citation.source_type,
        "title": citation.title,
        "path": citation.path or [],
    }
