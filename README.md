<a id="readme-top"></a>

# SchoolCRM

SchoolCRM is a service-oriented school operations platform with Go APIs, an Angular administration portal, PostgreSQL persistence, authentication, metrics, and a Python RAG service.

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Angular](https://img.shields.io/badge/Angular-21-DD0031?style=for-the-badge&logo=angular&logoColor=white)](https://angular.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker Compose](https://img.shields.io/badge/Docker_Compose-local-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docs.docker.com/compose/)

## Table of Contents

- [About the project](#about-the-project)
- [Repository layout](#repository-layout)
- [Getting started](#getting-started)
- [Local development](#local-development)
- [Service endpoints](#service-endpoints)
- [Default users](#default-users)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## About the project

The repository contains the core APIs and the web administration application for school workflows. The backend uses layered Go services under `api/services`, shared infrastructure under `foundation`, and application/domain adapters under `app` and `business`. The RAG service is a separate FastAPI application with ports-and-adapters boundaries.

The local stack is defined in [`zarf/compose/docker_compose.yaml`](zarf/compose/docker_compose.yaml). A rendered overview is available in [`docs/diagrams/schoolcrm-local-architecture.html`](docs/diagrams/schoolcrm-local-architecture.html), with its source specification in [`docs/diagrams/schoolcrm-local-architecture.architecture.json`](docs/diagrams/schoolcrm-local-architecture.architecture.json).

## Repository layout

```text
api/services/              Go services and the Python RAG service
api/frontends/web-admin/   Angular administration portal
api/tooling/               Development and operational tooling
app/                       Transport and application adapters
business/                  Domain logic and persistence support
foundation/                Shared infrastructure primitives
zarf/compose/              Local Docker Compose stack
docs/                      Architecture and developer documentation
```

## Getting started

### Prerequisites

- Go 1.26 or newer
- Docker and Docker Compose
- Node.js and npm for the Angular admin application
- Python 3.11 or newer for standalone RAG development

### Install tooling and run checks

From the repository root:

```bash
make dev-gotooling
make dev-brew
make dev-docker
make test
```

## Local development

Start the Kubernetes-oriented development workflow:

```bash
make dev-up
make dev-update-apply
```

Build and start the local Docker Compose stack:

```bash
make compose-build-up
```

For a host-process workflow, inspect available targets and run each long-lived service in its own terminal:

```bash
make local-run-help
make local-seed
make local-auth
make local-schoolcrm
make local-rag
make local-web-admin
```

The detailed runbook is in [`docs/local-development.md`](docs/local-development.md).

Stop the development stack with:

```bash
make dev-down
```

## Service endpoints

┌──────────────────────┬────────────────────────────────────────────┐
│ Service              │ Host endpoint                              │
├──────────────────────┼────────────────────────────────────────────┤
│ Web admin            │ http://localhost:8080                      │
│ SchoolCRM API        │ http://localhost:3000                      │
│ SchoolCRM debug vars │ http://localhost:3010/debug/vars           │
│ Auth REST API        │ http://localhost:6000                      │
│ Auth gRPC            │ localhost:6001                             │
│ RAG API              │ http://localhost:4545                      │
│ PostgreSQL           │ localhost:5454                             │
│ Metrics              │ localhost:4000, :4010, :4020               │
└──────────────────────┴────────────────────────────────────────────┘

The Compose RAG container listens on port `7000` internally and is published as `4545` on the host. PostgreSQL uses `postgres` as the local development password unless overridden.

## Default users

Seed data creates these local development accounts. All use the password `gophers`.

┌────────────────────────────┬──────────────┐
│ Email                      │ Role         │
├────────────────────────────┼──────────────┤
│ superadmin@example.com     │ SUPER_ADMIN  │
│ admin@example.com          │ SCHOOL_ADMIN │
│ teacher@example.com        │ TEACHER      │
│ user@example.com           │ STUDENT      │
└────────────────────────────┴──────────────┘

Example login:

```bash
curl -i -X POST http://localhost:6000/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"superadmin@example.com","password":"gophers"}'
```

## Documentation

- [`docs/local-development.md`](docs/local-development.md) — local setup and service runbook
- [`docs/diagrams/schoolcrm-local-architecture.html`](docs/diagrams/schoolcrm-local-architecture.html) — local architecture visualization
- [`api/frontends/web-admin/README.md`](api/frontends/web-admin/README.md) — Angular development commands
- [`api/services/RAG/README.md`](api/services/RAG/README.md) — RAG architecture and standalone run instructions

## Contributing

1. Create a focused branch from the current development branch.
2. Run the relevant tests and checks before opening a pull request.
3. Document changes to service contracts, ports, migrations, or local workflows.
4. Open a pull request with a concise description and verification notes.

## License

See [`LICENSE`](LICENSE).

[Back to top](#readme-top)
