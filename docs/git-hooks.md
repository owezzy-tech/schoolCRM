# Git Hooks and Commit Workflow

This repository uses **Husky** for git hooks and **Commitizen + commitlint** for conventional commit messages.

The goal is simple:

- run the right checks for the files you changed
- avoid running unrelated toolchains
- keep commit messages consistent

## Installed Tooling

At the repo root:

- `husky`
- `commitizen`
- `cz-git`
- `@commitlint/cli`
- `@commitlint/config-conventional`

Main config files:

```text
package.json
commitlint.config.cjs
.husky/pre-commit
.husky/pre-push
.husky/commit-msg
scripts/git-hooks/common.sh
scripts/git-hooks/pre-commit.sh
scripts/git-hooks/pre-push.sh
```

## Initial Setup

Run this once after cloning or after pulling hook changes:

```bash
npm install
```

If you work on the Python RAG service, also install its dev dependencies:

```bash
python3 -m pip install -e api/services/RAG[dev]
```

That installs:

- `pytest`
- `ruff`

Without those, the RAG hooks will fail when `api/services/RAG` files are changed.

## Hook Behavior

### Pre-commit

Runs only for **staged** files.

#### Frontend changes

If staged files are under:

```text
api/frontends/
```

the hook runs:

```bash
npm --prefix api/frontends/web-admin run lint
```

#### RAG changes

If staged files are under:

```text
api/services/RAG/
```

the hook runs:

```bash
python3 -m ruff check api/services/RAG
python3 -m compileall api/services/RAG
```

#### Go changes

If any staged file ends with `.go`, the hook runs:

```bash
CGO_ENABLED=0 go vet ./...
staticcheck -checks=all ./...
```

### Pre-push

Runs only for files changed between your branch and upstream.

#### Frontend changes

If changed files include:

```text
api/frontends/
```

the hook runs:

```bash
npm --prefix api/frontends/web-admin run test -- --watch=false --browsers=ChromeHeadless
```

#### RAG changes

If changed files include:

```text
api/services/RAG/
```

the hook runs:

```bash
python3 -m pytest api/services/RAG/tests -v
```

#### Go changes

If changed files include any `.go` file, the hook runs:

```bash
CGO_ENABLED=0 go test -count=1 ./...
govulncheck ./...
```

## Commit Messages

This repo enforces **conventional commits**.

Examples:

```text
feat(rag): add document ingestion endpoint
fix(frontend): correct login form validation
ci(config): add path-aware husky hooks
docs(docs): describe git hook workflow
```

### Recommended way to commit

Use the interactive prompt:

```bash
npm run commit
```

That starts Commitizen and helps you build a valid commit message.

### Direct commit messages are still checked

If you run:

```bash
git commit -m "feat(rag): add query endpoint"
```

the `commit-msg` hook still validates the message with `commitlint`.

## Scopes

Configured scopes include:

- `frontend`
- `rag`
- `go`
- `auth`
- `sales`
- `metrics`
- `docs`
- `config`
- `deps`
- `ci`

Custom scopes are also allowed when needed.

## Typical Workflow

### Frontend work

```bash
git add api/frontends/web-admin/src/...
git commit
```

Before commit:

- Angular lint runs

Before push:

- Angular tests run

### RAG work

```bash
git add api/services/RAG/...
git commit
```

Before commit:

- Ruff runs
- Python compile check runs

Before push:

- pytest runs

### Go work

```bash
git add app/... api/services/... business/...
git commit
```

Before commit:

- `go vet`
- `staticcheck`

Before push:

- `go test`
- `govulncheck`

## Troubleshooting

### `ruff` or `pytest` not found

Install the RAG dev dependencies:

```bash
python3 -m pip install -e api/services/RAG[dev]
```

### `staticcheck` or `govulncheck` not found

Install the Go tooling already referenced by the repo:

```bash
make dev-gotooling
```

### Frontend test fails because Chrome is unavailable

The hook uses Angular/Karma with `ChromeHeadless`. Make sure Chrome or Chromium is installed locally.

### Hooks are not running

Reinstall root dependencies:

```bash
npm install
```

That re-runs the Husky `prepare` step.

## Design Notes

These hooks are intentionally **path-aware**:

- frontend checks do not run for Go-only changes
- Go checks do not run for frontend-only changes
- RAG checks do not run unless `api/services/RAG` is touched

This keeps the workflow fast while still enforcing quality at commit and push time.
