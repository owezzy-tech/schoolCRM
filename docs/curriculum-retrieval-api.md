# Curriculum retrieval API

Bead `schoolCRM-o58.2` adds PDF ingestion, source review and semantic retrieval to the existing Python RAG service. The [School AI ARD](school-ai-ard.md) defines the wider lesson workflow. Generation and assistant screens remain separate features.

## Authority and review

Every request needs a bearer token validated by the Go authentication service. The Python service asks Go for current school and department permissions. A current teacher, HOD or dean capability permits uploading and reading department sources. Current school management permits approval and withdrawal. Global role names in a token do not grant school access. Anonymous mode does not grant curriculum access.

Uploads stay pending. The reviewer must check the original's publication authority, revision, framework, stage, subject, topic and selected page range. Approval records the reviewer and their note. The API records provenance supplied by staff; it cannot establish copyright permission or factual accuracy from a filename. The original Grades 1–3 English document has multiple grades, so select and review the pages for each intended grade separately.

Only approved sources enter search. Withdrawal removes a source from new queries and preserves its original, chunks and citation identifiers. There is no destructive source endpoint. Replacing a publication means uploading and reviewing a new revision.

## Configuration

Install locked dependencies with `uv sync --extra dev` in `api/services/RAG`. The database must support pgvector. Use a separate Python-owned database or a database role limited to the `rag_curriculum_*` tables. Python does not write Go business records.

| Setting | Purpose |
| --- | --- |
| `RAG_CURRICULUM_DATABASE_URL` | PostgreSQL connection URL. Without it, curriculum endpoints return 503 |
| `RAG_AUTH_SERVICE_URL` | Existing Go authentication endpoint |
| `RAG_SCHOOL_SERVICE_URL` | Go SchoolCRM endpoint, default `http://schoolcrm:3000` |
| `RAG_CURRICULUM_EMBEDDING_PROVIDER` | `ollama` for development or `cloudflare` for production |
| `RAG_OLLAMA_BASE_URL` | Local Ollama endpoint with `bge-m3:latest` installed |
| `RAG_CLOUDFLARE_ACCOUNT_ID` | Workers AI account, required for Cloudflare embeddings |
| `RAG_CLOUDFLARE_API_TOKEN` | Workers AI credential, supplied through protected environment configuration |

Run the explicit additive migration before enabling the endpoints:

```sh
uv run python -m infrastructure.curriculum_migrate
```

The migration creates the pgvector extension and Python-owned source, chunk and audit tables. Use a migration role with extension privileges, then a runtime role with table permissions. Schema creation is transactional and repeatable. Deployment rollback can restore the preceding application while retaining these tables. Preserve them and their backups after use.

The original PDF bytes, checksum, immutable parsing metadata and chunks are stored in PostgreSQL. Database backup and restore therefore include the originals. Approval writes vectors, review state, model identity and audit history in one transaction. A failed audit rolls back the index and approval. Concurrent reviews require the expected state and allow only one successful transition.

## Endpoints

The current OpenAPI schema is available at `/openapi.json`. Success and error documents use JSON:API v1.1 and `application/vnd.api+json`. The original download returns a PDF with private, non-cacheable attachment headers.

| Method and path | Behaviour |
| --- | --- |
| `POST /v1/rag/curriculum/sources` | Multipart upload with `file` and JSON `metadata`. Returns a pending source |
| `GET /v1/rag/curriculum/sources/{id}` | Authorised source metadata and current review/index identity |
| `GET /v1/rag/curriculum/sources/{id}/original` | Authorised download of the exact original |
| `POST /v1/rag/curriculum/sources/{id}/review` | School management approves a pending source or withdraws an approved source |
| `POST /v1/rag/curriculum/search` | Exact-scope semantic search with page citations, scores and index identity |

Upload metadata example:

```json
{
  "school_id": "11111111-1111-4111-8111-111111111111",
  "department_id": "22222222-2222-4222-8222-222222222222",
  "framework": "kenya-cbc",
  "stage": "grade-1",
  "subject": "english",
  "revision": "2024",
  "title": "English Language Activities, Grades 1–3",
  "topic": "School vocabulary",
  "authority": "Kenya Institute of Curriculum Development",
  "source_url": "https://kicd.ac.ke/wp-content/uploads/2025/07/Grade-1-3-English-Activities-Revised-Sept.pdf",
  "sha256": "eeea001277e1d952f4a6ebd4e42a21072cbb1e261d2e7228e50fe699ef041d0d",
  "first_page": 17,
  "last_page": 17
}
```

PDF viewer page numbers are one-based and can differ from printed pages. Uploads have a 10 MiB limit and must match their acquisition checksum. Encrypted, invalid and pages without sufficient extracted text return 422. Scanned pages need separate OCR and human review; the parser does not silently invent or approve extracted text.

Approve with `{"expected_status":"pending","decision":"approved","note":"Publication, revision and selected pages checked"}`. Withdraw with `{"expected_status":"approved","decision":"withdrawn","note":"Superseded by revised publication"}`. A stale or invalid transition returns 409. Duplicate upload of the same original and page range in the same scope returns 409.

Search uses the six scope fields from the example plus `question` and an optional `limit` between 1 and 20. Scope filters run in PostgreSQL before passages are returned. No framework equivalence is inferred. Cambridge accepts Early Years and Years 1–9; acquired IGCSE files do not expand that scope.

Evidence includes the source URL, authority, revision, checksum, original page, chunk ordinal, similarity score, embedding identity and index revision. The index revision is the immutable source UUID. The model identity and dimensions are recorded separately. Empty results return `meta.abstained = true`; this endpoint retrieves evidence and does not generate an answer. Generation must retain citations and recheck active sources before model exposure.

## Embeddings and validation

Both providers must return 1024-dimensional, finite, nonzero vectors. Ollama model digests are checked before and after embedding and stored with the index. A changed digest fails indexing; a new model identity cannot search an old index. Cloudflare uses the managed `@cf/baai/bge-m3` identity in a separate provider namespace. Local and Workers vectors are never mixed. Cloudflare does not expose the same local weight digest; validate its managed model and rebuild a separate production corpus before rollout.

The initial similarity floor is 0.55. It passed the recorded Grade 1 English query and an unrelated question; this is limited validation, not proof for every subject or question. Evaluate retrieval fixtures for each added corpus. Do not treat a similarity score as proof of a factual claim.

Run unit and database tests:

```sh
uv run ruff check .
uv run pytest tests -q
CURRICULUM_TEST_DATABASE_URL=postgresql://... uv run pytest tests/integration/test_curriculum.py -q
```

Set `CURRICULUM_LIVE_PDF` to the authorised local Grades 1–3 English PDF to include the live Ollama test. The committed optional HTTP journey also needs `CURRICULUM_LIVE_AUTH_URL`, `CURRICULUM_LIVE_SCHOOL_URL` and `CURRICULUM_LIVE_RAG_URL`. Run it only against isolated, seeded test services. It creates a school and membership, retrieves the reviewed source, revokes the teacher with their token still active, withdraws the source and verifies original preservation.

Cloudflare uses its documented [OpenAI-compatible embeddings endpoint](https://developers.cloudflare.com/workers-ai/configuration/open-ai-compatibility/), verified on 7 October 2026, with `@cf/baai/bge-m3` and strict response-index validation. Request/response handling has contract tests. Live Workers credentials, production backup restoration and broader authorised source acquisition remain separate validation work. Preview summaries and unverified publisher candidates are not automatically approved curriculum evidence.
