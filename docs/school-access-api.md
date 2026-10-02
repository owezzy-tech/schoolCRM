# School access API

Bead `schoolCRM-o58.4.1` provides the school membership foundation for lesson publishing. It does not implement lesson creation, version history or approval transitions. The [School AI ARD](school-ai-ard.md) defines those features.

## Authority

`SUPER_ADMIN` may create schools, create departments and grant or revoke school administration. Management delegation requires the target user to have the current `SCHOOL_ADMIN` role. That global role alone grants no school authority.

A delegated school administrator may create departments and grant or revoke teaching, HOD-review and dean-approval capabilities within their assigned school. They cannot grant or revoke management delegation, including their own, or operate in another school. Capabilities are explicit memberships, not names inferred from a user's free-text department.

| Capability | Scope | Meaning |
| --- | --- | --- |
| `manage_members` | School, with no department | Delegated membership administration |
| `teach` | One school department | Teaching responsibility |
| `review_lessons` | One school department | HOD review responsibility |
| `approve_lessons` | One school department | Dean approval responsibility |

The capability records define responsibilities. The subsequent lesson workflow must enforce three different users for author, HOD reviewer and dean approver on every version, even when a user holds multiple responsibilities. Holding `SUPER_ADMIN` does not replace these version-specific approvals.

## Endpoints

All endpoints require a bearer token from the existing authentication service. Mutation bodies are plain JSON request objects, as in existing SchoolCRM APIs. Responses go through `foundation/web.Respond` and use JSON:API v1.1. Single resources are in `data.attributes`, collections in `data` with `meta`, and safe failure messages in `errors[0].detail`.

| Method and path | Behaviour |
| --- | --- |
| `GET /v1/schools` | All schools for SUPER_ADMIN; active membership schools for other users |
| `POST /v1/schools` | Create a school. Body: `{"name":"Example School"}`. SUPER_ADMIN only |
| `GET /v1/schools/{school_id}/departments` | List departments for a current member of that school or SUPER_ADMIN |
| `POST /v1/schools/{school_id}/departments` | Create a department. Body: `{"name":"Sciences"}`. Scoped manager or SUPER_ADMIN |
| `GET /v1/schools/{school_id}/memberships` | List current and revoked memberships. Scoped manager or SUPER_ADMIN |
| `POST /v1/schools/{school_id}/memberships` | Grant or reactivate one capability. Body contains `userID`, `capability`, and `departmentID` for department capabilities |
| `DELETE /v1/schools/{school_id}/memberships/{membership_id}` | Revoke a capability, returning 204. Repeating a revocation is harmless |

Names contain 1 to 200 characters after trimming. IDs must be non-zero UUIDs. Mutation bodies have a 16 KiB limit and reject unknown fields and multiple JSON objects. A department must belong to the supplied school. A management grant must omit `departmentID`; the other capabilities require it. Granting to a disabled or missing user fails.

Grants reuse the same membership ID for the same school, department, user and capability. Revocation sets `active` to false and retains history. Collections are currently unpaginated and return the complete permitted scope; no page/filter parameters are accepted by this contract.

## Transactions and migration

Permission decisions read current user roles/enabled state and current memberships from PostgreSQL, rather than relying on token roles or the user cache. School mutations lock the school row; user checks take a shared row lock. This serialises grant/revocation decisions and prevents a concurrent role change from taking effect halfway through a committed operation. A revoked delegation or removed role prevents the next management command using an existing token.

Every successful mutation writes an actor/resource audit record through the existing audit store in the same transaction. Audit failure rolls back the business write. SQL constraints enforce capability scope, same-school departments and unique memberships, including school-scoped grants with a null department.

Darwin migration `1.30` adds `schools`, `school_departments` and `school_memberships`. It leaves existing users and admissions records unchanged. Foreign keys retain membership history by preventing referenced schools, departments and users from being deleted. Disable users and revoke memberships instead. No destructive history endpoint is provided.

Rollback can deploy the preceding application version while retaining these additive tables. Do not drop them after use without an approved data-retention decision. The feature is wired into the default and `crud` service builds; the reporting-only build stays read-only.

## Verification

Run the real PostgreSQL and authenticated HTTP contract tests:

```sh
go test -tags=integration -count=1 ./business/domain/schoolbus/stores/schooldb ./app/domain/schoolapp
go test -race -tags=integration -count=1 ./business/domain/schoolbus/stores/schooldb ./app/domain/schoolapp
```

These cover delegated management, cross-school denial, incorrect department scope, current-role checks, disabled users, duplicate concurrent grants, revocation, atomic audit rollback, real authentication and JSON:API success/error/collection responses. Apply the affected repository lint, test, build and vulnerability gates before the PR.
