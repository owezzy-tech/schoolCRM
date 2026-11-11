# SchoolCRM on the owezzy-db k3s node

SchoolCRM runs on a single-node k3s cluster on the `owezzy-db` VPS (bead
`schoolCRM-5ho`). Railway is retired; the 30 October 2026 database dump is the
authoritative data (see `docs/deployment/2026-10-31-database-preparation.md`).
The manifests are the Kustomize overlay in `zarf/k8s/vps`, built on `zarf/k8s/base`.
The earlier Docker Compose files for this host were removed in `schoolCRM-5ho.11`.

## Host and isolation

The VPS has 4 vCPU, 7.9 GB RAM and about 87 GB free disk. It is shared:

- **GG-Garage's PostgreSQL and the portfolio's Neo4j run in Docker, outside the
  cluster.** Their containers, volumes, identities and limits must stay intact.
  Never recreate them as part of a SchoolCRM release.
- **k3s v1.37.2+k3s1** (`/etc/rancher/k3s/config.yaml`) has `servicelb` and
  `traefik` disabled, so ports 80/443 stay free: certbot renews
  `db.owenadirah.com` in standalone mode on port 80. Secrets are encrypted at rest.
  The kubelet reserves 4608 Mi for Docker, k3s and the OS, which leaves 2828 Mi
  for pods across all namespaces (SchoolCRM, GG-Garage, HRAS).
- **ufw** allows the pod (`10.42.0.0/16`) and service (`10.43.0.0/16`) ranges.
  The Kubernetes API (6443) and kubelet (10250) are not reachable from outside.
- **Namespace `schoolcrm`** enforces the `restricted` pod-security level. Every
  container runs as a non-root user with a read-only root filesystem, no
  privilege escalation and all capabilities dropped. A ResourceQuota (1 Gi
  requests, 2 Gi limits) keeps SchoolCRM inside its share. A NetworkPolicy
  admits traffic only from pods in the same namespace.
- **Nothing binds a host port.** Public traffic arrives through the outbound-only
  `cloudflared` connector in the namespace.

## Database access

The `schoolcrm` database and login role on the shared server have a fresh
credential, a 25-connection limit and no superuser, role-creation,
database-creation or replication rights. pgvector is installed by the database
administrator.

Pods reach PostgreSQL through the `postgres` ExternalName Service, which points
at `db.owenadirah.com` (a DNS-only record for the node). Docker's published port
re-addresses that traffic, so PostgreSQL sees every pod as `172.19.0.1`, the
`postgres_default` gateway. `/opt/postgres/pg_hba.conf` therefore allows:

```text
hostssl schoolcrm schoolcrm 172.19.0.1/32 scram-sha-256
```

That rule sits before the explicit rejections of the `schoolcrm` role, and the
connection must use TLS. `172.19.0.1` covers any process on the host, not only
pods; the password and TLS remain the boundary. The file is a single-file bind
mount: edit it in place (keep the inode), validate with `pg_hba_file_rules`,
then `select pg_reload_conf()`. Never restart or recreate the container for a
rule change.

Services set `*_DB_HOST=postgres:5432`, `*_DB_NAME=schoolcrm`,
`*_DB_USER=schoolcrm`, `*_DB_DISABLE_TLS=false` and at most 10 open connections.
The base manifests' `*_DB_HOST_PORT` names are not read by the services; the
overlay removes them.

## Access from a workstation

```sh
ssh -f -N -L 16443:127.0.0.1:6443 owezzy-db
export KUBECONFIG=~/.kube/owezzy-k3s.yaml   # mode 0600; full cluster admin
kubectl get pods -n schoolcrm
```

Close the tunnel afterwards (`pkill -f '16443:127.0.0.1:6443'`).

## Secrets

The secret files prepared on 30 October live in `/opt/schoolcrm/secrets`
(root-owned, 0700/0600). Create the Kubernetes Secrets **on the VPS** so values
never leave it or reach a terminal, log or the repository:

```sh
cd /opt/schoolcrm/secrets
k3s kubectl -n schoolcrm create secret generic api-env --from-env-file=api.env
k3s kubectl -n schoolcrm create secret generic auth-env --from-env-file=auth.env
k3s kubectl -n schoolcrm create secret generic auth-keys --from-file=auth-keys/
python3 - <<'EOF' | k3s kubectl -n schoolcrm create secret generic rag-env --from-env-file=/dev/stdin
from urllib.parse import quote
env = dict(l.split("=", 1) for l in open("api.env").read().splitlines() if "=" in l)
print(open("rag-provider.env").read().rstrip("\n"))
print("RAG_CURRICULUM_DATABASE_URL=postgresql://schoolcrm:"
      + quote(env["SCHOOLCRM_DB_PASSWORD"], safe="")
      + "@postgres:5432/schoolcrm?sslmode=require")
EOF
```

