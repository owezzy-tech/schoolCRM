# Hosted curriculum and lesson proof

Bead `schoolCRM-kce.4`. This proof uses the deployed staging HTTPS endpoints, real Cloudflare Workers AI embeddings, real DeepSeek generation and the Go/PostgreSQL lesson store. It does not use `TestClient`, a model fixture or local Ollama.

## Scope and evidence

The staging teacher is a synthetic identity at `schoolai-teacher@staging.invalid`. It has only the global `TEACHER` role and a department-scoped `teach` membership. The deployment administrator reviewed the curriculum source. No student data was supplied.

| Resource | Recorded ID |
| --- | --- |
| School | `ed78d739-d93a-431d-8957-da3e99c78cef` |
| Languages department | `495cf156-ee17-4d06-b271-c4ccbaa47fe4` |
| Synthetic teacher | `701bd8ad-cc6c-4062-a8ca-08e4986ba0df` |
| Approved curriculum source | `b1f3d4e7-1c1d-4811-a18e-b87a5ba6a7c5` |
| Persisted generated plan | `495b0af4-b207-4cf9-b321-382e9ca5cc00` |
| Generation request | `b99268af-c194-47c1-80aa-07b731d2b815` |

The official source is [KICD English Grades 1-3](https://kicd.ac.ke/wp-content/uploads/2025/07/Grade-1-3-English-Activities-Revised-Sept.pdf), restricted to reviewed PDF page 17 for Kenya CBC, Grade 1 English, revision 2024. The original checksum is `eeea001277e1d952f4a6ebd4e42a21072cbb1e261d2e7228e50fe699ef041d0d`.

The hosted sequence passed:

1. The teacher uploaded the checksum-verified PDF with the selected school and department scope.
2. The administrator approved the source. Source readback identified `cloudflare:@cf/baai/bge-m3:1024`.
3. Scoped retrieval returned two approved passages.
4. Real `DeepSeek-V4.1-Flash` generated a 30-minute lesson. The Go version endpoint returned the persisted content with two citations to that source and the same Cloudflare embedding identity. Activity minutes totalled 30.
5. Repeating the same request returned the same plan. The Go store contained one version, and the durable generation receipt reported `completed`.

This proves hosted ingestion, reviewed indexing, retrieval, generation, persistence and request replay. It does not prove separate HOD/dean approval, publication, PDF export or post-restart recovery. Those checks remain in `schoolCRM-kce.5`.

## Repeating the proof

Run `scripts/verify_hosted_ai.py` through the existing RAG environment:

```sh
uv run --project api/services/RAG --frozen python scripts/verify_hosted_ai.py
```

Inject `SCHOOLCRM_HOSTED_PROOF=staging`, `SCHOOLCRM_HOSTED_PDF`, `SCHOOLCRM_STAGING_ADMIN_PASSWORD` and `SCHOOLCRM_STAGING_QA_PASSWORD` through the process environment. For this completed scenario, also set `SCHOOLCRM_HOSTED_SOURCE_ID=b1f3d4e7-1c1d-4811-a18e-b87a5ba6a7c5`. That reuses the approved source rather than uploading another original. The fixed request ID replays the existing plan and does not invoke the model again.

Retrieve credentials privately from the staging `schoolcrm` service. The administrator credential is `SCHOOLCRM_BOOTSTRAP_PASSWORD`; the separate synthetic teacher credential is `SCHOOLCRM_STAGING_QA_PASSWORD`. Never put values in command arguments, source files, reports or screenshots. The script emits resource IDs and verdicts only. It targets the staging domain explicitly and refuses to run without its staging opt-in.

## Credential-safe runtime

Current source already excludes credential header values from successful auth-client request logs. The new regression test found that rejected JWTs still entered the auth failure log. Removing that token field preserves the failure diagnostic and subject while preventing future bearer-token output. The test failed before the repair and passed afterwards under the race detector.

Source commit `72f3ea7` was deployed to staging auth as `2d682d2a-dfae-4ffc-877a-e2783361c36d` and staging business API as `9a750f32-a902-4c21-abe2-6c7c0ef7d37c`. Both reached Railway `SUCCESS` before the hosted proof. The RAG runtime retained its configured Cloudflare/DeepSeek credentials. Production and local databases were untouched.

Earlier token-bearing failure logs, if present, are historical records; this source repair does not erase them. Restrict log access and apply the established retention policy. Do not claim historical log sanitisation or rotate signing keys implicitly.
