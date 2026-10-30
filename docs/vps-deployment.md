# SchoolCRM on the existing VPS

Bead `schoolCRM-5ho.5` prepares the Go API and auth service for `owezzy-db`.
The target frontend is Cloudflare Pages; the existing application still runs on
Railway. This configuration does not deploy the frontend or RAG service,
transfer data, change DNS, or stop Railway.

## Host and isolation

The 30 October 2026 capacity check found four vCPUs, 7,936 MiB RAM and 88 GiB free
disk. About 6,508 MiB RAM was available during that snapshot. PostgreSQL and Neo4j
belong to other applications. Their volumes, identities and existing limits must
remain intact. These measurements support a pilot; they are not a load test.

Reuse the existing PostgreSQL server with a separate `schoolcrm` database and
`schoolcrm` login role. Do not grant that role superuser, role-creation or
database-creation privileges. Do not grant it access to `garage` or `garagerag`.
PostgreSQL grants database `CONNECT` to `PUBLIC` by default. Check and enforce
database privileges as well as `pg_hba.conf`; a separate role alone is not
isolation. Do not change another application's database privileges without
checking its consumers.

The Compose project is named `schoolcrm`. API and auth use dedicated containers
and an external network named `schoolcrm-database`. Attach only the existing
`postgres` container to that network, with alias `postgres`. Creating the network
and attaching the container does not require recreating PostgreSQL. Add a TLS-only `pg_hba.conf`
rule for the SchoolCRM database/role and the new network's actual subnet. Validate
the rules with `pg_hba_file_rules`, then reload PostgreSQL. Preserve existing
rules. Other clients, including RAG, require their own reviewed access path.

API and auth have memory limits of 512 MiB and 256 MiB, CPU limits of 1 and 0.5,
bounded connection pools and capped logs. Debug and gRPC bind container loopback.
Only API port 6300 and auth port 6600 bind host loopback. They are intended for a
host-side Cloudflare Tunnel connector. Do not change the bindings to `0.0.0.0`.
The existing database and Neo4j limits plus these application limits still leave
host headroom, but measure performance under representative traffic before
cutover. Build images off-host to avoid competing with database workloads.

## Release inputs

Copy `zarf/vps/compose.yml` to `/opt/schoolcrm/compose.yml`. Keep the deployment
directory root-owned. Create `/opt/schoolcrm/secrets` with mode 0700. Place the
following files there through a secret channel:

- `api.env`: `SCHOOLCRM_DB_PASSWORD` for the dedicated role.
- `auth.env`: `AUTH_DB_PASSWORD`, `AUTH_AUTH_ACTIVE_KID` and `AUTH_AUTH_ISSUER`.
- `auth-keys/`: the source deployment's signing key files, readable by UID 1000
  inside auth. The bind mount is read-only and refuses to create a missing path.

Environment files must be mode 0600. Copy the existing issuer, active key ID and
signing keys from the source deployment so authentication remains compatible
during the migration. If the source stores JSON keys, export them securely using
the application's existing keystore format. Do not use repository development
keys. Give the destination role a fresh database credential. This does not revoke
the exposed source credential. Resolve its rotation/revocation separately under
`schoolCRM-i67`, after identifying its consumers. Keep the old deployment working
until cutover completes.

Set immutable image digests in `/opt/schoolcrm/deploy.env`, using
`zarf/vps/deploy.env.example` as the field reference. Build the existing
`zarf/docker/dockerfile.auth` default production target and
`zarf/docker/dockerfile.schoolcrm` for `linux/amd64` from the reviewed release
commit. Record that commit, image digests and build/test evidence in Beads. The
example image references are placeholders, not published images.

Never print `docker compose config`, container environments or connection URLs
with real secrets. Use `config --quiet` to validate the release configuration.

## Data preservation and release order

1. Restore Railway access and record the source deployment revisions, database
   version, schema versions, extensions, row counts and durable RAG file paths.
   Do not copy a database from the developer machine or seed demo users.
2. Export a consistent custom-format database backup and the RAG file volume
   through private channels. Verify checksums and restore a rehearsal copy.
   Record the source and rehearsal counts without exporting personal records to
   logs. Keep both the source and the backup.
3. Stop source application writes for final cutover. This includes Go API/auth
   writers and RAG generation, checkpoint and upload writers. A preliminary
   backup taken while the source accepts writes is only a rehearsal. Take a
   final backup after writes stop, restore into the isolated destination, and
   compare schema and counts. Run no destination application writers before
   the final restore is verified. Do not use `pg_restore --clean` against any
   existing shared database. Database recreation must target only an explicitly
   verified, disposable SchoolCRM rehearsal database.
4. Run the Go migrations explicitly, then verify the Python curriculum/run
   migrations and RAG storage before enabling any RAG traffic. The Go migration
   service runs `./admin migrate`; it never runs `seed` or `migrate-seed`.
5. Start API/auth on loopback. Verify readiness, schema versions, existing-user
   login, school-scoped permissions and published lesson history. Test denied
   requests and ensure JSON:API errors reach the browser. Readiness alone does
   not prove schema or application correctness.
6. Complete Cloudflare ingress, frontend API routing and RAG connectivity under
   `schoolCRM-5ho.2`. Configure `/v1/auth/*` to auth and other Go `/v1/*` paths to
   API without stripping the path. RAG needs its own explicit route. Keep auth,
   API and RAG on compatible signing keys/issuer during overlap. Validate TLS,
   CORS and streaming before changing the public hostname.
7. Enable destination writes only after the final read-only acceptance checks.
   Check representative load, memory, latency and database pool usage. Record
   the cutover timestamp and its single authoritative writer.

Commands below assume the database, network, image digests, keys and environment
files have already passed those gates. Run them in `/opt/schoolcrm`:

```sh
docker compose --env-file deploy.env config --quiet
docker compose --env-file deploy.env pull auth api migrate
docker compose --env-file deploy.env --profile maintenance run --rm --no-deps migrate
docker compose --env-file deploy.env up -d --wait auth api
curl --fail --silent --show-error http://127.0.0.1:6300/v1/readiness
curl --fail --silent --show-error http://127.0.0.1:6600/v1/readiness
```

## Rollback

Keep the old Railway deployment and volumes available. Before destination writes
are enabled, route traffic back to the source and re-enable its writers if a
cutover check fails. After destination writes begin, do not direct writes back
to the stale source database. Stop writes, reconcile the destination data into
the selected authoritative database, verify it, then change routing. Preserve
both data copies throughout recovery.

For an application-only rollback, restore the recorded previous image digests
and rerun `up -d --wait auth api` only if the old images are compatible with the
current schema. Retain additive tables and published lesson versions. Do not
delete volumes, drop tables, run development seeds or recreate the shared
PostgreSQL container as a rollback shortcut.

The VPS's existing nightly backup covers only `garage` and `garagerag`. Add and
verify a separate SchoolCRM backup/restore schedule before production cutover.
