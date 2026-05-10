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
uvicorn main:app --reload --host 0.0.0.0 --port 7000
```

### Notes

- LangChain is intentionally kept out of `domain/` and `use_cases/`.
- This scaffold uses simple in-memory adapters so the service boots and the shape is testable.
- Replace the adapters incrementally with real LangChain, vector store, auth, and persistence implementations.
