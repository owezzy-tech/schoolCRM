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

Shut it down:

```bash
make dev-down
```

## Images and deployment names

Local images use the `localhost/schoolcrm/*` prefix. Helm development images use `localhost:5001/schoolcrm/*`, and production chart defaults use `registry.example.com/schoolcrm/*`.

The default Kubernetes namespace is `schoolcrm-system`, and the default Kind cluster name is `schoolcrm-cluster`.

## Licensing

See [LICENSE](LICENSE).
