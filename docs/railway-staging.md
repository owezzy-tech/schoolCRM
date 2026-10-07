# Railway staging deployment

Deployment scope is Beads epic `schoolCRM-kce`: Railway staging, a new empty database and Cloudflare Workers AI BGE-M3 embeddings. Production promotion and migration of existing data are excluded.

## Services

Keep the repository root as the build context. Set each service's config-file path to the corresponding file under `zarf/railway`. Build from a reviewed commit on `develop`, not the dirty architecture checkout. Railway terminates public HTTPS; only `web-admin` needs a public domain. Auth, business APIs, RAG and PostgreSQL use private service networking. Do not publish Go debug or gRPC ports.

| Service | Config file | Internal port | Required wiring |
| --- | --- | --- | --- |
| auth | `zarf/railway/auth.toml` | 6000 | Database variables, fresh signing key and issuer |
| schoolcrm | `zarf/railway/schoolcrm.toml` | 3000 | Database variables, `SCHOOLCRM_AUTH_HOST=http://auth.railway.internal:6000` |
| rag | `zarf/railway/rag.toml` | 7000 | Auth/business URLs, database URL, provider credentials and persistent files |
| web-admin | `zarf/railway/web-admin.toml` | `PORT` | Private auth/business/RAG upstreams |

Set `PORT=6000`, `3000` and `7000` respectively for the fixed-port API services so Railway health checks target the listening port. Caddy reads `PORT` dynamically. Set `AUTH_WEB_DEBUG_HOST` and `SCHOOLCRM_WEB_DEBUG_HOST` to loopback-only addresses. Keep the auth gRPC listener private and publish no domain for it.

The frontend forwards `/v1/auth/*` to auth, `/v1/rag/*` to RAG and the remaining `/v1/*` to schoolcrm without rewriting paths. Set `AUTH_API_UPSTREAM=auth.railway.internal:6000`, `RAG_API_UPSTREAM=rag.railway.internal:7000` and `SCHOOLCRM_API_UPSTREAM=schoolcrm.railway.internal:3000`. All browser requests stay on the same HTTPS origin.

## Credentials and data

The default auth image has an empty key directory. Generate a new staging RSA key and supply it as the masked `AUTH_AUTH_KEYS_JSON` variable with the existing `{ "key": "<kid>", "pem": "<private PEM>" }` contract. Set matching `AUTH_AUTH_ACTIVE_KID` and a staging-specific issuer. Never use repository development keys. `make auth` explicitly selects the development target to retain the local workflow.

Use Railway variable references for database credentials. Configure `AUTH_DB_*` and `SCHOOLCRM_DB_*` separately, including host, name, password and explicit TLS policy for the chosen private PostgreSQL endpoint. The database must have pgvector support and durable storage. Installing an extension is a user-only operation under the Railway skill; verify it before applying the Python curriculum migration.

For RAG set `RAG_ALLOW_ANONYMOUS=false`, `RAG_AUTH_SERVICE_URL=http://auth.railway.internal:6000`, `RAG_SCHOOL_SERVICE_URL=http://schoolcrm.railway.internal:3000` and `RAG_CURRICULUM_DATABASE_URL` to the private PostgreSQL connection. Mount persistent storage at `/service/var/files`, writable by UID 1000, and set `RAG_FILE_STORAGE_DIR` accordingly. Verify mounted ownership before accepting uploads. Set `RAG_CURRICULUM_EMBEDDING_PROVIDER=cloudflare`, `RAG_CLOUDFLARE_ACCOUNT_ID`, a scoped `RAG_CLOUDFLARE_API_TOKEN`, and `RAG_DEEPSEEK_API_KEY` through secret variables. Do not upload `.env` files. Keep external workflow tracing disabled.

Apply Go migrations without development seed data, then the explicit Python curriculum and lesson migrations. Create the initial administrator through a supported secure bootstrap path. Never run `migrate-seed` in hosted staging. This bootstrap and database verification belongs to `schoolCRM-kce.3`.

## Release checks and rollback

`.railwayignore` excludes secrets, repository keys, local databases, source files uploaded by users, agent state and dependency directories. `.dockerignore` separately excludes environment secrets and local data from image build contexts. The RAG image installs the committed `uv.lock` without editable path dependencies.

Build all four Dockerfiles, inspect the default auth image for absent PEM files and check RAG imports from its installed environment. Validate Caddy configuration and verify that all three API route groups reach the expected private upstream. Liveness alone does not prove database/provider readiness; the deployment acceptance journey must test them explicitly.

Record the exact deployment ID and commit per service and observe Railway `SUCCESS`. Then test HTTPS login, school-scoped denials, approved curriculum retrieval, real Cloudflare embeddings and DeepSeek generation, retries, reuse, distinct HOD/dean approval, publication and PDF export. Restart services and prove original files, versions and checkpoints persist.

Rollback each application service to its recorded previous successful deployment. Keep additive schema tables and immutable records; do not drop history or roll databases back by deleting volumes. Record a database backup and restore rehearsal before accepting the deployment. A successful image build or upload is not a deployed-system verdict.
