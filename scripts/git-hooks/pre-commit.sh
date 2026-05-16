#!/usr/bin/env sh

set -eu

. "$(dirname "$0")/common.sh"

if frontend_changed_in_staged; then
  run_frontend_lint
fi

if rag_changed_in_staged; then
  run_rag_lint
fi

if go_changed_in_staged; then
  run_go_lint
fi
