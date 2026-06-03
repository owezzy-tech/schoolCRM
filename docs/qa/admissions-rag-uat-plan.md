# Admissions CRM + RAG UAT Plan

## Scope

This plan covers the Admissions CRM and Admissions RAG flows tracked by Beads epic `schoolCRM-zg8`.

Environment:

- Kubernetes: kind cluster `schoolcrm-cluster`
- Namespace: `schoolcrm-system`
- Auth service: `auth-service:6000`
- RAG service: `rag-service:7000`
- RAG image: `localhost/schoolcrm/rag:0.0.1`
- Auth image: `localhost/schoolcrm/auth:0.0.1`
- Ollama model: `nemotron-3-super:cloud`
- RAG answer provider: `ollama`
- RAG graph retriever: `memory` for current dev smoke

## UAT Matrix

### Backend Smoke

| Area | Scenario | Expected result | Status |
| --- | --- | --- | --- |
| RAG health | `GET /v1/liveness` | `200 OK` | Pass |
| RAG readiness | `GET /v1/readiness` | `200 OK` | Pass |
| RAG to Ollama | RAG pod calls `host.docker.internal:11434/api/generate` with `nemotron-3-super:cloud` | Model responds | Pass |
| Auth login | `admin@example.com` / `gophers` | JSON:API login response with bearer token | Pass |
| Auth authenticate | Admin bearer token | JSON:API authenticate response includes system and admissions roles | Pass |
| Staff RAG | Staff token calls `/v1/rag/admissions/staff/query` | JSON:API `rag-result` | Pass |
| Applicant RAG | Applicant token calls `/v1/rag/admissions/applicant/applications/{id}/query` | Own-application scoped JSON:API result | Blocked by `schoolCRM-zg8.3` |
| Reviewer RAG | Reviewer token calls `/v1/rag/admissions/reviewer/applications/{id}/query` | Assigned-application scoped JSON:API result | Blocked by `schoolCRM-zg8.3` |

### Frontend UAT

| Area | Route/screen | Checks | Status |
| --- | --- | --- | --- |
| Staff applications | `/applications` | List, filter, status badges, navigation to detail | Pending |
| Application detail | `/applications/:id` | Applicant summary, documents, timeline, decisions | Pending |
| Reviewer workspace | Application review screen | Assigned reviewer workflow and RAG assistant affordance | Pending |
| Leads | `/leads` | Static data renders and primary actions are reachable | Pending |
| Admissions settings | `/admissions-settings` | Configuration screen renders static settings | Pending |
| Applicant portal | `/applicant/...` | Application/status screens render and remain applicant scoped | Pending |

## Execution Order

1. Bring up canonical kind environment.
2. Apply database, Auth, and RAG overlays.
3. Run database migration and seed data.
4. Patch RAG dev config to use the Ollama answer provider.
5. Run backend smoke in this order:
   - RAG health/readiness
   - RAG pod to Ollama
   - Auth login
   - Auth authenticate
   - Staff RAG protected query
   - Applicant and reviewer RAG queries when fixtures exist
6. Run frontend build and admissions route smoke.
7. Add Playwright coverage for admissions happy paths once the frontend E2E harness is available.
8. Record defects as Beads children under `schoolCRM-zg8`.

## Acceptance Criteria

- Smoke results are recorded in `docs/qa/admissions-rag-uat-report.md`.
- One Beads issue exists for each failed or blocked acceptance item.
- Staff RAG protected smoke passes through Auth, RAG, and Ollama.
- Applicant and reviewer RAG scope smoke pass once `schoolCRM-zg8.3` is resolved.
- Frontend admissions Playwright coverage is added or a blocker bead documents why it cannot run.
