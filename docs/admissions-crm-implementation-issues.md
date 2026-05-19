# Admissions CRM v1 Implementation Issue Breakdown

This issue breakdown turns the Admissions CRM v1 PRD and architecture decisions into independently grabbable implementation slices. Each issue should produce a vertical, testable increment.

Related planning artifacts:

- [`docs/admissions-crm-v1-prd.md`](admissions-crm-v1-prd.md)
- [`docs/admissions-crm-architecture.md`](admissions-crm-architecture.md)
- [`docs/diagrams/admissions-crm-bounded-contexts.excalidraw`](diagrams/admissions-crm-bounded-contexts.excalidraw)
- [`docs/diagrams/admissions-crm-application-workflow.excalidraw`](diagrams/admissions-crm-application-workflow.excalidraw)

## Milestone 1: Constituent and Application Core

### Issue 1: Add admissions bounded context scaffold

Goal: create the package/module structure for admissions CRM without implementing all behavior.

Deliverables:

- Admissions domain package.
- Initial route/API registration seam.
- Storage/repository interfaces.
- Basic health or placeholder endpoint if needed.

Acceptance criteria:

- The application builds.
- The admissions package does not depend on UI or infrastructure details.
- The new context has clear naming: constituent, inquiry, application, checklist, document, decision.

### Issue 2: Implement Constituent aggregate

Goal: model durable person identity for admissions.

Deliverables:

- Constituent entity with required identity fields.
- Lifecycle stage enum.
- Create/update/query use cases.
- Persistence migration.

Acceptance criteria:

- Constituent requires first name, last name, date of birth, primary email, and primary phone.
- Lifecycle stage supports prospect, inquiry, applicant, admitted, enrolled, alumni.
- Tests cover validation and lifecycle changes.

### Issue 3: Implement duplicate detection queue

Goal: detect exact and fuzzy duplicate candidates before conflicting records are created.

Deliverables:

- Exact match checks for verified email, external SIS ID, and portal account.
- Fuzzy candidate model and review queue.
- Staff resolution actions: merge, reject, link, defer.

Acceptance criteria:

- Exact trusted matches can auto-link where safe.
- Fuzzy matches never auto-merge without staff approval.
- Duplicate decisions are audit-logged.

### Issue 4: Add Program and AcademicTerm reference models

Goal: represent SIS-owned program and term data as read-only CRM reference data.

Deliverables:

- Program and AcademicTerm entities.
- External SIS ID fields.
- Read-only admin/API behavior.
- Local display metadata fields.

Acceptance criteria:

- CRM cannot create authoritative programs/terms outside sync/import paths.
- Applications must reference active term/program records.

### Issue 5: Implement Application aggregate and uniqueness invariant

Goal: support one active program application per constituent, term, and program.

Deliverables:

- Application entity.
- Application type enum: freshman, transfer, graduate.
- Active application uniqueness constraint.
- Create draft application use case.

Acceptance criteria:

- A constituent cannot have two active applications for the same term/program.
- Closed applications do not block future applications.
- Tests cover duplicate active application prevention.

### Issue 6: Implement application state machine

Goal: enforce allowed transitions and audit application workflow changes.

Deliverables:

- State transition table.
- Transition command/use case.
- Transition audit history.
- Side-effect hooks for notifications and sync events.

Acceptance criteria:

- Invalid transitions are rejected.
- Every transition records actor, timestamp, reason/note, from-state, and to-state.
- Tests cover happy paths and invalid transitions.

## Milestone 2: Applicant Portal and Documents

### Issue 7: Implement anonymous inquiry forms

Goal: let anonymous prospects submit inquiries and create/match constituents.

Deliverables:

- Public inquiry endpoint/form.
- Source/UTM capture.
- Constituent create-or-match behavior.

Acceptance criteria:

- Anonymous inquiry does not require login.
- Duplicate detection runs before creating a new constituent.
- Inquiry source attribution is stored.

### Issue 8: Implement applicant identity/profile context

Goal: connect authenticated portal users to constituents without separate account tables.

Deliverables:

- ApplicantProfile model.
- Identity/User to Constituent link.
- Portal access guard.

Acceptance criteria:

- Staff and applicant accounts share the same identity table.
- Applicant-only actions require authenticated applicant context.

### Issue 9: Implement application form templates

Goal: support mostly fixed forms with configurable required fields and checklist items.

Deliverables:

- ApplicationFormTemplate entity.
- Required field rules.
- Checklist template link.
- Versioning/active flag.

Acceptance criteria:

- Freshman, transfer, and graduate templates can differ.
- Required fields can be changed without code deployment.
- Core identity/program/term/status fields remain first-class.

### Issue 10: Implement checklist and document intake

Goal: collect and verify admissions documents.

Deliverables:

- ChecklistItem model.
- Document metadata model.
- Upload endpoint and storage adapter.
- Verification statuses and reviewer actions.

Acceptance criteria:

- Documents are tied to application/checklist item.
- Accepted/rejected/waived states are supported.
- Document view/download/verification is audited.

### Issue 11: Implement application fee payment tracking

Goal: track application fees without building full event payments.

Deliverables:

- ApplicationFee model.
- Payment status transitions.
- Waiver fields.
- Payment provider seam.

Acceptance criteria:

- Application fee can be pending, paid, failed, waived, refunded, or not required.
- Fee waiver records actor and reason.

