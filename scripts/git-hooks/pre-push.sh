#!/usr/bin/env sh

set -eu

. "$(dirname "$0")/common.sh"

if frontend_changed_in_push; then
  run_frontend_test
fi

if rag_changed_in_push; then
  run_rag_test
fi

if go_changed_in_push; then
  run_go_push_checks
fi
