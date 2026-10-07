"""Explicit additive migration for the Python-owned pgvector retrieval database."""

import asyncio

from adapters.repositories.curriculum_postgres import CurriculumPostgres
from infrastructure.config import Settings


async def migrate() -> None:
    settings = Settings()
    if not settings.curriculum_database_url:
        raise SystemExit("RAG_CURRICULUM_DATABASE_URL is required")
    await CurriculumPostgres(settings.curriculum_database_url).migrate()


if __name__ == "__main__":
    asyncio.run(migrate())
