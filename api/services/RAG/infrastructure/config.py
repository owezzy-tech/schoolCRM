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

    model_config = SettingsConfigDict(env_prefix="RAG_", extra="ignore")


@lru_cache
def get_settings() -> Settings:
    return Settings()
