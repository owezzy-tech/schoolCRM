from contextlib import asynccontextmanager
from pathlib import Path
from uuid import UUID

import psycopg
from psycopg.rows import dict_row
from psycopg.types.json import Jsonb

from domain.entities.curriculum import CurriculumError


class LessonRuns:
    def __init__(self, dsn: str):
        self.dsn = dsn

    async def migrate(self) -> None:
        async with await psycopg.AsyncConnection.connect(self.dsn) as conn:
            await conn.execute(Path(__file__).with_name("lesson_runs.sql").read_text())

    @asynccontextmanager
    async def command(self, actor: UUID, request_id: UUID, input_hash: str, request: dict):
        # A session lock survives the checkpoint's independent transactions.
        conn = await psycopg.AsyncConnection.connect(
            self.dsn, autocommit=True, row_factory=dict_row
        )
        try:
            key = f"rag-lesson-generation:{actor}:{request_id}"
            await conn.execute("SELECT pg_advisory_lock(hashtextextended(%s,0))", (key,))
            await conn.execute(
                """INSERT INTO rag_lesson_runs(actor_id,request_id,input_hash,request)
                VALUES (%s,%s,%s,%s) ON CONFLICT DO NOTHING""",
                (actor, request_id, input_hash, Jsonb(request)),
            )
            row = await self._row(conn, actor, request_id)
            if row["input_hash"] != input_hash:
                raise CurriculumError(
                    409, "request_id already used with different generation input"
                )
            yield conn, row
        finally:
            # Ending the session releases its lock. close() does no awaited I/O,
            # so a cancelled stream (client disconnect) cannot skip the release.
            await conn.close()

    async def finish(
        self, conn: psycopg.AsyncConnection, actor: UUID, request_id: UUID, result: dict
    ) -> None:
        await conn.execute(
            "UPDATE rag_lesson_runs SET result=%s WHERE actor_id=%s AND request_id=%s",
            (Jsonb(result), actor, request_id),
        )

    async def scope(self, actor: UUID, request_id: UUID) -> tuple[UUID, UUID]:
        # Resolve only immutable scope metadata before checking current authority.
        async with await psycopg.AsyncConnection.connect(self.dsn) as conn:
            cursor = await conn.execute(
                "SELECT request->'scope'->>'school_id', "
                "request->'scope'->>'department_id' FROM rag_lesson_runs "
                "WHERE actor_id=%s AND request_id=%s",
                (actor, request_id),
            )
            row = await cursor.fetchone()
            if not row:
                raise CurriculumError(404, "Generation request not found")
            return UUID(row[0]), UUID(row[1])

    async def get(self, actor: UUID, request_id: UUID) -> dict:
        async with await psycopg.AsyncConnection.connect(self.dsn, row_factory=dict_row) as conn:
            return await self._row(conn, actor, request_id)

    async def latest(self, actor: UUID, school_id: UUID, department_id: UUID) -> list[dict]:
        async with await psycopg.AsyncConnection.connect(self.dsn, row_factory=dict_row) as conn:
            cursor = await conn.execute(
                """SELECT request,result FROM rag_lesson_runs
                WHERE actor_id=%s AND request->'scope'->>'school_id'=%s
                AND request->'scope'->>'department_id'=%s
                ORDER BY created_at DESC, request_id LIMIT 50""",
                (actor, str(school_id), str(department_id)),
            )
            return await cursor.fetchall()

    async def _row(self, conn: psycopg.AsyncConnection, actor: UUID, request_id: UUID) -> dict:
        cursor = await conn.execute(
            "SELECT request,result,input_hash FROM rag_lesson_runs "
            "WHERE actor_id=%s AND request_id=%s",
            (actor, request_id),
        )
        row = await cursor.fetchone()
        if not row:
            raise CurriculumError(404, "Generation request not found")
        return row
