# SchoolCRM Backend

SchoolCRM is a backend platform for managing school operations with service-oriented APIs, PostgreSQL storage, observability, and local Kubernetes/Docker workflows.

## Backend services

- `api/services/auth` — authentication and authorization service.
- `api/services/schoolcrm` — primary REST API service for SchoolCRM workflows.
- `api/services/metrics` — metrics collection service.
- `api/services/RAG` — Python RAG service for document ingestion and query workflows.

## Architecture

The backend follows a layered Go service architecture:

- `api/services` contains service entry points.
- `api/tooling` contains operational tooling.
- `app` contains transport and application adapters.
- `business` contains domain logic and persistence support.
- `foundation` contains shared infrastructure primitives.

The Go module is `github.com/owezzy/schoolCRM`.

## Development

Install local tooling and run checks from the repository root:

```bash
make dev-gotooling
make dev-brew
make dev-docker
make test
```

Run the local stack:

```bash
make dev-up
make dev-update-apply
```

Run with Docker Compose:

```bash
make compose-build-up
```

Host ports used by the compose stack:

- SchoolCRM API: `http://localhost:3000`
- Auth API: `http://localhost:6000`
- Auth gRPC: `localhost:6001`
- RAG API: `http://localhost:4545`
- Web admin: `http://localhost:8080`
- PostgreSQL: `localhost:5454`

Run services locally without Docker or Kubernetes:

```bash
make local-run-help
make local-seed
make local-auth
make local-schoolcrm
make local-rag
make local-web-admin
```

Run each long-running service target in its own terminal.

Shut it down:

```bash
make dev-down
```

## Default users

The seed data creates these default users for local development and tests:

┌──────────────────────────────────────┬────────────────────┬────────────────────────┬──────────────┬──────────┐
│ User ID                              │ Name               │ Email                  │ Role         │ Password │
├──────────────────────────────────────┼────────────────────┼────────────────────────┼──────────────┼──────────┤
│ b7f1d86f-0f1f-4c7b-84d7-99a570f14b6f │ Super Admin Gopher │ superadmin@example.com │ SUPER_ADMIN  │ gophers  │
│ 5cf37266-3473-4006-984f-9325122678b7 │ Admin Gopher       │ admin@example.com      │ SCHOOL_ADMIN │ gophers  │
│ c41fa5d3-d61f-45f1-b054-d2c7a3704019 │ Teacher Gopher     │ teacher@example.com    │ TEACHER      │ gophers  │
│ 45b5fbd3-755f-4379-8f07-a58d4a30fa2f │ User Gopher        │ user@example.com       │ STUDENT      │ gophers  │
└──────────────────────────────────────┴────────────────────┴────────────────────────┴──────────────┴──────────┘

Login through the auth service:

```bash
curl -i -X POST http://localhost:6000/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"superadmin@example.com","password":"gophers"}'
```

Legacy roles were migrated as follows:

- `ADMIN` -> `SCHOOL_ADMIN`
- `USER` -> `STUDENT`

## Images and deployment names

Local images use the `localhost/schoolcrm/*` prefix. Helm development images use `localhost:5001/schoolcrm/*`, and production chart defaults use `registry.example.com/schoolcrm/*`.

The default Kubernetes namespace is `schoolcrm-system`, and the default Kind cluster name is `schoolcrm-cluster`.

## Licensing

See [LICENSE](LICENSE).
