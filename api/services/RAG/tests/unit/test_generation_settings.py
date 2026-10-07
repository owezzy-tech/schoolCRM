from pydantic import SecretStr

from infrastructure.config import Settings


def test_deepseek_key_is_masked_and_can_be_injected_without_environment_names():
    settings = Settings(deepseek_api_key=SecretStr("test-private-key"), _env_file=None)
    assert settings.deepseek_api_key.get_secret_value() == "test-private-key"
    assert "test-private-key" not in repr(settings)


def test_both_service_and_standard_key_names_load_from_dotenv(tmp_path):
    env_file = tmp_path / "provider.env"
    env_file.write_text("DEEPSEEK_API_KEY=standard-key\n", encoding="utf-8")
    settings = Settings(_env_file=env_file)
    assert settings.deepseek_api_key.get_secret_value() == "standard-key"
    env_file.write_text(
        "RAG_DEEPSEEK_API_KEY=service-key\nDEEPSEEK_API_KEY=standard-key\n", encoding="utf-8"
    )
    assert Settings(_env_file=env_file).deepseek_api_key.get_secret_value() == "service-key"
