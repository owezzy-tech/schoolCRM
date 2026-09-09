# Local Development Guide

This guide describes the supported local workflows for SchoolCRM. Run commands from the repository root unless a section says otherwise.

## Choose a workflow

Use Docker Compose when you want the complete stack with reproducible service images. Use the local targets when you are actively changing one or more services and need fast reloads. The Kubernetes-oriented targets are useful when validating the deployment workflow.

## Prerequisites

- Go 1.26 or newer
- Docker with Compose support
- Node.js and npm for `web-admin`
- Python 3.12 or newer for standalone RAG development

Install the repository tooling once:

```bash
make dev-gotooling
make dev-brew
make dev-docker
```

Run the checks for each toolchain from the repository root:

```bash
make test             # Go tests, vet, staticcheck, and vulnerability checks
npm run test:frontend # Angular unit tests
npm run test:rag      # RAG pytest suite
```

## Docker Compose

Build and start the stack:

```bash
make compose-build-up
```

The Compose definition is [`zarf/compose/docker_compose.yaml`](../zarf/compose/docker_compose.yaml). It starts PostgreSQL, migration and seed initialization, auth, SchoolCRM, RAG, the Angular web admin, and metrics on a shared bridge network.

Verify health endpoints:

```bash
curl http://localhost:3000/v1/liveness
curl http://localhost:6000/v1/liveness
curl http://localhost:4545/v1/liveness
```

Stop the Compose stack:

```bash
make compose-down
```

For the Kubernetes-oriented workflow, remove the KIND cluster with:

```bash
make dev-down
```

## Local service targets

List the available targets before starting a host-process workflow:

```bash
make local-run-help
```

Seed the database, then start each long-running process in its own terminal:

```bash
make local-seed
make local-auth
make local-schoolcrm
make local-rag
make local-web-admin
```

The local targets expect the configured PostgreSQL instance and service environment variables. If a target fails to bind, check that the corresponding Compose service is stopped and that the host port is free.

## Angular administration portal

The frontend lives in `api/frontends/web-admin`.

```bash
cd api/frontends/web-admin
npm install
npm start
```

The development server uses port `4400` and the configured proxy. For the production-like Compose image, use `make compose-build-up` instead. The Angular project README contains the complete build and test command reference.

## RAG service

The RAG service can run independently of the Compose image:

```bash
cd api/services/RAG
python3 -m venv .venv
source .venv/bin/activate
pip install -e .
uvicorn main:app --reload --host 0.0.0.0 --port 4545
```

The committed contract exposes `GET /v1/liveness`, `GET /v1/readiness`, document create/delete endpoints, and `POST /v1/rag/query`. See [`api/services/RAG/README.md`](../api/services/RAG/README.md) for the service architecture.

## Authentication and seed data

The migration-and-seed step creates four local users. Their password is `gophers`:

- `superadmin@example.com` — `SUPER_ADMIN`
- `admin@example.com` — `SCHOOL_ADMIN`
- `teacher@example.com` — `TEACHER`
- `user@example.com` — `STUDENT`

Request a token with:

```bash
curl -i -X POST http://localhost:6000/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"superadmin@example.com","password":"gophers"}'
```

## Troubleshooting

- Check `docker compose -f zarf/compose/docker_compose.yaml ps` when a container does not become healthy.
- Check the host port list in the root README before starting local processes.
- Re-run `make local-seed` after recreating the database volume.
- Inspect service logs with `docker compose -f zarf/compose/docker_compose.yaml logs -f <service>`.
