from functools import lru_cache

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    service_name: str = "rag"
    api_host: str = "0.0.0.0"
    api_port: int = 7000
    allow_anonymous: bool = False
    auth_service_url: str = "http://auth-service:6000"
    auth_request_timeout_seconds: float = 5.0
    file_storage_dir: str = "./var/files"
    admissions_answer_provider: str = "stub"
    ollama_base_url: str = "http://localhost:11434"
    ollama_model: str = "nemotron-3-super:cloud"
    ollama_timeout_seconds: float = 60.0
    graph_retriever: str = "memory"
    neo4j_uri: str = "bolt://localhost:7687"
    neo4j_username: str = "neo4j"
    neo4j_password: str = "password"
    langsmith_enabled: bool = False
    langsmith_api_key: str | None = None
    langsmith_project: str = "schoolcrm-rag"
    langsmith_endpoint: str | None = None

    model_config = SettingsConfigDict(
        env_prefix="RAG_",
        env_file=(".env", "api/services/RAG/.env"),
        env_file_encoding="utf-8",
        extra="ignore",
    )


@lru_cache
def get_settings() -> Settings:
    return Settings()
