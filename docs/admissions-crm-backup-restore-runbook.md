# Admissions CRM Backup and Restore Runbook

This runbook documents the restore process required by the Admissions CRM release gate.

Targets:

- RPO: less than or equal to 24 hours.
- RTO: less than or equal to 8 business hours.
- Backup cadence: daily automated database backups before production launch, then verified daily in production.
- Drill cadence: quarterly before production launch and semiannual after production launch.

## Systems in scope

- PostgreSQL database used by SchoolCRM services.
- Admissions CRM tables for constituents, applications, documents, checklist state, SIS sync state, audit state, campaigns, events, and reports.
- Document object storage metadata and any storage bucket or volume used by admissions uploads.
- Migration state managed by `business/sdk/migrate` and the admin tooling migrate/seed commands.

## Pre-backup checks

1. Confirm database readiness.

   ```bash
   kubectl -n schoolcrm-system exec statefulset/database -- pg_isready
   ```

2. Confirm service readiness.

   ```bash
   make readiness
   ```

3. Record the schema migration version from the migration table.
4. Record the current git commit deployed to the environment.
5. Confirm the last successful backup is less than 24 hours old.

## Backup procedure

1. Create a timestamped dump from the target environment.

   ```bash
   pg_dump --format=custom --file schoolcrm-YYYYMMDD-HHMM.dump "$DATABASE_URL"
   ```

2. Store the dump in the approved encrypted backup location.
3. Record metadata next to the artifact:

   ```json
   {
     "environment": "preprod",
     "createdAt": "YYYY-MM-DDTHH:MM:SSZ",
     "gitSha": "<deployed-sha>",
     "schemaVersion": "<darwin-version>",
     "postgresVersion": "18.3",
     "operator": "<name>"
   }
   ```

4. Verify the dump can be listed.

   ```bash
   pg_restore --list schoolcrm-YYYYMMDD-HHMM.dump
   ```

5. Mark backup status as verified only after step 4 succeeds.

## Restore procedure

Start the RTO timer before step 1.

1. Create an isolated restore database or namespace.
2. Restore the selected backup.

   ```bash
   pg_restore --clean --if-exists --no-owner --dbname "$RESTORE_DATABASE_URL" schoolcrm-YYYYMMDD-HHMM.dump
   ```

3. Run migrations against the restored database to align schema with the deployed application.

   ```bash
   ./admin migrate
   ```

4. Start application services pointed at the restored database.
5. Confirm service readiness.

   ```bash
   make readiness
   ```

6. Run admissions data verification queries:
   - Constituent count is greater than zero.
   - Active application count is greater than zero.
   - Document metadata count matches expected environment baseline.
   - Audit records exist for recent application transitions.
   - SIS sync jobs and retry-ready events are queryable.

7. Smoke test critical workflows:
   - Applicant portal status lookup.
   - Staff review workspace load.
   - Application detail load.
   - Document checklist load.
   - SIS sync operations screen load.

8. Stop the RTO timer after all readiness and smoke checks pass.

## RPO validation

1. Compare the selected backup timestamp with the incident declaration time.
2. Confirm the age is less than or equal to 24 hours.
3. If older than 24 hours, mark the drill or incident as an RPO failure and create a remediation issue.

## RTO validation

1. Record restore start and restore complete timestamps.
2. Confirm elapsed business time is less than or equal to 8 hours.
3. If elapsed time exceeds target, mark the drill or incident as an RTO failure and create a remediation issue.

## Rollback and abort criteria

Abort restore and escalate when:

- Backup integrity check fails.
- Restore database cannot start.
- Migrations fail and cannot be rolled forward safely.
- Authentication or authorization checks fail after restore.
- Admissions application, document, or audit data is missing from the restored environment.

## Required incident notes

```text
Incident or drill ID:
Environment:
Backup artifact:
Backup timestamp:
Incident declaration timestamp:
Restore start:
Restore complete:
RPO result:
RTO result:
Operator:
Approver:

Verification results:
- Database readiness:
- Service readiness:
- Constituent query:
- Application query:
- Document query:
- Audit query:
- SIS sync query:

Follow-up issues:
```
