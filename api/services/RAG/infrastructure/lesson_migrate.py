import asyncio

from langgraph.checkpoint.postgres.aio import AsyncPostgresSaver

from adapters.repositories.lesson_runs import LessonRuns
from infrastructure.config import Settings


async def migrate() -> None:
    settings = Settings()
    if not settings.curriculum_database_url:
        raise SystemExit("RAG_CURRICULUM_DATABASE_URL is required")
    await LessonRuns(settings.curriculum_database_url).migrate()
    async with AsyncPostgresSaver.from_conn_string(settings.curriculum_database_url) as saver:
        await saver.setup()


if __name__ == "__main__":
    asyncio.run(migrate())
