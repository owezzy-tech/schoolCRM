#!/usr/bin/env sh

set -eu

go_package_scopes() {
  printf '%s\n' "./api/services/... ./api/tooling/... ./app/... ./business/... ./foundation/..."
}

rag_python() {
  if [ -x "api/services/RAG/.venv/bin/python" ]; then
    printf '%s\n' "api/services/RAG/.venv/bin/python"
    return
  fi

  printf '%s\n' "python3"
}

run_rag_uv() {
  (
    cd api/services/RAG
    uv run "$@"
  )
}

rag_uses_uv() {
  command -v uv >/dev/null 2>&1 && [ -f "api/services/RAG/uv.lock" ]
}

frontend_changed_in_staged() {
  git diff --cached --name-only --diff-filter=ACMR | grep -q '^api/frontends/'
}

frontend_changed_in_push() {
  changed_files_between_refs | grep -q '^api/frontends/'
}

rag_changed_in_staged() {
  git diff --cached --name-only --diff-filter=ACMR | grep -q '^api/services/RAG/'
}

rag_changed_in_push() {
  changed_files_between_refs | grep -q '^api/services/RAG/'
}

go_changed_in_staged() {
  git diff --cached --name-only --diff-filter=ACMR | grep -q '\.go$'
}

go_changed_in_push() {
  changed_files_between_refs | grep -q '\.go$'
}

changed_files_between_refs() {
  upstream_ref=""
  if git rev-parse --abbrev-ref '@{upstream}' >/dev/null 2>&1; then
    upstream_ref=$(git rev-parse --abbrev-ref '@{upstream}')
  fi

  if [ -n "$upstream_ref" ]; then
    git diff --name-only "$upstream_ref"...HEAD
    return
  fi

  if git rev-parse --verify HEAD^ >/dev/null 2>&1; then
    git diff --name-only HEAD^...HEAD
    return
  fi

  git diff --cached --name-only --diff-filter=ACMR
}

run_frontend_lint() {
  echo "▶ Frontend changed: running Angular lint"
  npm --prefix api/frontends/web-admin run lint
}

run_frontend_test() {
  echo "▶ Frontend changed: running Angular tests"
  npm --prefix api/frontends/web-admin run test -- --watch=false --browsers=ChromeHeadless
}

run_go_lint() {
  echo "▶ Go files changed: running go vet and staticcheck"
  scopes=$(go_package_scopes)
  CGO_ENABLED=0 go vet $scopes

  if ! go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all $scopes; then
    echo "⚠️ staticcheck failed. If this is a tooling issue, verify your local Go installation supports the repo's Go version."
    exit 1
  fi
}

run_go_test() {
  echo "▶ Go files changed: running go test"
  scopes=$(go_package_scopes)
  CGO_ENABLED=0 go test -count=1 $scopes
}

run_go_push_checks() {
  run_go_test
  echo "▶ Go files changed: running govulncheck"
  scopes=$(go_package_scopes)
  go run golang.org/x/vuln/cmd/govulncheck@latest $scopes
}

run_rag_lint() {
  echo "▶ RAG files changed: running Ruff and Python compile validation"

  if rag_uses_uv; then
    run_rag_uv ruff check .
    run_rag_uv python -m compileall .
    return
  fi

  rag_py=$(rag_python)
  if ! "$rag_py" -m ruff --version >/dev/null 2>&1; then
    echo "❌ Ruff is required for RAG linting. If you use uv, run: (cd api/services/RAG && uv sync --extra dev). Otherwise create api/services/RAG/.venv and install dev deps."
    exit 1
  fi

  "$rag_py" -m ruff check api/services/RAG
  "$rag_py" -m compileall api/services/RAG
}

run_rag_test() {
  echo "▶ RAG files changed: running pytest"

  if rag_uses_uv; then
    run_rag_uv pytest tests -v
    return
  fi

  rag_py=$(rag_python)
  if ! "$rag_py" -m pytest --version >/dev/null 2>&1; then
    echo "❌ Pytest is required for RAG tests. If you use uv, run: (cd api/services/RAG && uv sync --extra dev). Otherwise create api/services/RAG/.venv and install dev deps."
    exit 1
  fi

  "$rag_py" -m pytest api/services/RAG/tests -v
}
