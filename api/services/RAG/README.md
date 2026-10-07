## RAG Service

Python FastAPI service for document ingestion and retrieval-augmented generation in the `schoolCRM` monorepo.

### Architecture

This service uses **Hexagonal Architecture (Ports and Adapters)** inside a single bounded context: **Document Intelligence**.

Full architecture documentation lives in `docs/architecture.md`.

### Final Structure

```text
api/services/RAG/
├── main.py
├── openapi.yaml
├── pyproject.toml
├── README.md
├── docs/
│   └── architecture.md
├── domain/
│   ├── entities/
│   ├── value_objects/
│   ├── ports/
│   └── types.py
├── use_cases/
├── adapters/
│   ├── controllers/
│   ├── parsers/
│   ├── embeddings/
│   ├── vector_stores/
│   ├── llm/
│   ├── repositories/
│   └── storage/
├── infrastructure/
└── tests/
```

### Contract

The committed contract lives in `openapi.yaml` and exposes:

- `GET /v1/liveness`
- `GET /v1/readiness`
- `POST /v1/rag/documents`
- `DELETE /v1/rag/documents/{document_id}`
- `POST /v1/rag/query`

### Local Run

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -e .
uvicorn main:app --reload --host 0.0.0.0 --port 4545
```

### Notes

- The [curriculum retrieval API](../../../docs/curriculum-retrieval-api.md) uses LlamaIndex PDF chunking, reviewed PostgreSQL/pgvector sources and BGE-M3 embeddings. Configure its database and run the explicit migration before use. `/openapi.json` includes its current request schemas.

- LangChain is intentionally kept out of `domain/` and `use_cases/`.
- Generic document endpoints still use in-memory and placeholder adapters. They do not establish curriculum review or persistence.
