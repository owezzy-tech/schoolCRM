from fastapi import APIRouter, Depends
from pydantic import BaseModel, Field

from infrastructure.auth import AuthContext, get_auth_context
from infrastructure.dependencies import get_query_documents_use_case
from use_cases.query_documents import QueryDocumentsCommand, QueryDocumentsUseCase


class QueryRequestDTO(BaseModel):
    question: str = Field(min_length=1)
    top_k: int = Field(default=3, ge=1, le=10)


router = APIRouter(prefix="/v1/rag", tags=["rag-query"])


@router.post("/query")
async def query_documents(
    dto: QueryRequestDTO,
    auth: AuthContext = Depends(get_auth_context),
    use_case: QueryDocumentsUseCase = Depends(get_query_documents_use_case),
) -> dict[str, str | list[str]]:
    result = await use_case.execute(
        QueryDocumentsCommand(question=dto.question, requested_by=auth.subject, top_k=dto.top_k)
    )

    return {
        "answer": result.answer,
        "document_ids": result.document_ids,
        "snippets": result.snippets,
    }
