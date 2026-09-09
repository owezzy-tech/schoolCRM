from fastapi import FastAPI

from adapters.controllers.admissions_query_controller import router as admissions_query_router
from adapters.controllers.health_controller import router as health_router
from adapters.controllers.ingest_controller import router as ingest_router
from adapters.controllers.query_controller import router as query_router
from infrastructure.lifespan import lifespan


def build_app() -> FastAPI:
    app = FastAPI(
        title="schoolCRM RAG Service",
        version="0.1.0",
        description="Document Intelligence service for schoolCRM.",
        lifespan=lifespan,
    )
    app.include_router(health_router)
    app.include_router(admissions_query_router)
    app.include_router(ingest_router)
    app.include_router(query_router)
    return app
