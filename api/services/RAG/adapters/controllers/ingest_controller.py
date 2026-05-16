from fastapi import APIRouter, Depends, File, Form, HTTPException, UploadFile, status

from infrastructure.auth import AuthContext, get_auth_context
from infrastructure.dependencies import (
    get_delete_document_use_case,
    get_ingest_document_use_case,
)
from use_cases.delete_document import DeleteDocumentUseCase
from use_cases.ingest_document import IngestDocumentCommand, IngestDocumentUseCase

router = APIRouter(prefix="/v1/rag", tags=["rag-documents"])


@router.post("/documents", status_code=status.HTTP_202_ACCEPTED)
async def ingest_document(
    title: str = Form(...),
    source: str = Form(...),
    file: UploadFile = File(...),
    auth: AuthContext = Depends(get_auth_context),
    use_case: IngestDocumentUseCase = Depends(get_ingest_document_use_case),
) -> dict[str, str | int]:
    payload = await file.read()
    if not payload:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Empty file upload")

    result = await use_case.execute(
        IngestDocumentCommand(
            title=title,
            source=source,
            uploaded_by=auth.subject,
            filename=file.filename or "upload.bin",
            content_type=file.content_type or "application/octet-stream",
            payload=payload,
        )
    )

    return {
        "document_id": result.document_id,
        "status": result.status,
        "chunk_count": result.chunk_count,
    }


@router.delete("/documents/{document_id}")
async def delete_document(
    document_id: str,
    _: AuthContext = Depends(get_auth_context),
    use_case: DeleteDocumentUseCase = Depends(get_delete_document_use_case),
) -> dict[str, str | bool]:
    result = await use_case.execute(document_id)
    if not result.deleted:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Document not found")

    return {"document_id": result.document_id, "deleted": result.deleted}
