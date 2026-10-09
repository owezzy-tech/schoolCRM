from dataclasses import asdict
from datetime import UTC, datetime
from typing import Annotated

from pydantic import BaseModel, ConfigDict, Field, ValidationError

from domain.entities.curriculum import CurriculumError, CurriculumEvidence
from domain.entities.lesson_generation import GeneratedLesson, LessonGenerationRequest

Text = Annotated[str, Field(min_length=1, max_length=2000)]


class Activity(BaseModel):
    model_config = ConfigDict(extra="forbid", str_strip_whitespace=True, strict=True)
    minutes: int = Field(ge=1, le=180)
    title: str = Field(min_length=1, max_length=200)
    detail: Text


class ModelLesson(BaseModel):
    model_config = ConfigDict(extra="forbid", str_strip_whitespace=True, strict=True)
    title: str = Field(min_length=1, max_length=200)
    objectives: list[Text] = Field(min_length=1, max_length=20)
    prerequisites: list[Text] = Field(max_length=20)
    materials: list[Text] = Field(max_length=20)
    activities: list[Activity] = Field(min_length=1, max_length=60)
    differentiation: Text
    assessment: Text
    citation_ids: list[str] = Field(min_length=1, max_length=20)


def evidence_id(evidence: CurriculumEvidence) -> str:
    return f"{evidence.source.id}:{evidence.passage.ordinal}"


def citation(item: CurriculumEvidence) -> dict:
    source = item.source
    return {
        "sourceID": str(source.id),
        "indexRevision": str(item.index_revision),
        "chunkOrdinal": item.passage.ordinal,
        "page": item.passage.page,
        "title": source.title,
        "authority": source.authority,
        "sourceURL": source.source_url,
        "sourceSHA256": source.sha256,
        "framework": source.scope.framework,
        "stage": source.scope.stage,
        "subject": source.scope.subject,
        "revision": source.scope.revision,
        "embeddingModel": item.model_identity,
    }


def evidence_citation(item: CurriculumEvidence) -> dict:
    """Retrieved evidence for display, keyed by the ID the model may cite."""
    return {"id": evidence_id(item), **citation(item), "passage": item.passage.text}


def assemble_lesson(
    raw_json: str,
    request: LessonGenerationRequest,
    evidence: list[CurriculumEvidence],
    completion_id: str,
) -> GeneratedLesson:
    try:
        parsed = ModelLesson.model_validate_json(raw_json)
    except ValidationError:
        raise CurriculumError(422, "DeepSeek returned an invalid structured lesson") from None
    if sum(a.minutes for a in parsed.activities) != request.duration_minutes:
        raise CurriculumError(422, "Generated activity times do not total the requested duration")
    available = {evidence_id(item): item for item in evidence}
    if any(identifier not in available for identifier in parsed.citation_ids):
        raise CurriculumError(
            422, "Generated lesson refers to evidence outside the retrieved scope"
        )
    citations = [citation(available[i]) for i in dict.fromkeys(parsed.citation_ids)]
    content = parsed.model_dump(exclude={"title", "citation_ids"})
    content.update(
        {
            "schemaVersion": 1,
            "framework": request.scope.framework,
            "stage": request.scope.stage,
            "subject": request.scope.subject,
            "curriculumRevision": request.scope.revision,
            "durationMinutes": request.duration_minutes,
            "citations": citations,
            "generation": {
                "provider": "deepseek",
                "modelID": "deepseek-flash",
                "modelVersion": "DeepSeek-V4.1-Flash",
                "completionID": completion_id,
                "generatedAt": datetime.now(UTC).isoformat(),
            },
        }
    )
    return GeneratedLesson(parsed.title, content)


def model_evidence(evidence: list[CurriculumEvidence]) -> list[dict]:
    return [
        {
            "id": evidence_id(e),
            "scope": asdict(e.source.scope),
            "title": e.source.title,
            "page": e.passage.page,
            "text": e.passage.text,
        }
        for e in evidence
    ]
