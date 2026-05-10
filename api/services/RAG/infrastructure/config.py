from functools import lru_cache

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    service_name: str = "rag"
    api_host: str = "0.0.0.0"
    api_port: int = 7000
    allow_anonymous: bool = True
    file_storage_dir: str = "./var/files"

    model_config = SettingsConfigDict(env_prefix="RAG_", extra="ignore")


@lru_cache
def get_settings() -> Settings:
    return Settings()
