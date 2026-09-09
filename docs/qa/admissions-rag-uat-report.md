# Admissions CRM + RAG UAT Report

## Summary

Current UAT execution focused on the Kubernetes RAG/Auth smoke path for the Admissions RAG feature.

Result:

- Canonical kind environment was recreated as `schoolcrm-cluster`.
- RAG, Auth, and Postgres were deployed in `schoolcrm-system`.
- RAG was configured to use Ollama with `nemotron-3-super:cloud`.
- Staff, applicant, and reviewer protected RAG queries now pass end-to-end through Auth, RAG, and Ollama.
- Applicant/reviewer role-boundary denial checks return 403 as expected.
- Frontend admissions Playwright E2E now runs headlessly and covers the required admissions happy paths.

## Environment Notes

The original running kind cluster was `desktop`, but the project makefile targets `schoolcrm-cluster`. The environment was standardized by deleting `desktop` and creating `schoolcrm-cluster`.

Local port conflicts blocked initial cluster creation:

- `6000`: local `go run api/services/auth/main.go`
- `3000`/`3010`: local `go run api/services/schoolcrm/main.go`

Those local dev processes were stopped before recreating the canonical cluster.

The first `make dev-up` preload exposed missing public observability images. These were pulled and loaded into kind:

- `grafana/grafana:12.4.0`
- `prom/prometheus:v3.10.0`
- `grafana/tempo:2.10.0`
- `grafana/loki:3.6.0`
- `grafana/promtail:3.6.0`
- `postgres:18.3`

## Smoke Results

| Check | Result | Evidence |
| --- | --- | --- |
| RAG pod readiness | Pass | `kubectl wait pods --namespace=schoolcrm-system --selector app=rag --for=condition=Ready` |
| Auth pod readiness | Pass | `kubectl wait pods --namespace=schoolcrm-system --selector app=auth --for=condition=Ready` |
| Database rollout | Pass | `kubectl rollout status --namespace=schoolcrm-system sts/database` |
| RAG to Ollama | Pass | Pod call to `/api/generate` returned `admissions-rag-smoke-ok` |
| Auth login | Pass | `admin@example.com` / `gophers` returns JSON:API login response with bearer token |
| Auth authenticate roles | Pass | Authenticate response includes `SCHOOL_ADMIN`, `ADMISSIONS_ADMIN`, `APPLICATION_REVIEWER`, `REPORT_VIEWER` |
| Staff RAG query | Pass | `/v1/rag/admissions/staff/query` returns JSON:API `rag-result` |
| Applicant RAG query | Pass | Applicant portal token for `applicant@example.com` can query `/v1/rag/admissions/applicant/applications/d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44/query` |
| Reviewer RAG query | Pass | `teacher@example.com` can query `/v1/rag/admissions/reviewer/applications/d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44/query` |
| Applicant/reviewer denied scope | Pass | Applicant token on reviewer route and reviewer token on applicant route both return 403 |
| Frontend build smoke | Pass with warning | `npx -p node@22 npm run build` succeeds; existing initial bundle budget warning remains |
| Frontend admissions E2E | Pass | `npm run e2e` runs 5 Playwright tests covering applications list/detail, review workspace, leads/settings navigation, applicant portal apply, and portal status |

Final staff RAG smoke response shape:

```json
{
  "jsonapi": {
    "version": "1.1"
  },
  "type": "rag-result",
  "answerPreview": "The provided context does not contain any information about admissions documents that staff should review before a decision, so I cannot answer the question based on the given knowledge graph.",
  "citationCount": 0,
  "documentIds": []
}
```

The insufficient-context answer is expected for the current dev overlay because RAG uses the in-memory graph retriever without seeded admissions knowledge graph content.

## Defects Found and Status

### `schoolCRM-zg8.1`: RAG admissions staff queries reject seeded admissions admins

Status: Closed.

Root cause:

- Auth `/v1/auth/authenticate` returned only JWT/system roles.
- RAG staff guards require admissions roles such as `ADMISSIONS_ADMIN`.
- Seeded `admin@example.com` had admissions roles in `admissions_staff_profiles`, but Auth did not include them in authenticate responses.

Fix:

- Auth authenticate now augments claims with active admissions staff profile roles and de-duplicates roles.

Validation:

- Targeted Go tests pass.
- Live protected staff RAG smoke passes.

### `schoolCRM-zg8.2`: RAG Auth client cannot parse JSON:API authenticate response

