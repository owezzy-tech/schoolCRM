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
