import asyncio
from pathlib import Path

from domain.ports.file_store import IFileStore


class LocalFileStore(IFileStore):
    def __init__(self, base_dir: str) -> None:
        self._base_dir = Path(base_dir)
        self._base_dir.mkdir(parents=True, exist_ok=True)

    async def save(self, *, document_id: str, filename: str, payload: bytes) -> str:
        target = self._base_dir / f"{document_id}-{filename}"
        await asyncio.to_thread(target.write_bytes, payload)
        return str(target)

    async def delete(self, path: str | None) -> None:
        if not path:
            return

        target = Path(path)
        if target.exists():
            await asyncio.to_thread(target.unlink)
