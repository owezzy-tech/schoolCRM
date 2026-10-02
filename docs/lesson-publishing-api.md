# Lesson publishing API

Bead `schoolCRM-o58.4` enforces the lesson approval and publication workflow defined in the [School AI ARD](school-ai-ard.md). It builds on the department capabilities from the [school access API](school-access-api.md). Lesson generation, the final lesson schema, reuse and PDF export belong to later beads.

## Workflow

Every save creates a new, immutable version. A plan points to its latest version (`currentVersion`) and, separately, to its live publication (`publishedVersion`). Workflow commands act only on the current version, so a revision immediately makes earlier versions ineligible for review, approval or publication.

| Command | Who | From | To |
| --- | --- | --- | --- |
| Create plan or revise | Plan author with `teach` | Any current state | New `draft` version |
| Submit | Version author with `teach` | `draft` | `hod_review` |
| Review | `review_lessons`, not the author | `hod_review` | `dean_approval` or `changes_requested` |
| Approve | `approve_lessons`, not the author or that version's reviewer | `dean_approval` | `approved` or `changes_requested` |
| Publish | Version author with `teach` | `approved` | `published` |

The author, HOD reviewer and dean approver of a version are always three different people, even when one user holds several capabilities. Holding `SUPER_ADMIN` grants no lesson authority. A changes-requested version cannot be resubmitted; the author creates a revision, which carries no earlier decisions. Publishing a newer version moves `publishedVersion`; until then, the earlier publication stays live and unchanged.

## Endpoints

All endpoints require a bearer token and return JSON:API v1.1 through `foundation/web.Respond`. Plans have type `lessonplan`; versions have type `lessonversion` with ID `<planID>:<version>`.

| Method and path | Body | Behaviour |
| --- | --- | --- |
| `GET /v1/schools/{school_id}/departments/{department_id}/lessons` | | Department plans, newest activity first |
| `POST /v1/schools/{school_id}/departments/{department_id}/lessons` | `title`, `content`, `changeSummary` | Create a plan with draft version 1 |
| `GET /v1/lessons/{plan_id}` | | Plan with current title, status and both version pointers |
| `GET /v1/lessons/{plan_id}/versions` | | Every version, newest first, with decisions |
| `POST /v1/lessons/{plan_id}/versions` | `baseVersion`, `title`, `content`, `changeSummary` | Create the next draft; `baseVersion` must be current |
| `POST /v1/lessons/{plan_id}/versions/{version}/submit` | | Submit for HOD review |
| `POST /v1/lessons/{plan_id}/versions/{version}/review` | `decision`, `feedback` | HOD decision |
| `POST /v1/lessons/{plan_id}/versions/{version}/approval` | `decision`, `feedback` | Dean decision |
| `POST /v1/lessons/{plan_id}/versions/{version}/publish` | | Publish; repeating a completed publication returns 200 and writes nothing |

`decision` is `approve` or `request_changes`; requesting changes requires `feedback` of at most 2000 characters. Titles contain 1 to 200 characters, change summaries at most 1000. Until `schoolCRM-o58.3` confirms the lesson schema, `content` is any JSON object of at most 256 KiB.

Reads require `teach`, `review_lessons` or `approve_lessons` in the plan's department. An unsubmitted draft is visible only to its author. Other readers see a plan's newest submitted or published version as its `currentVersion` and title, draft versions are left out of its history, and a plan that has only drafts returns 404 and is left out of department lists. Failures return 401 without a token, 400 for invalid input, 403 without current authority, 404 for unknown plans, versions or departments, and 409 for a stale version or a command that does not match the version's state.

## Consistency and audit

Each command locks the plan row and takes a shared lock on its school row, so it serialises with membership grants and revocations. It then reads the actor's enabled flag and current department capabilities inside the same transaction. A revocation therefore applies to the next command, even with an existing token. The version decision, plan pointers and an audit record holding the version's workflow metadata commit together; an audit failure rolls the command back.

Darwin migration `1.31` adds `lesson_plans` and `lesson_plan_versions`. The database also enforces the core rules: a trigger rejects any change to version content or authorship; checks keep author, reviewer and approver distinct and require both decisions before approval or publication; and deferred foreign keys keep both plan pointers on existing versions. Foreign keys retain history; no delete endpoint exists. Rolling back the application leaves these additive tables in place.

## Verification

```sh
go test -race -tags=integration -count=1 ./business/domain/lessonbus/... ./app/domain/lessonapp
```

These run against real PostgreSQL and authenticated HTTP. They cover the full lifecycle, the three-person rule, cross-department and cross-school denial, stale revisions and decisions, requested changes, revoked authority, concurrent publication retries, atomic audit rollback, database-level immutability and a publication waiting behind a committed revocation.
