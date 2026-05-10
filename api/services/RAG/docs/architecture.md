# RAG Service Architecture

## Overview

The `RAG` service is a Python FastAPI service inside the `schoolCRM` monorepo. Its job is to ingest documents and answer questions using retrieval-augmented generation.

This service uses **Hexagonal Architecture** in a single bounded context: **Document Intelligence**.

Why this shape fits the repo:

- the service remains deployable and independent under `api/services/`
- the business core stays separate from FastAPI, LangChain, and storage details
- Go services can call the Python service through a stable HTTP contract
- adapters can be swapped later without changing the use cases

## Bounded Context

This is **one bounded context**, not two.

The service has two workflows inside that context:

1. **Ingestion**
   - upload raw file
   - persist original file
   - parse text
   - chunk text
   - embed chunks
   - store chunks and metadata

2. **Query**
   - accept question
   - embed question
   - retrieve top matching chunks
   - generate answer from retrieved context

These are separate workflows, but they operate on the same document model, vector store, and embedding space, so they should stay in one context.

## Layer Model

```text
┌──────────────────────────────────────────────────────────────┐
│ Driving adapters                                             │
│ FastAPI controllers, request parsing, auth dependency        │
├──────────────────────────────────────────────────────────────┤
│ Application layer                                            │
│ IngestDocumentUseCase, QueryDocumentsUseCase,                │
│ DeleteDocumentUseCase                                        │
├──────────────────────────────────────────────────────────────┤
│ Domain core                                                  │
│ Entities, value objects, ports, simple types                 │
├──────────────────────────────────────────────────────────────┤
│ Driven adapters                                              │
│ Parser, embedding provider, vector store, LLM, file store,   │
│ metadata repository                                          │
├──────────────────────────────────────────────────────────────┤
│ Infrastructure                                               │
│ App assembly, config, lifespan, observability, auth wiring   │
└──────────────────────────────────────────────────────────────┘
```

## Dependency Rule

Imports must always point inward:

- `domain/` imports nothing from adapters or infrastructure
- `use_cases/` imports only from `domain/`
- `adapters/` implement `domain/ports/`
- `infrastructure/` wires everything together

That rule is the main reason this service will stay easy to test and easy to evolve.

## Directory Responsibilities

```text
api/services/RAG/
├── main.py                    # Entrypoint, exposes FastAPI app
├── openapi.yaml               # Committed API contract
├── pyproject.toml             # Python package + dependencies
├── README.md                  # Quick start and summary
├── docs/
│   └── architecture.md        # This document
├── domain/
│   ├── entities/              # Core business models
│   ├── value_objects/         # Small immutable domain concepts
│   ├── ports/                 # Abstract interfaces to the outside world
│   └── types.py               # Lightweight shared types
├── use_cases/                 # Application orchestration
├── adapters/
│   ├── controllers/           # HTTP entrypoints
│   ├── parsers/               # File parsing implementations
│   ├── embeddings/            # Embedding adapter implementations
│   ├── vector_stores/         # Vector store implementations
│   ├── llm/                   # LLM adapter implementations
│   ├── repositories/          # Metadata persistence implementations
│   └── storage/               # Raw file persistence implementations
├── infrastructure/
│   ├── app.py                 # FastAPI assembly
│   ├── auth.py                # Auth dependency placeholder
│   ├── config.py              # Environment-driven settings
│   ├── dependencies.py        # DI / container wiring
│   ├── lifespan.py            # Startup and shutdown lifecycle
│   └── observability.py       # Logging / telemetry setup
└── tests/
    └── unit/                  # Fast unit tests
```

## Core Domain Concepts

### Entities

- `DocumentRecord`
  - metadata for an ingested document
  - document id, title, source, status, uploader, file location, chunk count

- `DocumentChunk`
  - one stored retrieval chunk
  - includes source metadata and owning document id

- `Query`
  - question plus caller identity and retrieval settings

- `QueryResult`
  - generated answer plus provenance-like document ids and snippets

### Value Objects

- `ChunkMetadata`
  - source, chunk index, start offset, end offset

- `DocumentType`
  - content type classification for upload handling

## Ports

The domain defines these external contracts:

- `IDocumentParser`
- `IEmbeddingProvider`
- `IVectorStore`
- `ILLMProvider`
- `IDocumentRepository`
- `IFileStore`

These are the seams that keep LangChain and infrastructure details out of the core.

## Use Cases

### `IngestDocumentUseCase`

Responsibilities:

- persist raw upload through `IFileStore`
- parse text through `IDocumentParser`
- split text internally using service config
- embed each chunk through `IEmbeddingProvider`
- persist vectors through `IVectorStore`
- persist document metadata through `IDocumentRepository`

### `QueryDocumentsUseCase`

Responsibilities:

- embed the question
- retrieve relevant chunks
- generate an answer from retrieved context
- return answer plus snippet evidence

### `DeleteDocumentUseCase`

Responsibilities:

- remove metadata
- remove vectors
- remove raw file

## Why some abstractions were intentionally not added

To keep this pragmatic for a school CRM:

- there is **no reranker port yet**
- there is **no text splitter port yet**
- embeddings are a lightweight type alias, not a heavy value object

These can be added later if actual requirements justify them.

## Contract

The canonical contract is in `openapi.yaml`.

### Endpoints

```text
GET    /v1/liveness
GET    /v1/readiness
POST   /v1/rag/documents
DELETE /v1/rag/documents/{document_id}
POST   /v1/rag/query
```

### Contract principles

- stable JSON responses
- bearer auth support
- health probes aligned with existing repo conventions
- explicit error shape with `code` and `message`

## Integration with the Go Monorepo

The repo already uses interface-based service clients like `authclient`.

The RAG service should follow the same pattern on the Go side:

```text
app/sdk/ragclient/
├── ragclient.go
└── http/
    └── ragclient.go
```

Recommended Go-side interface:

```text
type Client interface {
    Query(ctx, req) (resp, error)
    Ingest(ctx, req) (resp, error)
    Delete(ctx, documentID) error
    Close() error
}
```

That preserves the same monorepo contract style already used by `authclient`.

## Deployment Shape

The scaffold includes:

- `zarf/docker/dockerfile.rag`
- `rag` entry in `zarf/compose/docker_compose.yaml`

The container currently serves FastAPI on port `7000`.

Health checks:

- `/v1/liveness`
- `/v1/readiness`

## Current Scaffold Status

What is real right now:

- route structure
- use case boundaries
- OpenAPI contract
- local file persistence
- in-memory metadata + vector store
- placeholder parser, embedding provider, and LLM

What should be replaced next:

- auth placeholder with JWKS validation against auth service
- in-memory repository with persistent metadata storage
- in-memory vector store with real vector database
- noop parser with LangChain-backed loaders/parsers
- echo LLM with real model provider

## Recommended next implementation steps

1. Create Go `ragclient` package
2. Replace auth placeholder with real JWT validation
3. Add persistent metadata repository
4. Add LangChain-backed parser + embedding adapter
5. Add real vector store adapter
6. Add integration tests around ingestion and query flows
