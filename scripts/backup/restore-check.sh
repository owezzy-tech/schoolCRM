#!/bin/sh
# Restore check for a SchoolCRM backup (bead schoolCRM-5ho.9). Run it on the
# machine that holds the "SchoolCRM Backups" private key, never on the VPS.
#
#   scripts/backup/restore-check.sh <dir>
#
# <dir> holds one night's objects downloaded from R2 database/daily/<time>/:
# schoolcrm.dump.gpg, rag-files.tar.gz.gpg and SHA256SUMS. The script checks
# the checksums, decrypts (gpg asks for the passphrase), restores the dump into
# a throwaway PostgreSQL container and prints each table's row count, to
# compare with production. Decrypted files stay in a temporary directory that
# is removed on exit; the container is removed too.
set -eu

src=${1:?usage: restore-check.sh <dir with schoolcrm.dump.gpg, rag-files.tar.gz.gpg, SHA256SUMS>}
image=pgvector/pgvector:0.8.1-pg18   # same PostgreSQL and pgvector as owezzy-db
name=schoolcrm-restore-check

(cd "$src" && shasum -a 256 -c SHA256SUMS)

umask 077
work=$(mktemp -d)
cleanup() { docker rm -f "$name" >/dev/null 2>&1 || true; rm -rf "$work"; }
trap cleanup EXIT

gpg --output "$work/schoolcrm.dump" --decrypt "$src/schoolcrm.dump.gpg"
gpg --output "$work/rag-files.tar.gz" --decrypt "$src/rag-files.tar.gz.gpg"
echo "RAG files: $(tar -tzf "$work/rag-files.tar.gz" | grep -vc '/$') files"

docker run -d --name "$name" -e POSTGRES_PASSWORD=restore-check -e POSTGRES_DB=schoolcrm "$image" >/dev/null
# Wait on TCP: the image's first-start init server answers on the socket only.
until docker exec "$name" pg_isready -h 127.0.0.1 -U postgres -d schoolcrm >/dev/null 2>&1; do sleep 1; done
docker exec -i "$name" pg_restore -U postgres -d schoolcrm --no-owner --no-privileges --exit-on-error <"$work/schoolcrm.dump"

# Exact row counts, the same query as the production comparison in the runbook.
docker exec "$name" psql -U postgres -d schoolcrm -At -F ' ' -c "
  select format('%s.%s %s', schemaname, relname,
         (xpath('/row/c/text()', query_to_xml(format('select count(*) as c from %I.%I', schemaname, relname), false, true, '')))[1]::text)
  from pg_stat_user_tables order by schemaname, relname"
