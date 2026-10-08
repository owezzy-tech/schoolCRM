# Railway staging inventory

Bead `schoolCRM-kce.2`, deployment epic `schoolCRM-kce`. Snapshot: 8 October 2026. This records provisioned infrastructure, not completed acceptance of lesson generation or administrator login.

Project: `174f33cc-0bb5-4d8b-88d8-1288f60792a6`.
Staging environment: `10e68a24-8d12-4210-8b23-ea532fb477a2`.
Production service instances remain empty.

Public frontend: <https://web-admin-staging-c78a.up.railway.app>.
Only the frontend has a public domain. PostgreSQL has no TCP proxy; Go debug and auth gRPC listeners bind loopback.

| Service | Service ID | Observed successful deployment |
| --- | --- | --- |
| Postgres | `da09d1cf-b68e-4fd3-82fb-44fbecde0136` | `8b5695c5-a730-4f4b-a584-4bc22cc0ab77` |
| auth | `2a0f7cac-8d37-4b93-89a4-f15d259f52d6` | `b89b9e41-25ae-46ad-84c9-7d7e47d61432` |
| schoolcrm | `e332dc5d-ea3c-41a5-ae01-b415fccd80aa` | `74c51663-701f-4992-b84b-8bb90e57d05b` |
| rag | `7acac89f-84e3-4418-8963-71e9a7657ef7` | `c74ce8cf-4afe-4eaf-82bb-8c9917c844c5` |
| web-admin | `cdbe91fd-b0ea-43fb-b6aa-a0501426da23` | `cf3f58b7-197f-434c-b1c0-d027e040f1ac` |

Application images use staging verification commits from [PR #73](https://github.com/owezzy-tech/schoolCRM/pull/73), still draft pending authenticated image-based proof. Auth/API/frontend builds used `4d89c83`; the RAG storage-evidence build used `a057236`. No feature was merged or promoted to production by these operations.

## Storage and credentials

PostgreSQL uses `pgvector/pgvector:pg18`. A new database named `schoolcrm` and a strong random password were created; no old database was copied. Its volume `4491ee85-318c-435f-ab26-d7da4ea6657b` is mounted at `/var/lib/postgresql`, with `PGDATA=/var/lib/postgresql/18/docker`.

RAG volume `75e59aae-6e77-4f76-a513-4e0dc07c3829` is mounted at `/service/var/files`. Both volumes report `READY`, with 5 GB capacity. All services have one configured/running replica in the platform's default `sfo` region. Each Go service has a maximum of ten database connections and two idle connections.

The RAG entrypoint in PR #73 prepares the root-owned mount, then starts the server as UID/GID 1000. The live startup log confirms `rag-container uid=1000 storage=writable`. It changes only the storage directory's ownership, not historical files recursively.

Auth uses a newly generated RSA-3072 signing key and staging-specific issuer, injected through Railway variables. The repository development keys are not used. Database credentials and private upstreams use Railway variable references. Cloudflare BGE-M3 and DeepSeek credentials were transferred through CLI stdin and verified without displaying values. Anonymous RAG access and external workflow tracing are disabled.

Private database traffic explicitly disables PostgreSQL TLS within Railway's private network; public traffic uses HTTPS at the Railway edge. No public database endpoint was created.

## Verified behaviour and remaining gates

- HTTPS `/` and `/sign-in` return 200, and the Angular sign-in form renders in a real browser.
- Proxied Go `/v1/readiness` returns 204. Unknown-user auth login returns 401.
- Protected Go user and RAG generation requests without credentials return JSON:API 401 responses.
- The initial scalar pre-deploy setting was omitted from the deployed manifest. A browser login caught the resulting undefined-table error despite healthy readiness. Configuring the native command list `preDeployCommand: ["./admin migrate"]` corrected that omission, and deployment `74c51663-701f-4992-b84b-8bb90e57d05b` records the command and reached `SUCCESS`. It never runs `seed` or `migrate-seed`. A readiness response alone is not schema proof.
- The user selected `owezzy@owenadirah.com` for the initial staging administrator. That account has not yet been created.
- The RAG runtime database URL remains unset until Python curriculum/checkpoint tables are migrated. Its successful readiness check alone does not prove retrieval or generation readiness.
- pgvector is available in the PostgreSQL image, but extension creation and the explicit Python migrations still require their separate gate. The Railway skill reserves extension installation for a user-run operation.
- A restart-persistence probe through `railway volume files` was attempted, but SSH authentication failed because the account has no registered key. No remote probe was uploaded and no volume data was deleted. Permission to register a dedicated account SSH key is pending; that permission has not been assumed.

Keep `schoolCRM-kce.2` open until private-network and persistence verification is complete. Authenticated browser/API journeys belong to the remaining `kce.1`, `kce.3`, `kce.4` and `kce.5` acceptance gates. Existing local PostgreSQL 17 and proof data remain untouched.

## Rollback boundary

Use the recorded per-service deployment history for application rollback. Retain both volumes and additive schema tables. Do not delete services, volumes, published history or signing keys as a rollback shortcut. A database backup/restore rehearsal is still required before the deployment epic can close.
