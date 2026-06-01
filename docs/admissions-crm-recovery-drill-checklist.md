# Admissions CRM Recovery Drill Checklist

Use this checklist for quarterly pre-production drills and semiannual post-production drills.

## Drill setup

- [ ] Drill owner assigned.
- [ ] Environment selected.
- [ ] Backup artifact selected.
- [ ] Restore target isolated from production traffic.
- [ ] Current deployed git SHA recorded.
- [ ] Current schema migration version recorded.
- [ ] RTO timer owner assigned.

## Execute restore

- [ ] Start RTO timer.
- [ ] Validate backup artifact with `pg_restore --list`.
- [ ] Restore database into isolated target.
- [ ] Run migrations with admin tooling.
- [ ] Start application services against restored database.
- [ ] Confirm database readiness.
- [ ] Confirm service readiness.

## Validate RPO and RTO

- [ ] Backup timestamp is less than or equal to 24 hours before incident or drill start.
- [ ] Restore completed within 8 business hours.
- [ ] RPO result recorded.
- [ ] RTO result recorded.

## Admissions data verification

- [ ] Constituent lookup by email works.
- [ ] Constituent lookup by external SIS ID works when SIS-linked records exist.
- [ ] Application list and application detail load.
- [ ] Staff review workspace loads assigned applications.
- [ ] Document checklist loads metadata and status.
- [ ] Audit timeline contains application transition records.
- [ ] SIS sync jobs and events are visible.
- [ ] Retry-ready SIS sync events remain retryable after restore.

## Workflow smoke tests

- [ ] Applicant portal status lookup works.
- [ ] Staff reviewer can open an authorized application.
- [ ] Staff reviewer can draft a decision.
- [ ] Document reviewer can view checklist statuses.
- [ ] Operations user can open SIS sync operations screen.

## Accessibility gate spot-check after restore

- [ ] `/portal` keyboard script passes.
- [ ] `/portal/apply` keyboard script passes.
- [ ] `/applications/review` keyboard script passes.
- [ ] Document verification keyboard script passes.
- [ ] Error and status messages remain text-visible and screen-reader friendly.

## Drill closeout

- [ ] Failures and warnings documented.
- [ ] Remediation Beads issues created for failures.
- [ ] Runbook updates captured.
- [ ] Drill evidence attached to release record.
- [ ] Drill owner signs off.

## Result summary

```text
Drill ID:
Environment:
Owner:
Date:
Backup artifact:
Backup timestamp:
Restore start:
Restore complete:
RPO PASS/FAIL:
RTO PASS/FAIL:
Critical workflow PASS/FAIL:
Accessibility spot-check PASS/FAIL:

Open remediation issues:
```
