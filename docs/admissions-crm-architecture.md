# Admissions CRM v1 Architecture and Domain Model

## Architectural Intent

Admissions CRM v1 is a bounded context focused on enrollment and admissions workflows. It should not be modeled as a generic school CRM or as a broad school operations module. The core language is `Constituent`, `Inquiry`, `Application`, `Checklist`, `Document`, `Decision`, `Campaign`, `Event`, and `SIS Sync`.

## Bounded Contexts

Rendered diagram: [`docs/diagrams/admissions-crm-bounded-contexts.svg`](diagrams/admissions-crm-bounded-contexts.svg)  
Editable diagram: [`docs/diagrams/admissions-crm-bounded-contexts.excalidraw`](diagrams/admissions-crm-bounded-contexts.excalidraw)

```mermaid
flowchart LR
  subgraph Identity[Identity and Access]
    User[Identity/User]
    GlobalRoles[Global school roles]
    AdmissionsRoles[Admissions CRM roles]
  end

  subgraph Admissions[Admissions CRM]
    Constituent[Constituent]
    Inquiry[Inquiry]
    Application[Program Application]
    Checklist[Checklist]
    Document[Admissions Document]
    Decision[Decision]
    DuplicateReview[Duplicate Review]
  end

  subgraph Engagement[Engagement]
    Communication[Communication]
    Campaign[Segmented Campaign]
    Event[Event]
    Notification[Notification]
  end

  subgraph Reporting[Reporting and Audit]
    Dashboard[Dashboards]
    Export[CSV/PDF Export]
    Audit[Audit Log]
  end

  subgraph Integrations[External Integrations]
    SIS[PeopleSoft/SIS]
    Email[Email Provider]
    SMS[SMS Provider]
    Payment[Payment Provider]
    RAG[RAG Service]
  end

  User --> AdmissionsRoles
  AdmissionsRoles --> Application
  Constituent --> Inquiry
  Constituent --> Application
  Application --> Checklist
  Checklist --> Document
  Application --> Decision
  Application --> DuplicateReview
  Application --> Communication
  Campaign --> Communication
  Event --> Communication
  Admissions --> Dashboard
  Admissions --> Audit
  Reporting --> Export
  SIS <--> Admissions
  Communication --> Email
  Communication --> SMS
  Application --> Payment
  Document --> RAG
```

## Core Aggregate Boundaries

### Constituent Aggregate

Owns durable person identity and lifecycle state.

Fields:

- Legal name.
- Preferred name.
- Date of birth.
- Emails.
- Phones.
- Addresses.
- Lifecycle stage.
- Duplicate status.
- External SIS ID.

Invariants:

- A constituent must have first name, last name, date of birth, primary email, and primary phone for admissions flows.
- A constituent may have multiple applications over time.
- A constituent can be staff and applicant through shared identity links.

### Application Aggregate

Owns one program application for one term/program/application type.

Fields:

- Constituent ID.
- Academic term ID.
- Program ID.
- Application type.
- Status state.
- Assigned reviewer.
- Decision.
- Submitted timestamp.

Invariant:

- Only one active application per constituent, academic term, and program.

### Checklist and Document Aggregate

Owns admissions document collection and verification.

Document metadata is owned by CRM. File bytes are stored in object/file storage. Official archive/sync status is tracked for SIS.

### Engagement Aggregates

- `Communication` records email, SMS, and phone-call interactions.
- `Campaign` owns simple segmented campaign configuration and metrics.
- `Event` owns registration, capacity, reminders, and check-in.
- `Notification` owns in-app notification delivery.

### Analytics Integration

Google Analytics 4 (GA4) is an external analytics integration, not the admissions system of record. CRM remains authoritative for constituents, inquiries, applications, campaign audience snapshots, send audit, and provider delivery metrics. GA4 contributes web and campaign behavior metrics for reporting: source/medium/campaign attribution, page and form engagement, public event registration funnel steps, application funnel events, and aggregate conversion reporting.

