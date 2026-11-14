#!/bin/sh
# Nightly SchoolCRM backup, run by schoolcrm-backup.timer (bead schoolCRM-5ho.9).
# Dumps the schoolcrm database and archives the RAG files volume, checks the
# dump reads back, encrypts both to the "SchoolCRM Backups" GnuPG key and
# uploads them to R2 under database/daily/<UTC time>/. The VPS holds only the
# public key: it can write backups but never read them back.
# Logs go to the journal:  journalctl -u schoolcrm-backup
set -eu

# R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_ENDPOINT and R2_BUCKET, from an R2
# API token with Object Read & Write on the SchoolCRM bucket only.
. /opt/schoolcrm/secrets/r2.env
# The remote is configured entirely from the environment, with no rclone.conf.
export RCLONE_CONFIG=/dev/null
export RCLONE_CONFIG_R2_TYPE=s3
export RCLONE_CONFIG_R2_PROVIDER=Cloudflare
export RCLONE_CONFIG_R2_ENDPOINT="$R2_ENDPOINT"
export RCLONE_CONFIG_R2_ACCESS_KEY_ID="$R2_ACCESS_KEY_ID"
export RCLONE_CONFIG_R2_SECRET_ACCESS_KEY="$R2_SECRET_ACCESS_KEY"
# The token can't list buckets, so don't let rclone check that the bucket exists.
export RCLONE_S3_NO_CHECK_BUCKET=true
# After an upload rclone re-reads the object by version id, which R2 rejects
# (501). Skip it: the upload's Content-MD5 already makes R2 check the bytes.
export RCLONE_S3_NO_HEAD=true

GNUPGHOME=/opt/schoolcrm/backup-gpg
export GNUPGHOME
# "SchoolCRM Backups" key (gpg encrypts to its encryption subkey); the private
# key is held off the VPS.
RECIPIENT=652B9E603ACD86BBB95787E7951DA669B31CB5A4
RAG_FILES=$(echo /var/lib/rancher/k3s/storage/pvc-*_schoolcrm_rag-files)
test -d "$RAG_FILES" || { echo "RAG files volume not found: $RAG_FILES" >&2; exit 1; }

stamp=$(date -u +%Y-%m-%dT%H%MZ)
prefix="r2:$R2_BUCKET/database/daily/$stamp"
umask 077
dir=$(mktemp -d /var/tmp/schoolcrm-backup.XXXXXX)
trap 'rm -rf "$dir"' EXIT

docker exec -u postgres postgres pg_dump --format=custom schoolcrm >"$dir/schoolcrm.dump"
docker exec -i -u postgres postgres pg_restore --list >/dev/null <"$dir/schoolcrm.dump"
tar -C "$RAG_FILES" -czf "$dir/rag-files.tar.gz" .

for f in schoolcrm.dump rag-files.tar.gz; do
  gpg --batch --yes --trust-model always --recipient "$RECIPIENT" \
    --output "$dir/$f.gpg" --encrypt "$dir/$f"
  rm -f "$dir/$f"
done
(cd "$dir" && sha256sum schoolcrm.dump.gpg rag-files.tar.gz.gpg >SHA256SUMS)

for f in schoolcrm.dump.gpg rag-files.tar.gz.gpg SHA256SUMS; do
  rclone copyto "$dir/$f" "$prefix/$f"
  remote=$(rclone lsf --format s "$prefix/$f")
  local_size=$(stat -c %s "$dir/$f")
  test "$remote" = "$local_size" || { echo "$f: R2 has ${remote:-nothing}, expected $local_size bytes" >&2; exit 1; }
done
echo "uploaded schoolcrm.dump.gpg and rag-files.tar.gz.gpg to $R2_BUCKET/database/daily/$stamp"
