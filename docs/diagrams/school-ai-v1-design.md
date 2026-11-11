# School assistant v1 diagram specification

These diagrams describe the proposed v1, not deployed agent functionality.
They adapt the EduPlan blueprint to the existing SchoolCRM repository.
The user confirmed the curriculum, approval and calendar decisions on 2026-10-01.

## Runtime ownership

Angular provides student, teacher and admin workspaces, chat, citations and approval cards.
Extend the existing authenticated application and user-management screens.
Python FastAPI hosts LangGraph, with a role-aware router selecting student adviser,
teaching assistant or operations assistant workflows. Agents are bounded workflows,
not separate services or sources of permission.
Both Go and Python validate identity through the existing Go authentication service.
Go owns school/department membership, resource authorisation and business rules.
Every retrieval and tool call resolves scope from trusted backend identity, never model text.
LlamaIndex operates within Python adapters for parsing, indexing and evidence retrieval.
Go owns deterministic qualification calculations, lesson persistence and publication,
class/timetable conflicts, attendance/exam aggregates, appointments and audit records.
Python agents call typed, allowlisted Go JSON:API tools using the authorised user's identity.
Write tools require an exact-payload approval, fresh permission/conflict checks and idempotency.
PostgreSQL stores school records, immutable lesson versions and business audit history.
Separate PostgreSQL workflow tables persist LangGraph checkpoints, private thread ownership
and pending approvals. Raw bearer tokens and provider credentials are not checkpoint content.
One configured model provider performs bounded reasoning and validated structured generation.
Its output cannot set official numeric points, grant access or bypass a publication gate.
The provider is not yet selected and should receive only scoped, permitted content.

## Curriculum sources and evidence

Kenya National Curriculum scope is Playgroup through secondary school.
Keep CBC/CBE and 8-4-4 framework identities and their applicable stages/revisions distinct.
CBC/CBE lesson design emphasises competencies and student-centred learning.
Cambridge scope is Early Years through Year 9; higher years are future expansion.
Cambridge lesson design supports inquiry and project-based learning.
Store curriculum, framework revision, stage/year, subject, topic, objectives and access scope.
Do not silently equate Kenyan grades with Cambridge years or reuse sources across frameworks.
Actual curriculum documents are still required; the user's scope description is not a syllabus.
Authorised staff upload source files, which retain checksums, versions and page/record citations.
Parsing/OCR preserves original text; uncertain extraction remains reviewable.
Source review precedes indexing; only approved, active resource revisions enter retrieval.
LlamaIndex indexes curriculum chunks into PostgreSQL with pgvector using matching model metadata.
The retrieval gate filters school, access, framework, level, subject and revision before model exposure.
Approved evidence returns as citations to LangGraph, which creates validated lesson drafts.
Existing Maseno admissions graphs remain separate historical evidence with reviewed rule projections.
Go evaluates supported qualification paths against reviewed rules and labels the public cluster estimate.
Cambridge, CBC/CBE and 8-4-4 qualification equivalence needs an explicit authoritative mapping.
Architecture proposals are excluded from student and curriculum retrieval collections.

## Agent action workflow

The user submits a request in Angular. Python authenticates and resolves authoritative scope.
LangGraph routes only to an allowed student, teacher or operations workflow.
Read-only tools retrieve scoped evidence or invoke Go calculation/report APIs.
Their result becomes a cited answer, eligibility breakdown, lesson draft or report.
Write requests first produce a structured proposal specifying resources, arguments and versions.
LangGraph checkpoints before requesting the required human approval in Angular.
Rejection ends the proposed action without a business write; approved actions resume the checkpoint.
Go revalidates permission, approval payload, resource versions and conflicts before execution.
Idempotent execution produces an audit receipt; conflicts or errors return a visible failure.
Refresh authority on every resumed run; changing roles invalidates previously sufficient authority.
Document text cannot register tools, select an approver or instruct an unapproved business change.

## Lesson version lifecycle

Teacher creates or edits a draft through Go; AI generation is optional draft assistance.
Teacher submits a specific immutable version to the assigned Head of Department.
Head of Department reviews that version and either requests changes or forwards it to the dean.
Dean approves every version, including the initial version, or requests changes.
Teacher is the publisher and can publish only the exact HOD-reviewed, dean-approved version.
Publication rechecks the teacher's school/class authority and the version's approvals.
The main lifecycle is Draft -> HOD review -> Dean approval -> Approved -> Published.
Requested changes return to teacher editing; resubmission restarts HOD review and dean approval.
Editing a published plan creates a new draft version; the published version remains immutable/live.
No approval carries over to changed content. Retention/deletion rules must preserve required history.
Represent HOD review and dean approval as scoped capabilities/memberships in the existing identity model.
The current global role enum does not itself implement these new scoped responsibilities.

## Calendar integration

Google Calendar is the initial provider, using user OAuth consent and protected credential storage.
The Go calendar integration checks availability and supports create/update/cancel,
recurrence, attendees, timezone-aware events and sync status, subject to Google account access.
Persist command keys and provider event IDs; recover ambiguous timeouts without duplicate appointments.
Use an outbox/retry worker for provider writes, keeping pending/failed/succeeded status visible.
Google Calendar API standard use currently has no additional API charge; quotas still apply.
Google documentation says quota excess charging is planned later in 2026; recheck before deployment.
Account/Workspace licensing is separate from API pricing.
Nextcloud Calendar is an open-source alternative with CalDAV, recurring events,
invitations, shared calendars and appointment features, but requires hosting and operations.
Its documented free/busy visibility is limited to users on the same Nextcloud instance.
Do not treat it as a drop-in provider with universal feature parity. It is not implemented in v1.

Research references checked on 2026-10-01:
- https://developers.google.com/workspace/calendar/api/guides/quota
- https://developers.google.com/workspace/calendar/api/guides/overview
- https://docs.nextcloud.com/server/latest/user_manual/en/groupware/calendar.html

## Current repository anchors

- Angular routes: api/frontends/web-admin/src/app/app.routes.ts:119.
- Existing user APIs: app/domain/userapp/route.go:32.
- RAG identity validation: api/services/RAG/infrastructure/auth.py:18.
- Scaffold assembly: api/services/RAG/infrastructure/dependencies.py:29.
- Current retrieval port lacks scoped filters: api/services/RAG/domain/ports/vector_store.py:20.
- Public KCSE approximation: business/domain/admissionsbus/values_ke.go:282.
- Source graphs and review issues: docs/RAGData/knowledge-graph/REPORT.md:1.
- PostgreSQL deployment: zarf/compose/docker_compose.yaml:3.

Delivery tracking: Beads task schoolCRM-rii under epic schoolCRM-o58.
