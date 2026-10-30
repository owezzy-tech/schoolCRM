import json
from dataclasses import replace
from uuid import uuid4

from adapters.llm.lesson_schema import assemble_lesson
from domain.entities.curriculum import (
    CurriculumError,
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


class CurriculumFixture:
    def __init__(self, evidence):
        self.evidence = evidence
        self.withdrawn = False

    async def search(self, token, scope, question, limit):
        return [self.evidence] if scope == self.evidence.source.scope else []

    async def source(self, token, source_id):
        return (
            replace(self.evidence.source, status="withdrawn")
            if self.withdrawn
            else self.evidence.source
        )


class ModelFixture:
    def __init__(self, output):
        self.output = output
        self.calls = 0

    async def generate(self, request, evidence):
        self.calls += 1
        return assemble_lesson(json.dumps(self.output), request, evidence, "fixture-completion")

    async def close(self):
        pass


class GoFixture:
    def __init__(self):
        self.receipts = {}
        self.calls = 0
        self.fail_once = False
        self.revoked = False
        self.tokens = []

    async def authorize(self, token, school_id, department_id):
        self.tokens.append(token)
        if self.revoked:
            raise CurriculumError(403, "Teaching permission revoked")

    async def create(self, token, scope, request_id, draft):
        await self.authorize(token, scope.school_id, scope.department_id)
        self.calls += 1
        if request_id not in self.receipts:
            self.receipts[request_id] = {"planID": str(uuid4()), "version": 1, "title": draft.title}
        if self.fail_once:
            self.fail_once = False
            raise CurriculumError(503, "Ambiguous Go response after commit")
        return self.receipts[request_id]

    async def close(self):
        pass
