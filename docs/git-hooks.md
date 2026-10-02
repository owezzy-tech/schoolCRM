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

### Beads and Git

Beads already uses this repository's Git-backed Dolt remote. Issue data syncs separately
from the application branch; committing `.beads/config.yaml` alone does not upload tasks.
The local Dolt database and runtime files remain ignored by Git.

The active hooks are managed by Husky at `core.hooksPath=.husky/_`. Tracked `.husky/`
scripts call `scripts/git-hooks/beads.sh` after the existing quality gates:

- `pre-commit` runs the native Beads callback and commits pending task changes locally.
- `pre-push` commits remaining task changes and runs `bd dolt push`. A sync failure blocks
  the application push, so code is not pushed while its task data is unsynchronised.
- `post-merge` and a `post-rewrite` caused by rebase commit local task changes before
  `bd dolt pull`. A local amend does not trigger a network pull.
- `post-checkout` and `prepare-commit-msg` delegate to native Beads callbacks.
- A recursion guard prevents Beads' own Git transport from starting another sync.

Install a compatible `bd` CLI on PATH before working in this repository. The integration
was verified with Beads 1.1.2. Then run `npm install` at the repository root to activate
Husky, and verify the configured task remote:

```bash
bd dolt remote list
bd dolt commit
bd dolt pull
```

The configured origin is `git+https://github.com/owezzy-tech/schoolCRM`. Existing
`.beads/hooks/` files are inactive while Husky owns `core.hooksPath`; do not switch the
hook path or replace Husky with `bd hooks install`, which would bypass the quality gates.
Manual `bd dolt commit`, `bd dolt push` and `bd dolt pull` remain available for recovery.
A post-merge/rebase sync failure does not undo the completed Git operation; resolve the
reported Beads error and retry its sync explicitly.

Verify hook coordination without touching real task data or contacting a remote:

```bash
PYTHONDONTWRITEBYTECODE=1 python3 scripts/git-hooks/test_beads.py
```

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
npm --prefix api/frontends/web-admin run test -- --watch=false
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
- `schoolcrm`
- `metrics`
- `docs`
- `config`
- `deps`
- `ci`

Custom scopes are not accepted by the commit hook unless they are added to `commitlint.config.cjs`.

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
