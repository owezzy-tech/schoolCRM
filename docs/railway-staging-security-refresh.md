# Railway staging Go security refresh

Bead `schoolCRM-8fc`. On 9 October 2026, auth and the business API were deployed from clean `develop` commit `a19c646b03fdd87b432de7317dd482343a903c01`. This activates the Go and dependency patches in PR #79 and the adapter cancellation repair in PR #83. The other intervening merged changes are included in this source snapshot.

## Exact release scope

Every Railway operation used project `174f33cc-0bb5-4d8b-88d8-1288f60792a6` and staging environment `10e68a24-8d12-4210-8b23-ea532fb477a2` explicitly. The project readback listed staging as its only environment. No production promotion was performed.

| Service | Previous successful deployment | New deployment observed `SUCCESS` |
| --- | --- | --- |
| auth (`2a0f7cac-8d37-4b93-89a4-f15d259f52d6`) | `2d682d2a-dfae-4ffc-877a-e2783361c36d` | `f226c1ba-fba3-4b4c-89f7-da4852594a92` |
| schoolcrm (`e332dc5d-ea3c-41a5-ae01-b415fccd80aa`) | `9a750f32-a902-4c21-abe2-6c7c0ef7d37c` | `9ca76093-d2de-469c-aed2-430296ff92df` |

The CLI upload messages include the full source commit and bead ID. Deployment status was read back for these exact IDs. Auth reached success before the business API upload. Both builds used their existing service-specific Dockerfiles and `golang:1.27.2`, resolved as `sha256:5bc7f572bbaa98885a3a1fd9c0aa76b59e3e14e8628bfc316bbfd0c701e4818c`. Read-only SSH inspection of `/service/auth` and `/service/schoolcrm` found the embedded Go build metadata `go1.27.2` and `golang.org/x/net v0.60.0` in both live binaries.

The business deployment retained `preDeployCommand: ["./admin migrate"]` and `/v1/readiness` health checks. No seed command ran. Uploads honoured `.railwayignore`, excluding secrets, development keys, local databases, dependency directories and agent state. Credentials and signing keys were not changed.

RAG `c3e22a03-b2c3-401e-8ab3-9f72c4283f76`, frontend `60e78418-3742-46f5-ba33-1b6afb3fcf78` and PostgreSQL `8b5695c5-a730-4f4b-a584-4bc22cc0ab77` remained on their existing successful deployments. No volume or database was replaced.

## Post-deployment acceptance

The public HTTPS endpoint remains `https://web-admin-staging-c78a.up.railway.app`. The current committed verifiers ran with credentials injected privately through the process environment. They emitted IDs and verdicts only. The Python verifier reused the existing frozen RAG dependency environment; it executed the script from the clean release checkout.

`scripts/verify_hosted_ai.py` passed authenticated login, approved-source retrieval and completed-request replay. The original school `ed78d739-d93a-431d-8957-da3e99c78cef`, department `495cf156-ee17-4d06-b271-c4ccbaa47fe4`, source `b1f3d4e7-1c1d-4811-a18e-b87a5ba6a7c5`, plan `495b0af4-b207-4cf9-b321-382e9ca5cc00` and request `b99268af-c194-47c1-80aa-07b731d2b815` were retained. Its printed summary reported two retrieved passages, model `DeepSeek-V4.1-Flash`, two citations, two versions, same-plan replay, wrong-school denial, missing-revision abstention and verified source, evidence and citation scope. Behind that summary, the script asserts that every passage and citation comes from the approved source with the Cloudflare BGE-M3 embedding identity, that every citation carries source checksum `eeea001277e1d952f4a6ebd4e42a21072cbb1e261d2e7228e50fe699ef041d0d`, that replay leaves the version ID list unchanged and that the generation receipt is completed. Replay reused the completed request ID and did not request new model generation.

`scripts/verify_staging_acceptance.mjs` passed on fresh probe `bf913bcd-ab5c-43a4-9c82-8725a91447fc`. Its printed summary listed the five denial names, the teacher, HOD and dean principals, current version 2, published version 1, reused plan `7945a0ac-3cf4-4cbf-94e5-6b5cba696f5d` and a verified source checksum. Behind that summary, the script asserts each denial's exact status with a JSON:API `errors` body and no `data`: private read (404), private reuse (404), premature publication (409), teacher review (403) and HOD dean approval (403). It also asserts that published version 1 records the distinct teacher author, HOD reviewer and dean approver beside human draft version 2, that the HOD reuse plan is an unpublished draft, and that the downloaded original PDF matches SHA-256 `eeea001277e1d952f4a6ebd4e42a21072cbb1e261d2e7228e50fe699ef041d0d`.

The probe and its two versions remain in staging as synthetic acceptance history. No audit or version records were deleted. The two exact deployment IDs still reported `SUCCESS` at the 01:35 UTC platform readback.

## Recovery boundary

The deployment history records the preceding credential-safe auth/business builds above. They predate the PR #79 patches, so rolling back to them reintroduces the unpatched Go toolchain and `golang.org/x/net`; prefer a patched forward fix. Application rollback must retain additive schema, original files, audit history, lesson versions, generation receipts and both existing volumes. The earlier backup/restore and browser/PDF evidence in [staging acceptance](railway-staging-acceptance.md) remains historical; those checks were not repeated by this stateless application refresh. No live restore or destructive rollback was performed.

The complete School AI implementation remains tracked in its open Beads.