The analytics boundary is privacy-first:

- No applicant PII is sent to GA4. Names, emails, phone numbers, dates of birth, documents, reviewer notes, and raw application identifiers are excluded from event parameters and URLs.
- GA4 `user_id`, when used, is an opaque CRM identifier suitable for analytics joins, not an email address or human-readable ID.
- Consent Mode v2 or an equivalent consent layer controls analytics and ads storage before GA4 tags or Measurement Protocol events are sent.
- CRM audit logs record staff actions and report/export access. GA4 event collection is separate telemetry and is not a substitute for CRM audit trails.

Recommended reporting flow:

```mermaid
flowchart LR
  Browser[Applicant/public pages] --> GTM[Google Tag / server-side GTM]
  CRM[CRM backend events] --> MP[GA4 Measurement Protocol]
  GTM --> GA4[GA4 property]
  MP --> GA4
  GA4 --> DataAPI[GA4 Data API]
  GA4 --> BigQuery[GA4 BigQuery export]
  CRMDB[(CRM database)] --> Warehouse[Reporting warehouse]
  BigQuery --> Warehouse
  DataAPI --> Reports[Embedded CRM reports]
  Warehouse --> Reports
```

Use the GA4 Data API for aggregate dashboard widgets and near-term campaign reports. Use BigQuery export for CRM joins, history beyond standard API windows, funnel reconstruction, and attribution analysis that combines GA4 events with CRM application state.

### Integration Aggregate

`SyncJob` and `SyncEvent` track batch and real-time PeopleSoft/SIS synchronization.

## Domain Entity Relationship

Rendered diagram: [`docs/diagrams/admissions-crm-domain-entities.svg`](diagrams/admissions-crm-domain-entities.svg)

```mermaid
erDiagram
  IDENTITY_USER ||--o| STAFF_PROFILE : has
  IDENTITY_USER ||--o| APPLICANT_PROFILE : has
  APPLICANT_PROFILE }o--|| CONSTITUENT : maps_to
  CONSTITUENT ||--o{ APPLICATION : submits
  CONSTITUENT ||--o{ INQUIRY : creates
  CONSTITUENT ||--o{ EVENT_REGISTRATION : registers
  CONSTITUENT ||--o{ COMMUNICATION : receives
  APPLICATION }o--|| ACADEMIC_TERM : targets
  APPLICATION }o--|| PROGRAM : targets
  APPLICATION ||--o{ CHECKLIST_ITEM : requires
  CHECKLIST_ITEM ||--o{ DOCUMENT : satisfies
  APPLICATION ||--o{ APPLICATION_TRANSITION : changes_state
  APPLICATION ||--o| APPLICATION_FEE : requires
  APPLICATION ||--o| DECISION : receives
  APPLICATION ||--o{ SIS_SYNC_EVENT : syncs
  CAMPAIGN ||--o{ COMMUNICATION : sends
  EVENT ||--o{ EVENT_REGISTRATION : has
  DUPLICATE_REVIEW }o--|| CONSTITUENT : candidate_a
  DUPLICATE_REVIEW }o--|| CONSTITUENT : candidate_b

  IDENTITY_USER {
    uuid id PK
    string email
    string identity_provider
    string external_subject
  }

  CONSTITUENT {
    uuid id PK
    string first_name
    string last_name
    date date_of_birth
    string lifecycle_stage
    string external_sis_id
  }

  APPLICATION {
    uuid id PK
    uuid constituent_id FK
    uuid academic_term_id FK
    uuid program_id FK
    string application_type
    string status_state
  }

  DOCUMENT {
    uuid id PK
    uuid application_id FK
    string file_key
    string status
    string sis_sync_status
  }
```

## Application State Machine

