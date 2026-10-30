from dataclasses import dataclass


@dataclass(frozen=True, slots=True)
class ChunkMetadata:
    source: str
    chunk_index: int
    start_offset: int
    end_offset: int
    collection: str = "default"
