#!/usr/bin/env sh

set -eu

# Beads' Git transport can invoke Git itself. Do not recursively sync.
if [ "${BEADS_GIT_SYNC_ACTIVE:-0}" = "1" ] || [ ! -d .beads ]; then
    exit 0
fi

if ! command -v bd >/dev/null 2>&1; then
    echo >&2 "Beads Git integration requires bd on PATH. Install Beads before committing or pushing."
    exit 1
fi

export BEADS_GIT_SYNC_ACTIVE=1
export BD_GIT_HOOK=1

hook=$1
shift

case "$hook" in
    pre-commit)
        bd hooks run pre-commit "$@"
        bd dolt commit
        ;;
    pre-push)
        bd hooks run pre-push "$@"
        bd dolt commit
        bd dolt push
        ;;
    post-merge)
        bd hooks run post-merge "$@"
        bd dolt commit
        bd dolt pull
        ;;
    post-rewrite)
        # Amend is local; a rebase can be the result of git pull --rebase.
        if [ "${1:-}" = "rebase" ]; then
            bd dolt commit
            bd dolt pull
        fi
        ;;
    post-checkout|prepare-commit-msg)
        bd hooks run "$hook" "$@"
        ;;
    *)
        echo >&2 "Unsupported Beads Git hook: $hook"
        exit 2
        ;;
esac
