from contextlib import asynccontextmanager

from fastapi import FastAPI

from infrastructure.config import Settings
from infrastructure.dependencies import Container, build_container
from infrastructure.observability import configure_logging
from infrastructure.tracing import configure_langsmith


@asynccontextmanager
async def lifespan(app: FastAPI):
    settings = Settings()
    configure_logging(settings.service_name)
    configure_langsmith(settings)
    app.state.settings = settings
    container: Container = build_container(settings)
    app.state.container = container
    try:
        if container.lesson_workflow is not None:
            await container.lesson_workflow.start()
        yield
    finally:
        await container.close()
