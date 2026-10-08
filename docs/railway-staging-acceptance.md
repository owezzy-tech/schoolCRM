# Railway staging acceptance and recovery

Bead `schoolCRM-kce.5`. The deployment is staging-only. Production promotion and migration of existing local data are excluded.

## Hosted workflow evidence

`scripts/verify_staging_acceptance.mjs` exercised real public HTTPS endpoints with separate synthetic teacher, HOD, dean and outsider identities. Inject `SCHOOLCRM_HOSTED_PROOF=staging`, `SCHOOLCRM_STAGING_ADMIN_PASSWORD` and `SCHOOLCRM_STAGING_QA_PASSWORD` privately. Do not pass credentials as arguments.

| Requirement | Direct result |
| --- | --- |
| Distinct author, reviewer and approver | Teacher `701bd8ad-cc6c-4062-a8ca-08e4986ba0df`, HOD `bf0c78cd-6933-43c9-8fab-cb6bc5f0d1ef`, dean `a9ed588d-5066-46b2-b2fe-6c725cf8123e`; stored version 1 contains their matching author/reviewer/approver IDs |
| Scoped access | Unassigned outsider denied source and generation-receipt reads; anonymous original download denied |
| Private drafts | HOD could not read or reuse the teacher's unsubmitted draft; teacher could not read the HOD's reused draft |
| Approval boundaries | Teacher review and HOD dean-approval attempts denied; the separate HOD and dean then completed their authorised stages |
| Publication and immutable edits | Plan `495b0af4-b207-4cf9-b321-382e9ca5cc00` has published version 1 and current draft version 2; the human edit did not replace version 1 |
| Reuse | HOD-owned plan `7945a0ac-3cf4-4cbf-94e5-6b5cba696f5d` is an unpublished draft with the selected source snapshot |
| Original integrity | Authenticated original download matched SHA256 `eeea001277e1d952f4a6ebd4e42a21072cbb1e261d2e7228e50fe699ef041d0d` |
| Browser and PDF | Synthetic teacher signed in through the real browser and opened `/lesson-export/495b0af4-b207-4cf9-b321-382e9ca5cc00/1`; the page displayed version 1 as the live publication |

The browser's native PDF renderer produced `/private/tmp/schoolcrm-staging-version1.pdf`, 158,957 bytes and five pages. PDF text readback contains version 1, the approved source ID and `DeepSeek-V4.1-Flash`, and excludes `Teacher-edited staging objective` from the newer draft. This is a real PDF of the authenticated historical snapshot, not an HTML-only assertion. The desktop print-dialog Save interaction was not automated.

The source and generation IDs, provider identities and initial real-model proof are recorded in PR #77 and `docs/railway-hosted-ai-proof.md` on its branch. Replaying the completed request after the human edit preserved both version IDs and returned the same plan. The original generated version still has two matching Cloudflare citations.

Independent specification review found that the initial script could skip denial checks when replaying an already-published fixture. The corrected script creates a fresh teacher-owned probe from the real published output on every run, without a model call. It must observe five denials before reporting completion: private read, private reuse, premature publication, teacher review and HOD dean approval. The first corrected hosted run used probe `95d31f5c-45e2-41cd-a8b7-7af9cb1c1be9` and observed all five, then completed distinct-principal approval/publication and preserved published version 1 beside human draft version 2. The original generated plan remains unchanged.

## Restart and restore

RAG deployment `c3e22a03-b2c3-401e-8ab3-9f72c4283f76` restarted through the native Railway command. SSH inspection found a fresh runtime process, age 77 seconds, running as UID 1000. Hosted approved-source retrieval and completed-request replay then passed with the same plan/source/request IDs and both lesson versions retained.

After the workflow, `pg_basebackup` created `/var/lib/postgresql/schoolcrm-final-backup-20261008`, owned by `postgres` with mode 0700. `pg_verifybackup` reported success. A separate copy under `/tmp/schoolcrm-final-restore-proof.gbTSOy` started on port 5544 with TCP listening disabled and its own private Unix socket.

The restored database query returned `5|2|2|1|1|5`: five synthetic/operator identities, two versions of the generated plan, current version 2, published version 1, one approved source with the exact Cloudflare identity and five durable LangGraph checkpoints. The restored cluster stopped cleanly. No live data or volume was replaced. The backup and stopped proof copy remain protected for operator inspection.

Both physical backups share the live database volume. They prove recovery mechanics, not independent disaster recovery. Choose an independent destination, retention and schedule before production; that choice is outside this staging rollout.

## Current release and rollback boundaries

The public entry point is `https://web-admin-staging-c78a.up.railway.app`. Auth, business API, RAG and PostgreSQL remain private. The database has no public TCP proxy. Signing keys, provider keys and passwords are Railway secret variables.

Current successful application deployments are auth `2d682d2a-dfae-4ffc-877a-e2783361c36d`, business API `9a750f32-a902-4c21-abe2-6c7c0ef7d37c`, RAG `c3e22a03-b2c3-401e-8ab3-9f72c4283f76` and frontend `60e78418-3742-46f5-ba33-1b6afb3fcf78`. Auth and business API contain the credential-safe source from PR #77. Runtime success is not a claim that this PR has been merged.

For recovery, inspect the exact affected service and deployment in staging first. Restart the existing successful release for a process-only failure. For a bad application release, select a recorded, previously validated deployment through Railway's deployment history and verify its status and protected HTTPS behaviour before accepting rollback. Do not restore an old auth image that logs bearer values or an old RAG image with writable privileged startup code.

Keep additive schema tables, immutable versions, source originals, audit records and generation receipts. Never run `seed` or `migrate-seed`, delete volumes, or overwrite the live database as a rehearsal. Restore a backup into an isolated cluster, compare identities/source/version/checkpoint state, then obtain approval for any live connection switch. No live destructive rollback was performed for this proof.

## Review and remaining findings

Independent standards/specification and architecture review is a final gate. GitHub has no configured check results for these PRs; do not describe absent checks as passing CI.

The separate follow-ups `schoolCRM-31e` and `schoolCRM-buq` cover template dashboard content and a shared admissions test-fixture race. They are not hidden as passing tests. Earlier token-bearing failure logs, if any, are historical records; PR #77 stops future token logging but does not erase log history.
