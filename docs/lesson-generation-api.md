# Lesson generation and reuse

Bead `schoolCRM-o58.3` adds resumable lesson generation from approved curriculum evidence. Go owns business records and publication. Python validates DeepSeek output and calls Go with the teacher's current identity.

## Confirmed contract

The user approved objectives, prerequisites, materials, timed activities, differentiation, assessment and citations on 7 October 2026. Activity times total the requested duration. DeepSeek-V4.1-Flash is selected; its official API identifier is `deepseek-flash` at `https://api.deepseek.com`, verified against the [provider release notes](https://api-docs.deepseek.com/updates) and [model catalogue](https://api-docs.deepseek.com/api/list-models).

Only a current department `teach` capability permits generation. Global roles are insufficient. Retrieval checks school, department, framework, stage, subject and revision. Missing or withdrawn evidence cannot support a resumed request. Models cannot approve, publish or grant permissions.

## Configuration

Set `RAG_DEEPSEEK_API_KEY` in the RAG service's ignored `.env`. `DEEPSEEK_API_KEY` is also accepted; the service-prefixed name takes precedence. Keys are masked and never enter prompts or checkpoints. Restart the service after changing its environment. `.env.example` contains no key.

Configure `RAG_CURRICULUM_DATABASE_URL`, the existing embedding provider and Go service URLs. Apply Go migration `1.32` and the explicit Python migrations:

```sh
cd api/services/RAG
uv sync --extra dev --locked
uv run python -m infrastructure.curriculum_migrate
uv run python -m infrastructure.lesson_migrate
```

The Python migration creates actor-owned generation records and the installed LangGraph PostgreSQL checkpoint tables. Go migration `1.32` adds creation receipts. Application rollback retains these additive tables and immutable snapshots; do not discard their history after use.

## Generate a draft

`POST /v1/rag/lessons/generate` requires a bearer token and this body:

```json
{
  "request_id": "11111111-1111-4111-8111-111111111111",
  "school_id": "22222222-2222-4222-8222-222222222222",
  "department_id": "33333333-3333-4333-8333-333333333333",
  "framework": "kenya-cbc",
  "stage": "grade-1",
  "subject": "english",
  "revision": "2024",
  "topic": "School vocabulary and pronunciation",
  "duration_minutes": 30
}
```

Use a new nonzero UUID per intended generation. Duration is a whole number from 1 to 180 minutes. The service retrieves evidence, requests structured JSON, validates it and saves a draft through Go. Missing evidence returns 422. Invalid, unsupported or truncated model output is not saved. URLs, checksums, revision identities and page numbers come from retrieved records, not model-generated URLs.

Success is a JSON:API `lesson-generation` resource with `planID`, `version`, `title`, `creationStatus`, citations and generation provenance. This is the original creation receipt. Read Go plan/version APIs for current workflow status after edits or publication.

`GET /v1/rag/lessons/generations/{request_id}` returns the caller's request as pending or completed. Another user receives 404. Reads and retries recheck teaching authority. The same actor/request resumes its checkpoint or returns its receipt. Changed input with the same ID returns 409.

## Recovery and validation

LangGraph checkpoints retrieval, generation and persistence in PostgreSQL. Request records bind actor and input. Bearer tokens live only in transient runtime context; every resume supplies fresh credentials. External tracing is disabled for this private workflow. Permission to process student data does not grant access to other schools or students.

Python uses a session advisory lock per actor/request. Go uses a distinct transaction-lock namespace, preventing self-deadlock when both services share a database. Go checks authority before replaying receipts. The plan, version, pointer, audit and receipt commit atomically. An ambiguous reply can retry persistence without another plan or another checkpointed model call.

Model output may select only retrieved citation IDs. Unknown IDs, extra command fields, wrong types, empty required fields and incorrect activity totals fail validation. Go independently validates generated `schemaVersion: 1` content. Legacy manually-authored JSON remains compatible; old snapshots are not rewritten. Editors preserve citation/provenance fields and recalculate duration after activity edits.

## Reuse a version

`POST /v1/lessons/{plan_id}/reuse` accepts `version`, nonzero `requestID` and optional `title`. It requires teaching authority in the source department. Another author's unsubmitted draft remains private.

Reuse creates a new author-owned draft version 1. It preserves content/citations and records `reusedFrom.planID` and `reusedFrom.version`. It copies no review, approval or publication state. New publication requires the full teacher/HOD/dean workflow. Same-command retries return the original receipt; changed input returns 409.

## Proof and limits

Tests cover strict output, exact citations, model identity, secret masking, concurrent retries, actor isolation, fresh-token recovery, withdrawal, revocation and immutable history. The optional HTTP journey uses real Go authentication, PostgreSQL, the official KICD English original and local BGE-M3. Its default DeepSeek transport is a contract fixture. Set `DEEPSEEK_LIVE_TEST=1` with a configured key for real-provider verification. Do not describe fixture results as live DeepSeek proof.

Teacher assistant/evidence/reuse controls remain under `schoolCRM-o58.12`. Broader corpus acquisition and integrated production validation remain in their respective beads.
