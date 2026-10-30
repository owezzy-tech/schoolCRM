# Lesson assistant protocol v1

Bead `schoolCRM-o58.12` extends the actor-owned lesson generation runs. A lesson assistant thread is one immutable generation request and its progress, retrieved evidence and saved draft. Follow-up generation starts a new request; edits and human decisions use Go lesson APIs.

## Identity and history

`GET /v1/rag/lessons/threads?school_id=UUID&department_id=UUID` returns the caller's latest 50 runs in that department as JSON:API `lesson-thread` resources. Attributes: `request` (the original generate body, flattened), `status` (`pending` or `completed`), `result` (creation receipt or null). `GET /v1/rag/lessons/threads/{request_id}` returns the same resource. Only the actor can read a thread, including administrators and reviewers. Every read and reconnect requires current Go teaching authority. Revoked access fails before checkpoint contents are read. Existing generation/status endpoints remain compatible.

## Streaming and recovery

`POST /v1/rag/lessons/generate/stream` accepts the existing generate body and bearer authentication. It returns `text/event-stream`; validation/authentication failures before streaming use JSON:API errors. Each named SSE event has JSON data containing `protocolVersion: 1` and `requestID`. Names and additional fields:

| Event | Fields | Meaning |
| --- | --- | --- |
| `progress` | `stage`: `retrieve`, `generate` or `persist`; `detail` | Current workflow step |
| `citation` | `citations`: trusted citation array | Retrieved evidence, including source identity, page, revision, checksum and passage |
| `structured-result` | `result`: original Go creation receipt | A saved draft; UI reads its immutable Go version for the complete structured content |
| `approval-request` | `planID`, `version`, `action: "submit-for-review"` | A presentation prompt to open the saved draft for a human decision; never executes approval or submission |
| `completed` | no additional fields | Terminal success, following persistence of the receipt |
| `terminal-error` | `status`, `detail` | Terminal failure with a client-safe message |

Every connected stream finishes with one terminal event. Transport failure/disconnect is an incomplete observation, never proof of completion. Reconnect uses a fresh bearer token and exactly the same request ID and input. It resumes the existing PostgreSQL checkpoint or replays the saved receipt; no separate browser-only conversation database is used. Browser refresh discovers unfinished requests through private server history. Changed input for the same ID returns a conflict. Events are replaceable projections, not an append-only log; clients replace citations/results on replay. No Last-Event-ID cursor is promised. Idle streams carry `: ping` comments, which clients ignore. A `401` before streaming or a `terminal-error` with status `401` means Go rejected the bearer; the client signs in again and returns to the same request rather than retrying. Cancellation releases the request lock; Go receipts prevent duplicate business writes even after an ambiguous response. Private tokens never enter events, checkpoints, prompts or tracing.

Only teaching members see generation controls; HOD/dean members use existing review/approval queues and evidence on exact lesson versions. Reuse is a human command for an explicitly selected version and a stable request ID, preserving provenance while creating a fresh author-owned draft requiring new approvals.

Implementation uses the installed [LangGraph streaming API](https://reference.langchain.com/python/langgraph/pregel/main/Pregel/astream) and [FastAPI Server-Sent Events](https://fastapi.tiangolo.com/tutorial/server-sent-events/), which cancels the workflow when the client disconnects; PostgreSQL checkpoints retain synchronous durability.