Rendered diagram: [`docs/diagrams/admissions-crm-application-state-machine.svg`](diagrams/admissions-crm-application-state-machine.svg)  
Editable workflow diagram: [`docs/diagrams/admissions-crm-application-workflow.excalidraw`](diagrams/admissions-crm-application-workflow.excalidraw)

```mermaid
stateDiagram-v2
  [*] --> DRAFT
  DRAFT --> SUBMITTED: applicant submits
  SUBMITTED --> AWAITING_DOCUMENTS: checklist incomplete
  SUBMITTED --> READY_FOR_REVIEW: checklist complete
  AWAITING_DOCUMENTS --> READY_FOR_REVIEW: required docs accepted
  READY_FOR_REVIEW --> IN_REVIEW: reviewer starts
  IN_REVIEW --> DECISION_PENDING: review complete
  DECISION_PENDING --> ADMITTED: admit
  DECISION_PENDING --> DENIED: deny
  DECISION_PENDING --> WAITLISTED: waitlist
  DECISION_PENDING --> DEFERRED: defer
  ADMITTED --> ENROLLED: enrollment confirmed
  DRAFT --> WITHDRAWN: applicant withdraws
  SUBMITTED --> WITHDRAWN: applicant withdraws
  AWAITING_DOCUMENTS --> WITHDRAWN: applicant withdraws
  READY_FOR_REVIEW --> WITHDRAWN: applicant withdraws
  IN_REVIEW --> WITHDRAWN: applicant withdraws
  WAITLISTED --> ADMITTED: offer seat
  DEFERRED --> SUBMITTED: reactivate next term
  DENIED --> [*]
  WITHDRAWN --> [*]
  ENROLLED --> [*]
```

## SIS Sync Design

Rendered diagram: [`docs/diagrams/admissions-crm-sis-sync.svg`](diagrams/admissions-crm-sis-sync.svg)

```mermaid
sequenceDiagram
  participant CRM as Admissions CRM
  participant Queue as Sync Queue
  participant SIS as PeopleSoft/SIS
  participant Audit as Audit Log

  CRM->>Queue: enqueue real-time sync event
  Queue->>SIS: push approved field set
  alt success
    SIS-->>Queue: accepted with external ID/status
    Queue->>CRM: update sync status
    CRM->>Audit: record sync success
  else failure
    SIS-->>Queue: error
    Queue->>CRM: mark failed and retryable
    CRM->>Audit: record sync failure
  end
  CRM->>Queue: nightly batch reconciliation
  Queue->>SIS: pull terms/programs/enrollment status
  SIS-->>CRM: canonical reconciliation data
```

## RAG Security Boundary

Rendered diagram: [`docs/diagrams/admissions-crm-rag-security.svg`](diagrams/admissions-crm-rag-security.svg)

```mermaid
flowchart TD
  User[Authenticated user] --> Auth[Auth service]
  Auth --> Roles[Admissions role and scope]
  Roles --> Router{Document collection}
  Router --> StaffKB[Staff knowledge documents]
  Router --> ApplicantDocs[Applicant documents]
  StaffKB --> StaffAnswer[Policy Q&A]
  ApplicantDocs --> ScopeCheck[Application assignment/scope check]
  ScopeCheck --> ApplicantAnswer[Checklist/reviewer answer]
  StaffAnswer --> Audit[Audit query]
  ApplicantAnswer --> Audit
```

Rules:

- Staff knowledge documents are role/department scoped.
- Applicant documents are application/constituent scoped.
- Applicants can query only their own application context.
- Reviewers can query only assigned applications.
- Global applicant-document search requires elevated permission and audit logging.

## Implementation Guardrails

- Keep admissions permissions action-based, not menu-based.
- Keep core identity/status/program/term/decision fields as first-class columns.
- Use custom fields only for `Constituent` and `Application` in v1.
- Batch sync must reconcile failed real-time sync events.
- Every export and applicant-document RAG query must be audited.
- Section 508 and WCAG AA are non-negotiable UI requirements.
