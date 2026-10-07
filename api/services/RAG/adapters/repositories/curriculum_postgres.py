import json
from dataclasses import asdict
from pathlib import Path
from uuid import UUID

import psycopg
from psycopg.rows import dict_row
from psycopg.types.json import Jsonb

from domain.entities.curriculum import (
    REVIEW_TRANSITIONS,
    CurriculumError,
    CurriculumEvidence,
    CurriculumPassage,
    CurriculumScope,
    CurriculumSource,
)


def _json(value: object) -> Jsonb:
    return Jsonb(value, dumps=lambda v: json.dumps(v, default=str))


def _source(row: dict) -> CurriculumSource:
    metadata = row["metadata"].copy()
    scope = metadata.pop("scope")
    scope["school_id"] = UUID(scope["school_id"])
    scope["department_id"] = UUID(scope["department_id"])
    metadata["id"] = UUID(metadata["id"])
    return CurriculumSource(
        **metadata,
        scope=CurriculumScope(**scope),
        status=row["status"],
        reviewed_by=row["reviewed_by"],
        review_note=row["review_note"],
        model_identity=row["model_identity"],
    )


class CurriculumPostgres:
    def __init__(self, dsn: str) -> None:
        self.dsn = dsn

    async def migrate(self) -> None:
        schema = Path(__file__).with_name("curriculum.sql").read_text()
        async with await psycopg.AsyncConnection.connect(self.dsn) as conn:
            await conn.execute(schema)

    async def create(
        self, source: CurriculumSource, content: bytes, passages: list[CurriculumPassage]
    ) -> None:
        metadata = asdict(source)
        for key in ("status", "reviewed_by", "review_note", "model_identity"):
            metadata.pop(key)
        scope = source.scope
        try:
            async with await psycopg.AsyncConnection.connect(self.dsn) as conn:
                await conn.execute(
                    """INSERT INTO rag_curriculum_sources
                    (id, school_id, department_id, framework, stage, subject, revision,
                     sha256, metadata, original)
                    VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)""",
                    (
                        source.id,
                        scope.school_id,
                        scope.department_id,
                        scope.framework,
                        scope.stage,
                        scope.subject,
                        scope.revision,
                        source.sha256,
                        _json(metadata),
                        content,
                    ),
                )
                async with conn.cursor() as cursor:
                    await cursor.executemany(
                        """INSERT INTO rag_curriculum_chunks (source_id, ordinal, page, content)
                        VALUES (%s,%s,%s,%s)""",
                        [(source.id, p.ordinal, p.page, p.text) for p in passages],
                    )
                await self._audit(conn, source.id, source.uploaded_by, "uploaded", "pending")
        except psycopg.errors.UniqueViolation as exc:
            raise CurriculumError(
                409, "This original and page range already exist in this scope"
            ) from exc

    async def get(self, source_id: UUID) -> CurriculumSource:
        async with await psycopg.AsyncConnection.connect(self.dsn, row_factory=dict_row) as conn:
            cursor = await conn.execute(
                "SELECT * FROM rag_curriculum_sources WHERE id = %s", (source_id,)
            )
            row = await cursor.fetchone()
            if not row:
                raise CurriculumError(404, "Curriculum source not found")
            return _source(row)

    async def passages(self, source_id: UUID) -> list[CurriculumPassage]:
        async with await psycopg.AsyncConnection.connect(self.dsn, row_factory=dict_row) as conn:
            cursor = await conn.execute(
                """SELECT page, ordinal, content FROM rag_curriculum_chunks
                WHERE source_id = %s ORDER BY ordinal""",
                (source_id,),
            )
            return [
                CurriculumPassage(r["page"], r["ordinal"], r["content"])
                for r in await cursor.fetchall()
            ]

    async def review(
        self,
        source_id: UUID,
        actor: str,
        expected: str,
        decision: str,
        note: str,
        model_identity: str | None,
        vectors: list[list[float]],
    ) -> CurriculumSource:
        async with await psycopg.AsyncConnection.connect(self.dsn, row_factory=dict_row) as conn:
            cursor = await conn.execute(
                "SELECT * FROM rag_curriculum_sources WHERE id = %s FOR UPDATE", (source_id,)
            )
            row = await cursor.fetchone()
            if not row:
                raise CurriculumError(404, "Curriculum source not found")
            if row["status"] != expected:
                raise CurriculumError(409, "Source state changed; reload before reviewing")
            if (expected, decision) not in REVIEW_TRANSITIONS:
                raise CurriculumError(409, "Invalid curriculum review transition")
            if decision == "approved":
                cursor = await conn.execute(
                    "SELECT count(*) AS n FROM rag_curriculum_chunks WHERE source_id = %s",
                    (source_id,),
                )
                count = (await cursor.fetchone())["n"]
                if not model_identity or count != len(vectors):
                    raise CurriculumError(503, "Incomplete curriculum index")
                async with conn.cursor() as cursor:
                    await cursor.executemany(
                        """UPDATE rag_curriculum_chunks SET embedding = %s::vector
                        WHERE source_id = %s AND ordinal = %s""",
                        [(str(vector), source_id, i) for i, vector in enumerate(vectors)],
                    )
            await conn.execute(
                """UPDATE rag_curriculum_sources SET status=%s, reviewed_by=%s, review_note=%s,
                model_identity=COALESCE(%s, model_identity) WHERE id=%s""",
                (decision, actor, note, model_identity, source_id),
            )
            await self._audit(conn, source_id, actor, decision, note)
            cursor = await conn.execute(
                "SELECT * FROM rag_curriculum_sources WHERE id = %s", (source_id,)
            )
            return _source(await cursor.fetchone())

    async def search(
        self,
        scope: CurriculumScope,
        model_identity: str,
        vector: list[float],
        limit: int,
        minimum_score: float,
    ) -> list[CurriculumEvidence]:
        # Scope predicates run inside PostgreSQL, before any evidence reaches a model.
        async with await psycopg.AsyncConnection.connect(self.dsn, row_factory=dict_row) as conn:
            cursor = await conn.execute(
                """SELECT s.*, c.page, c.ordinal, c.content,
                1 - (c.embedding <=> %s::vector) AS score
                FROM rag_curriculum_sources s JOIN rag_curriculum_chunks c ON c.source_id=s.id
                WHERE s.school_id=%s AND s.department_id=%s AND s.framework=%s
                AND s.stage=%s AND s.subject=%s AND s.revision=%s
                AND s.status='approved' AND s.model_identity=%s
                AND c.embedding IS NOT NULL AND 1 - (c.embedding <=> %s::vector) >= %s
                ORDER BY c.embedding <=> %s::vector, s.id, c.ordinal LIMIT %s""",
                (
                    str(vector),
                    scope.school_id,
                    scope.department_id,
                    scope.framework,
                    scope.stage,
                    scope.subject,
                    scope.revision,
                    model_identity,
                    str(vector),
                    minimum_score,
                    str(vector),
                    limit,
                ),
            )
            return [
                CurriculumEvidence(
                    _source(r),
                    CurriculumPassage(r["page"], r["ordinal"], r["content"]),
                    r["score"],
                    r["model_identity"],
                    r["id"],
                )
                for r in await cursor.fetchall()
            ]

    async def original(self, source_id: UUID) -> bytes:
        async with await psycopg.AsyncConnection.connect(self.dsn) as conn:
            cursor = await conn.execute(
                "SELECT original FROM rag_curriculum_sources WHERE id=%s", (source_id,)
            )
            row = await cursor.fetchone()
            if not row:
                raise CurriculumError(404, "Curriculum source not found")
            return bytes(row[0])

    async def _audit(
        self, conn: psycopg.AsyncConnection, source_id: UUID, actor: str, action: str, note: str
    ) -> None:
        await conn.execute(
            """INSERT INTO rag_curriculum_audit (source_id, actor, action, note)
            VALUES (%s,%s,%s,%s)""",
            (source_id, actor, action, note),
        )
