from uuid import uuid4

from domain.entities.curriculum import (
    CurriculumEvidence,
    CurriculumPassage,
    CurriculumScope,
    CurriculumSource,
)
from domain.entities.lesson_generation import LessonGenerationRequest


def generation_fixture():
    scope = CurriculumScope(uuid4(), uuid4(), "kenya-cbc", "grade-1", "english", "2024")
    source = CurriculumSource(
        uuid4(),
        scope,
        "English",
        "School",
        "KICD",
        "https://kicd.ac.ke/english.pdf",
        "a" * 64,
        "english.pdf",
        17,
        17,
        "uploader",
        status="approved",
    )
    evidence = CurriculumEvidence(
        source,
        CurriculumPassage(17, 0, "Identify classroom objects"),
        0.9,
        "ollama:bge-m3:pinned",
        source.id,
    )
    request = LessonGenerationRequest(uuid4(), scope, "School vocabulary", 30)
    output = {
        "title": "Classroom objects",
        "objectives": ["Name a desk and chair"],
        "prerequisites": [],
        "materials": ["Picture cards"],
        "activities": [
            {"minutes": 30, "title": "Practise words", "detail": "Name classroom objects"}
        ],
        "differentiation": "Use picture prompts",
        "assessment": "Name three objects",
        "citation_ids": [f"{source.id}:0"],
    }
    return request, evidence, output
