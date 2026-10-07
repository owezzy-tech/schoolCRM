import math

import httpx

from domain.entities.curriculum import CurriculumError

DIMENSIONS = 1024


class BGEM3Embeddings:
    """Provider identity is part of every stored index and every search filter."""

    def __init__(
        self,
        provider: str,
        ollama_url: str,
        cloudflare_account: str | None = None,
        cloudflare_token: str | None = None,
        transport: httpx.AsyncBaseTransport | None = None,
    ) -> None:
        self.provider = provider
        self.ollama_url = ollama_url.rstrip("/")
        self.account = cloudflare_account
        self.token = cloudflare_token
        self.client = httpx.AsyncClient(timeout=120, transport=transport)

    async def embed(self, texts: list[str]) -> tuple[str, list[list[float]]]:
        try:
            if self.provider == "ollama":
                tags = await self.client.get(f"{self.ollama_url}/api/tags")
                tags.raise_for_status()
                model = next(
                    (m for m in tags.json()["models"] if m["name"] == "bge-m3:latest"), None
                )
                if not model or not model.get("digest"):
                    raise CurriculumError(503, "Install local Ollama bge-m3 before indexing")
                identity = f"ollama:bge-m3:{model['digest']}"
                response = await self.client.post(
                    f"{self.ollama_url}/api/embed",
                    json={"model": model["name"], "input": texts, "truncate": False},
                )
                response.raise_for_status()
                vectors = response.json()["embeddings"]
                after = await self.client.get(f"{self.ollama_url}/api/tags")
                after.raise_for_status()
                if not any(
                    m.get("name") == model["name"] and m.get("digest") == model["digest"]
                    for m in after.json()["models"]
                ):
                    raise CurriculumError(503, "Embedding model changed during request")
            elif self.provider == "cloudflare":
                if not self.account or not self.token:
                    raise CurriculumError(503, "Cloudflare embedding credentials are required")
                identity = "cloudflare:@cf/baai/bge-m3:1024"
                response = await self.client.post(
                    f"https://api.cloudflare.com/client/v4/accounts/{self.account}"
                    "/ai/v1/embeddings",
                    headers={"Authorization": f"Bearer {self.token}"},
                    json={"model": "@cf/baai/bge-m3", "input": texts},
                )
                response.raise_for_status()
                payload = response.json()
                rows = sorted(payload["data"], key=lambda row: row["index"])
                if [row["index"] for row in rows] != list(range(len(texts))):
                    raise CurriculumError(503, "Cloudflare embedding response has invalid indexes")
                vectors = [row["embedding"] for row in rows]
            else:
                raise CurriculumError(503, "Unsupported curriculum embedding provider")
            if len(vectors) != len(texts) or any(
                len(v) != DIMENSIONS
                or any(not isinstance(x, (int, float)) or not math.isfinite(x) for x in v)
                or not any(v)
                for v in vectors
            ):
                raise CurriculumError(503, "Embedding response has invalid dimensions or values")
            return identity, vectors
        except (httpx.HTTPError, KeyError, TypeError, ValueError) as exc:
            raise CurriculumError(503, "Curriculum embedding provider unavailable") from exc

    async def close(self) -> None:
        await self.client.aclose()
