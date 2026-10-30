# Secure staging administrator bootstrap

Bead `schoolCRM-kce.3`. The new `bootstrap-admin` command extends the existing Go admin tool; it is not a public HTTP endpoint or a seed command.

Provide `SCHOOLCRM_BOOTSTRAP_NAME`, `SCHOOLCRM_BOOTSTRAP_EMAIL` and masked `SCHOOLCRM_BOOTSTRAP_PASSWORD` as environment variables, together with the existing `SCHOOLCRM_DB_*` settings. Do not place credentials in command arguments or source files.

```sh
go run ./api/tooling/admin bootstrap-admin
```

The operator's bootstrap password must be 16-19 characters accepted by the existing application password type. That type permits letters, digits and `#@!-`. Generate a random19-character credential and store it in the selected Railway staging service's variables. Validation errors never include rejected password text.

The command locks the user table in a transaction and creates exactly one enabled `SUPER_ADMIN` only when it is empty. It reuses the canonical user business/store and bcrypt logic, and writes a credential-free `bootstrap_admin` audit through the canonical audit business/store in the same transaction. Actor UUID zero identifies the deployment operator's initial system action.

A retry for the exact already-audited administrator returns its ID without changing the password or adding another audit. Other populated databases are rejected; the command cannot upgrade an existing ordinary account. Audit failure rolls back account creation. An uncertain commit can safely retry the same administrator identity. No other users, school memberships or fixture credentials are seeded.

## Staging sequence and rollback

Target the SchoolCRM staging project/environment and its private database, never an existing production database. The chosen initial email is `owezzy@owenadirah.com`. Run the reviewed admin tool locally over an encrypted SSH forward, or from an image containing this command. Verify database identity, user count and extension/schema versions first.

The user-installed pgvector extension is a prerequisite. Apply the canonical curriculum table/index DDL without re-executing extension installation, then the explicit lesson/LangGraph migration. Enable the RAG runtime private database URL only after these tables exist. Keep model credentials, tokens and connection strings out of logs and model inputs.

Verify administrator login and permissions, audit count and safe same-identity retry. Retain a protected database backup and rehearse restoration to a separate database; never overwrite the live staging database for proof. Application rollback retains users, audits, additive tables and immutable history. If an initial account needs revocation, use the supported account-management workflow, not a direct password/role update hidden in a deployment script.

## Proof

Unit tests cover secret-safe validation, new-account commit, receipt replay, refusal to upgrade existing accounts, lock failure, user/audit rollback and uncertain commit. Opt-in PostgreSQL integration tests create only fresh UUID-named isolated schemas, run eight concurrent bootstrap attempts and prove one user/one audit, then test rollback when audit storage is unavailable. Test schemas are removed after proof; application tables are not touched.

### Hosted proof, 8 October 2026

The reviewed operator tool created administrator `15953eb5-2b0a-490b-b193-cbb2f3d90bc4` for the selected email. A second call returned the same ID. The database contained one user and one bootstrap audit. No bootstrap audit payload contained password or hash fields.

HTTPS authentication returned status 200 and the `SUPER_ADMIN` role. Browser sign-in opened the authenticated application and showed Owen Adirah's identity. The dashboard then displayed a data-load error. That failure remains a final hosted QA gate in `schoolCRM-kce.5`; sign-in does not prove the dashboard works.

The temporary credential is in the staging `schoolcrm` service variable `SCHOOLCRM_BOOTSTRAP_PASSWORD`. Retrieve it privately in Railway. Do not copy it into documentation, command arguments or logs.

PostgreSQL's `pg_basebackup` created `/var/lib/postgresql/schoolcrm-backup-proof-20261008`, owned by `postgres` with mode 0700. `pg_verifybackup` reported `backup successfully verified`. A separate copy started on port 5543 with TCP listening disabled and a private Unix socket under `/tmp/schoolcrm-restore-proof.U4WYek`.

A query against the restored cluster returned `1|1|0.8.7|8`: one user, one bootstrap audit, vector version 0.8.7 and all eight required curriculum, lesson and checkpoint tables. The restored cluster then stopped cleanly. The live database was not replaced or stopped. The protected backup and stopped proof copy remain available to the operator.

This is a recovery rehearsal, not an independent disaster-recovery backup. The backup shares the live database volume. Before production, choose an independent backup destination, retention and schedule. Restore into an isolated cluster first, verify the administrator and schemas, then approve any application connection switch. Never restore over the live database merely to repeat this proof.