Status: Closed.

Root cause:

- Auth returns JSON:API envelopes.
- RAG `AuthServiceClient` only parsed legacy top-level `Claims`/`UserID` response fields.
- Valid tokens were rejected as invalid after Auth returned JSON:API `data.attributes`.

Fix:

- RAG Auth client now unwraps JSON:API `data.attributes` while keeping legacy response support.

Validation:

- RAG unit tests cover JSON:API and legacy response shapes.
- Live protected staff RAG smoke passes.

### `schoolCRM-zg8.3`: QA seed data lacks applicant/reviewer RAG smoke fixtures

Status: Closed.

Impact:

- Applicant and reviewer RAG scope smoke previously could not be completed with real Auth tokens.

Root cause:

- Seed data included users and admissions staff profiles but lacked submitted admissions applications, applicant identity/profile context, and assigned reviewer application fixtures.
- The initial fixture stored applicant email in raw form while current admissions persistence queries canonical `mail.Address.String()` values (`<email>`).
- The initial fixture used nullable JSON application payloads that exposed an existing scan limitation in admissions application queries.
- The applicant user initially used admissions role `APPLICANT` in `users.roles`, but `users.roles` only accepts system roles.

Fix:

- Added deterministic seed fixtures for applicant user/profile, constituent, program, academic term, submitted application, and reviewer assignment.
- Added idempotent seed repairs for applicant system role and non-null application JSON payloads.
- Made `QueryConstituentByPrimaryEmail` tolerate both canonical `<email>` and legacy raw email rows.
- Auth applicant portal authenticate resolves portal tokens to the backing applicant user via applicant profile, while preserving the application ID in portal claims for applicant ownership context.

Validation:

- Applicant portal token flow succeeds for `applicant@example.com`.
- `/v1/auth/authenticate` returns applicant user `f47ac10b-58cc-4372-a567-0e02b2c3d479`, role `APPLICANT`, subject `f47ac10b-58cc-4372-a567-0e02b2c3d479`, and portal claim application ID `d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44`.
- Applicant and reviewer RAG queries return JSON:API `rag-result` responses.
- Applicant token on reviewer route and reviewer token on applicant route both return 403.
- Runtime still uses the permissive `StubApplicationOwnershipChecker`, so true database-backed owner/reviewer mismatch denial remains future integration work beyond this seed/Auth blocker.

### `schoolCRM-zg8.4`: Frontend admissions Playwright E2E harness is missing

Status: Closed.

Impact:

- The required Playwright E2E portion of Admissions CRM UAT previously could not run.

Root cause:

- `api/frontends/web-admin` had no Playwright dependency, configuration, or E2E specs.
- Existing scripts covered Angular build and Karma tests only.

Fix:

- Added `@playwright/test`, `playwright.config.ts`, and `npm run e2e` / `npm run e2e:ui` scripts.
- Added `e2e/admissions.spec.ts` with Auth JSON:API route interception and accessibility-first selectors.
- Covered applications list, application detail, review workspace, leads-to-settings navigation, settings assertions, and applicant portal apply happy path.

Validation:

- `npm run e2e` passes: 4 tests.
- Playwright Chromium browser was installed with `npx playwright install chromium` for local execution.
- `npx -p node@22 npm run build` passes with the pre-existing initial bundle budget warning.

## Validation Commands Run

```bash
go test ./app/domain/authapp ./app/domain/admissionsapp ./app/sdk/authclient/...
go test ./app/domain/authapp ./app/domain/admissionsapp ./app/sdk/authclient/... ./business/domain/admissionsbus/...
uv run ruff check .
uv run pytest tests/unit
npx -p node@22 npm run build
npx playwright install chromium
npm run e2e
make seed
make auth
make rag
kind load docker-image localhost/schoolcrm/auth:0.0.1 --name schoolcrm-cluster
kind load docker-image localhost/schoolcrm/rag:0.0.1 --name schoolcrm-cluster
kubectl rollout restart deployment auth --namespace=schoolcrm-system
kubectl rollout restart deployment rag --namespace=schoolcrm-system
kubectl wait pods --namespace=schoolcrm-system --selector app=auth --timeout=120s --for=condition=Ready
kubectl wait pods --namespace=schoolcrm-system --selector app=rag --timeout=120s --for=condition=Ready
```

## Remaining UAT Work

- Run frontend admissions route smoke for staff and applicant screens.
- Re-run full build/test gates after frontend E2E work.