## Milestone 3: Staff Workflows

### Issue 12: Add admissions CRM roles and permissions

Goal: add context-specific admissions roles while preserving global school roles.

Deliverables:

- Admissions role model.
- Action-based permissions.
- StaffProfile model.

Acceptance criteria:

- Global roles grant platform access.
- Admissions roles grant admissions actions.
- Permissions are not hardcoded only to sidebar menus.

### Issue 13: Implement assignment rules and manual override

Goal: assign leads/applications using configurable rules.

Deliverables:

- AssignmentRule model.
- Rule evaluation by priority.
- Manual override action.

Acceptance criteria:

- Rules support territory, program interest, application type, term, lead source, event attendance, and alpha range.
- Auto and manual assignments are audited.

### Issue 14: Implement rule-based lead scoring

Goal: score prospects using explainable rules.

Deliverables:

- LeadScoreRule model.
- Score recalculation use case.
- Score bands.

Acceptance criteria:

- Rules can be enabled/disabled.
- Score band is derived from total score.
- Staff can see why a score was assigned.

### Issue 15: Build staff review workspace

Goal: let reviewers manage assigned applications and decisions.

Deliverables:

- Assigned applications list.
- Application detail view.
- Document checklist review.
- Decision entry flow.

Acceptance criteria:

- Reviewers can only see authorized applications.
- Decisions create state transitions and audit records.

## Milestone 4: Communications, Events, and Notifications

### Issue 16: Implement communication log

Goal: capture email, SMS, and phone-call interactions.

Deliverables:

- Communication entity.
- Email/SMS provider seams.
- Phone call logging UI/API.

Acceptance criteria:

- Communications can be tied to constituent and optionally application.
- Delivery status is tracked where provider supports it.

### Issue 17: Implement simple segmented campaigns

Goal: support saved audience, message, schedule, and metrics.

Deliverables:

- Campaign entity.
- Audience segment criteria.
- Message template.
- Basic metrics.

Acceptance criteria:

- Segments can use application type/status, term, program, event attendance, lead score, recruiter, and territory.
- Campaign sends are audited.

### Issue 18: Implement admissions events core

Goal: support registration, capacity, check-in, confirmations, and reminders.

Deliverables:

- Event and EventRegistration models.
- Public event registration.
- Capacity enforcement.
- Check-in action.
- Reminder/confirmation notification hooks.

Acceptance criteria:

- Anonymous users can register for public events.
- Existing constituents are matched where possible.
- Capacity is enforced.

### Issue 19: Implement in-app notifications

Goal: notify staff/applicants of workflow events.

Deliverables:

- Notification entity.
- Notification list/read state.
- Event hooks for key applicant and staff events.

Acceptance criteria:

- Applicants receive application/document/decision notifications.
- Staff receive assignment/review/sync/duplicate notifications.

## Milestone 5: SIS Sync, Reporting, RAG, and Operations

### Issue 20: Implement SIS sync framework

Goal: support batch-canonical sync with selected real-time events.

Deliverables:

- SyncJob and SyncEvent models.
- Batch reconciliation runner.
- Real-time event enqueueing.
- Failure/retry visibility.

Acceptance criteria:

- Batch sync pulls terms, programs, existing person matches, and enrollment status.
- Real-time sync pushes application submission, decision, document status, and enrollment intent.
- Failures are visible and audited.

### Issue 21: Implement dashboards and exports

Goal: expose operational dashboards and role-based CSV/PDF exports.

Deliverables:

- Admissions funnel dashboard.
- Document bottleneck dashboard.
- Campaign/event analytics.
- Export service with audit logging.

Acceptance criteria:

- Exports respect permission and program/department scope.
- PII-heavy exports require elevated permission.
- Large exports run async.

### Issue 22: Implement custom fields for Constituent and Application

Goal: support user-defined fields without making core fields dynamic.

Deliverables:

- CustomFieldDefinition and CustomFieldValue models.
- Validation by data type/options.
- Search/report/import/export integration.

Acceptance criteria:

- Custom fields work for constituent and application only.
- Core identity/status/program/term/decision fields remain first-class.

### Issue 23: Implement CSV and Excel imports

Goal: import operational admissions data safely.

Deliverables:

- Upload/import endpoint.
- Field mapping.
- Validation preview.
- Duplicate detection before commit.
- Error report.

Acceptance criteria:

- CSV and Excel are supported.
- Imports create an audit log and import batch record.
- Invalid rows can be downloaded for correction.

### Issue 24: Phase RAG for admissions

Goal: release RAG safely across staff knowledge, applicant status, and reviewer assistant phases.

Deliverables:

- Separate collections for staff knowledge and applicant documents.
- Per-application authorization checks.
- Query audit log.
- Source citation display.

Acceptance criteria:

- Phase 1 supports staff-only policy Q&A.
- Phase 2 supports applicant checklist/status assistant only within own application scope.
- Phase 3 supports reviewer assistant only for assigned applications.

### Issue 25: Add accessibility and recovery gates

Goal: make Section 508/WCAG AA and recovery targets testable release criteria.

Deliverables:

- Accessibility checklist for admissions screens.
- Keyboard/focus acceptance tests where possible.
- Backup and restore runbook.
- Recovery drill checklist.

Acceptance criteria:

- Admissions portal and staff review flows meet WCAG AA basics.
- Restore process documents RPO <= 24 hours and RTO <= 8 business hours.
