from dataclasses import dataclass, field
from datetime import datetime, timezone

from domain.value_objects.chunk import ChunkMetadata


@dataclass(slots=True)
class DocumentRecord:
    id: str
    title: str
    source: str
    content_type: str
    uploaded_by: str
    status: str
    file_path: str | None = None
    chunk_count: int = 0
    created_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))


@dataclass(slots=True)
class DocumentChunk:
    id: str
    document_id: str
    text: str
    metadata: ChunkMetadata
