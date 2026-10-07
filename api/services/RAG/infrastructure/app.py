import psycopg
from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError
from starlette.exceptions import HTTPException

from adapters.controllers.admissions_query_controller import router as admissions_query_router
from adapters.controllers.curriculum_controller import (
    curriculum_error_handler,
    curriculum_http_error_handler,
    curriculum_storage_error_handler,
    curriculum_validation_error_handler,
)
from adapters.controllers.curriculum_controller import (
    router as curriculum_router,
)
from adapters.controllers.health_controller import router as health_router
from adapters.controllers.ingest_controller import router as ingest_router
from adapters.controllers.query_controller import router as query_router
from domain.entities.curriculum import CurriculumError
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
    app.include_router(curriculum_router)
    app.add_exception_handler(CurriculumError, curriculum_error_handler)
    app.add_exception_handler(psycopg.Error, curriculum_storage_error_handler)
    app.add_exception_handler(HTTPException, curriculum_http_error_handler)
    app.add_exception_handler(RequestValidationError, curriculum_validation_error_handler)
    return app