| Secret | Keys | Used by |
|---|---|---|
| `api-env` | `SCHOOLCRM_DB_PASSWORD` | API, `schoolcrm-migrate` |
| `auth-env` | `AUTH_DB_PASSWORD`, `AUTH_AUTH_ACTIVE_KID`, `AUTH_AUTH_ISSUER` | auth |
| `auth-keys` | signing key `.pem` files (mounted read-only at `/run/schoolcrm-keys`) | auth |
| `rag-env` | `RAG_CURRICULUM_DATABASE_URL`, Cloudflare and DeepSeek provider keys | RAG and its migrations |
| `web-admin-env` (optional) | `google_maps_api_key` (browser-visible; restrict by referrer) | web-admin |
| `cloudflared-token` | `token` from the Cloudflare tunnel | cloudflared |

Auth's issuer, active key ID and signing keys are the ones the Railway
deployment used, so existing tokens and the RAG service stay compatible. Do not
use the repository's development keys. Rotate a Secret with
`kubectl create secret ... --dry-run=client -o yaml | kubectl apply -f -` on the
VPS, then `kubectl rollout restart` the consumer.

## Release

Build images from a clean, reviewed commit; the tag is the commit's short SHA:

```sh
make vps-images                # linux/amd64 builds, imported into k3s over SSH
make vps-migrate               # Go ./admin migrate + RAG curriculum/lesson migrations
make vps-apply                 # platform, auth, API, RAG, web-admin; waits for rollouts
make vps-status
```

`vps-migrate` runs `./admin migrate` only; nothing on this host runs `seed` or
`migrate-seed`. Run it before `vps-apply` whenever the schema changes. Use
`VPS_TAG=<sha>` to redeploy or roll back to an earlier imported build.

## First deployment: data restore order

1. **Final backup.** The verified dump is
   `/opt/schoolcrm/backups/source-rehearsal-20261030.dump`
   (sha256 `e956756d…`), with an encrypted off-site copy in R2
   (`schoolcrm-backups`, see `schoolCRM-5ho.9`).
2. **Restore into the empty `schoolcrm` database** before any SchoolCRM pod
   writes. Restore as the PostgreSQL administrator with objects owned by the
   application role, and exclude the extension entries (the administrator owns
   pgvector):

   ```sh
   cd /opt/schoolcrm/backups
   docker exec -i postgres sh -c 'umask 077; cat > /tmp/source.dump' < source-rehearsal-20261030.dump
   docker exec postgres sh -c "pg_restore -l /tmp/source.dump | grep -vE ' EXTENSION | COMMENT - EXTENSION ' > /tmp/restore.list"
   docker exec -u postgres postgres pg_restore -d schoolcrm --no-owner --role=schoolcrm \
     --no-privileges --exit-on-error -L /tmp/restore.list /tmp/source.dump
   docker exec postgres rm -f /tmp/restore.list /tmp/source.dump
   ```

   Never use `--clean` against the shared server. Compare table counts with the
   receipt `final-source-decision-20261111.receipt.json`, check table ownership,
   constraints and a vector query.
3. **Create the Secrets**, then `make vps-images vps-migrate vps-apply`.
4. **Verify inside the cluster** before public traffic: readiness of every
   service, existing-user login through `/v1/auth`, school-scoped permissions
   and denied requests (JSON:API errors), published lesson history, and RAG
   streaming.
5. **Ingress.** Create the Cloudflare tunnel, store its token as
   `cloudflared-token`, point `schoolcrm.owenadirah.com` at
   `http://web-admin-service.schoolcrm.svc.cluster.local:80`, then
   `make vps-tunnel`. web-admin's Caddy routes `/v1/auth/*` to auth,
   `/v1/rag/*` to RAG and other `/v1/*` paths to the API.
6. Check memory and latency under representative use, and record the release
   SHA, image tags and cutover time in Beads.

## Rollback

- **Application:** `make vps-apply VPS_TAG=<previous sha>`, provided the
  previous images are compatible with the current schema. Migrations are
  additive; do not drop tables or published lesson versions.
- **Data:** the restore never touches GG-Garage databases. To redo a failed
  first restore, drop and recreate only the `schoolcrm` database as the
  administrator, then restore again from the verified dump.
- **Cluster:** `kubectl delete namespace schoolcrm` removes SchoolCRM's
  workloads. The RAG files volume uses `local-path-retain`, so its data stays
  on disk under `/var/lib/rancher/k3s/storage` until removed by hand.
  `/usr/local/bin/k3s-uninstall.sh` removes k3s entirely; it does not affect
  the Docker databases.

Never recreate the shared PostgreSQL container, delete its volume, or run
development seeds as a rollback shortcut. A recreation must keep
`/opt/schoolcrm/postgres-network.compose.yml`
(`zarf/vps/postgres-network.compose.yml`) until the `172.20.0.0/16` network and
its `pg_hba` rule, used only by the retired Compose setup, are removed.
