# SchoolCRM database migration preparation

`schoolCRM-5ho.8` records operational work performed on 30 October 2026 UTC
(31 October in Nairobi). The API/auth configuration is delivered separately in
[PR #91](https://github.com/owezzy-tech/schoolCRM/pull/91), which is open and
unmerged. The migration destination removes Railway completely.

## Verified state

Railway staging remains the authoritative writer. Its `schoolcrm` database runs
PostgreSQL 18.6 with pgvector 0.8.7. The source inventory found 52 public tables,
five users, one school, four memberships, nine lesson plans, 15 lesson versions
and saved AI checkpoints. The configured RAG directory `/service/var/files`
contained zero files. These counts were read while writers remained enabled;
they are not a final cutover snapshot.

The preliminary custom-format database backup is stored on `owezzy-db` at
`/opt/schoolcrm/backups/source-rehearsal-20261030.dump`, with mode 0600 inside a
0700 directory. It is 4,829,574 bytes. Its SHA-256 is
`e956756ddc0b8bae8f5cd67aed6ecd29a193463d7d7839f2a4df34176bf0e886`.
The archive catalog validates with the VPS's existing `pg_restore`.

The backup was restored into a disposable container using the shared server's
exact image, `sha256:508c5290cda481d4f5f846446a26e9c1b804766828a394a5861de1b348a18b4c`.
That image runs PostgreSQL 18.2 and pgvector 0.8.1. All 52 restored table counts
matched the preliminary inventory, with zero invalid public indexes. The
rehearsal used 512 MiB, half a CPU, temporary data storage, no network and no
published ports. It reported no OOM and was removed after verification.
This proves restore compatibility for this snapshot, not future use of newer
pgvector features or production traffic capacity.

## Isolated destination

The existing PostgreSQL container now joins the dedicated external Docker
network `schoolcrm-database`, subnet `172.20.0.0/16`, with alias `postgres`.
It was not recreated. A new `schoolcrm` database is owned by a new `schoolcrm`
login with a fresh credential, a 25-connection limit and no superuser,
role-creation, database-creation or replication privileges. Its pgvector
extension is installed by the PostgreSQL administrator. No application tables
or source records have been restored there yet.

Only the new database's public access was revoked. The Garage databases retain
their original owners and ACLs. PostgreSQL's existing `pg_hba.conf` was preserved
and extended with TLS-only SchoolCRM access from the new subnet, followed by
explicit rejection of other SchoolCRM database connections. The HBA parser
reported zero errors before reload. Live connection probes verified:

- `schoolcrm` connects to its own database over TLS.
- `schoolcrm` cannot connect to `garage` or `garagerag`.
- The `garage` login lacks `CONNECT` on the new SchoolCRM database.

Garage's databases still inherit their existing `PUBLIC` connection defaults;
the SchoolCRM cross-database restriction is enforced by HBA. No existing
application's ACLs were rewritten.

The network attachment survives a container restart. A container recreation
must include the installed `/opt/schoolcrm/postgres-network.compose.yml`
(source: `zarf/vps/postgres-network.compose.yml`) alongside the existing
`/opt/postgres/compose.yml`; otherwise Docker Compose can omit the extra network.
The overlay preserves the existing default network. Incorporate it into the
database operator's maintenance command before the next PostgreSQL recreation.
Compose parser comparison against the live base file confirmed that all service
fields except networks, and all volume definitions, remain identical.
Do not run a shared-database recreation as part of this migration.

## Secrets and receipts

The user explicitly approved a narrow direct transfer of signing material,
issuer/key ID and required provider credentials from Railway to
`/opt/schoolcrm/secrets`. No broad local variable export ran. Secret values were
not logged or committed. Environment files are root-owned mode 0600; the secret
directory is 0700. The signing key directory and key belong to UID 1000, with
modes 0700 and 0400 respectively, for the auth container's read-only mount.
Fresh Go database credentials are in `api.env` and `auth.env`.

Operational receipts are root-only files under `/opt/schoolcrm/backups`:

- `source-rehearsal-20261030.manifest.json`: archive size and checksum.
- `rehearsal-20261030.receipt.json`: restore counts, versions and isolation.
- `database-provision-20261031.receipt.json`: database/TLS/access evidence.
- `pg_hba.before-schoolcrm-20261031.conf`: preserved HBA configuration.

The temporary Railway SSH identity `schoolcrm-migration-20261030` remains
registered for final export. Revoke it, unload its agent identity and remove
its ignored local private key when migration recovery no longer needs it.

## Gates before cutover

No SchoolCRM application writer has started on the VPS. Stop all source writers
before taking and verifying the final database backup. Restore only into the
verified SchoolCRM destination; retain both data copies. PostgreSQL owns the
preinstalled vector extension, so an application-role restore must exclude
extension creation/comment entries from its reviewed archive selection. Verify
table ownership, constraints, counts and vector queries after restore.

Complete RAG and durable-file deployment under `schoolCRM-5ho.6`, independent
scheduled off-site backups and restore checks under `schoolCRM-5ho.9`, then
frontend/ingress and authenticated acceptance under `schoolCRM-5ho.2`.
The preliminary VPS-local dump does not provide ongoing off-site backup.
Existing-user login, school-scoped permissions, lesson history, streaming,
representative load and rollback remain live cutover acceptance conditions.
Railway retirement stays in `schoolCRM-5ho.7` until those conditions pass.
