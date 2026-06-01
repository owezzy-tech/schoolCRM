# Admissions CRM v1 Wireframes

These wireframes translate the Admissions CRM v1 PRD into UI screens that follow the current `web-admin` frontend design language. Every workflow in the PRD is captured as a distinct screen or sub-screen.

Primary artifacts:

- Screen PNG set: [`docs/diagrams/admissions-crm-screens-png/`](diagrams/admissions-crm-screens-png/)
- User journey diagram: [`docs/diagrams/admissions-crm-user-journey.svg`](diagrams/admissions-crm-user-journey.svg)

Use the PNG screen set as the implementation reference for frontend work. Each PNG is embedded in the matching screen section below using relative Markdown image paths so the previews render in GitHub and other Markdown viewers. The SVG user journey remains the cross-screen workflow reference.

## Existing Frontend Design Cues Used

Source patterns from the Angular admin frontend:

- Layout shell: fixed header/sidebar + `.content > .content-block` from `api/frontends/web-admin/src/app/layout/app-layout/main-layout/main-layout.component.html`.
- Breadcrumb: `app-breadcrumb` pattern from `api/frontends/web-admin/src/app/shared/components/breadcrumb/breadcrumb.component.html`.
- Dashboard cards: `.card.card-statistic-2` from `api/frontends/web-admin/src/app/admin/dashboard/main/main.component.html` and `src/assets/scss/ui/_card.scss`.
- Tables: `.table.table-hover`, status badges, and action buttons from `src/assets/scss/plugins/_tables.scss`.
- Forms: Material outline fields and responsive cards from `src/assets/scss/components/_formcomponents.scss`.
- RAG hero/card pattern: purple-to-cyan header and two-column grid from `api/frontends/web-admin/src/app/admin/rag/chat/rag-chat.component.html` and `.scss`.

---

## Screen Set — Staff Admin Screens

### 1. Admissions Dashboard

Route: `/admin/admissions/dashboard`  
Roles: `ADMISSIONS_ADMIN`, `RECRUITER`, `REPORT_VIEWER`

Purpose: admissions leadership and recruiters see funnel health at a glance.

Rendered screen:

![Admissions dashboard wireframe](diagrams/admissions-crm-screens-png/01-dashboard.png)

Key components:

- Breadcrumb: `Admissions / Dashboard`.
- Purple-to-cyan hero banner: "Admissions Command Center".
- Four statistic cards: New Inquiries (count, trend), Submitted Applications, Ready for Review, Admitted This Term.
- Funnel chart card (Prospects → Inquiries → Applicants → Admitted → Enrolled).
- Recruiter workload card (bar chart by recruiter).
- Recent applications table (5 rows) with status badges and row actions (View, Assign).
- Quick-action buttons: New Inquiry, New Application, View Reports.

### 2. Constituent List

Route: `/admin/admissions/constituents`  
Roles: `ADMISSIONS_ADMIN`, `RECRUITER`

Purpose: search and manage the prospect/applicant population.

Rendered screen:

![Constituent list wireframe](diagrams/admissions-crm-screens-png/02-constituents.png)

Key components:

- Search bar + advanced filter drawer (lifecycle stage, program interest, term, lead score, territory, assigned recruiter).
- Lifecycle stage chips as quick-filters.
- Table columns: Name, Email, Phone, Stage, Lead Score, Assigned To, Last Activity, Actions.
- Row actions: View, Edit, Merge (if flagged duplicate), Assign.
- Bulk actions toolbar: Assign Recruiter, Add to Campaign, Export.

### 3. Constituent Detail

Route: `/admin/admissions/constituents/:id`  
Roles: `ADMISSIONS_ADMIN`, `RECRUITER`

Purpose: full profile view for one constituent across their lifecycle.

Rendered screen:

![Constituent detail wireframe](diagrams/admissions-crm-screens-png/03-constituent-detail.png)

Key components:

- Identity summary card (name, DOB, email, phone, lifecycle stage badge, SIS ID, duplicate flag).
- Contact info card (emails, phones, addresses).
- Applications tab: list of all applications for this constituent with status badges.
- Communications tab: timeline of all emails, SMS, calls, notifications.
- Events tab: registered/attended events.
- Activity timeline card (audit trail: created, inquiry submitted, application started, etc.).
- Custom fields card (any custom fields defined for Constituent).
- Action panel: Edit, Create Application, Log Communication, Flag Duplicate.

### 4. Duplicate Review Queue

Route: `/admin/admissions/duplicates`  
Roles: `ADMISSIONS_ADMIN`

Purpose: staff review fuzzy duplicate matches before merging.

Rendered screen:

![Duplicate review queue wireframe](diagrams/admissions-crm-screens-png/04-duplicates.png)

Key components:

- Queue count badge in sidebar nav.
- Table columns: Candidate A, Candidate B, Match Score, Match Reason, Status, Actions.
- Row expand: side-by-side field comparison (name, DOB, email, phone, SIS ID).
- Actions per pair: Merge (choose primary), Not Duplicate (dismiss), Defer.
- Filters: match score range, date range.

### 5. Inquiry Management

Route: `/admin/admissions/inquiries`  
Roles: `ADMISSIONS_ADMIN`, `RECRUITER`

Purpose: manage anonymous and authenticated inquiries before they become applications.

Rendered screen:

![Inquiry management wireframe](diagrams/admissions-crm-screens-png/05-inquiries.png)

Key components:

- Filter card: date range, source, program interest, status (New, Contacted, Converted, Closed).
- Table columns: Name/Email, Source, Program Interest, Term, Status, Assigned To, Date, Actions.
- Row actions: View, Convert to Application, Assign, Log Communication.
- Bulk actions: Assign, Add to Campaign.

### 6. Online Applications List

Route: `/admin/admissions/online-applications`  
Roles: `ADMISSIONS_ADMIN`, `RECRUITER`, `APPLICATION_REVIEWER`

Purpose: staff triage and filter program applications.

Rendered screen:

![Online applications list wireframe](diagrams/admissions-crm-screens-png/06-applications.png)

Key components:

- Search and filter card with Material outline inputs (status, type, program, term, recruiter, date range).
- Status chips for quick-filter (Draft, Submitted, Awaiting Docs, Ready for Review, In Review, Decision Pending).
- Table columns: Applicant, Program, Term, Type (Freshman/Transfer/Graduate), Status, Checklist Progress (x/y), Assigned Reviewer, Last Updated, Actions.
- Row actions: View, Assign Reviewer, Change Status.
- Bulk actions: Assign Reviewer, Export.
- Table pagination and configurable page size.

### 7. Application Detail / Review Workspace

Route: `/admin/admissions/applications/:id`  
Roles: `ADMISSIONS_ADMIN`, `APPLICATION_REVIEWER`, `RECRUITER`

Purpose: reviewers inspect one application end-to-end and make decisions.

Rendered screen:

![Application detail review workspace wireframe](diagrams/admissions-crm-screens-png/07-application-detail.png)

Key components:

- **Header card**: Applicant name, program, term, type badge, current status badge, assigned reviewer.
- **Application state stepper/timeline**: visual progress through state machine (Draft → Submitted → Awaiting Docs → Ready for Review → In Review → Decision Pending → Final).
- **Applicant summary tab**: identity fields, contact info, academic background, custom fields.
- **Document checklist tab**: table with Document Name, Required (Y/N), Status (Uploaded/Pending Review/Accepted/Rejected/Waived/Expired), Uploaded Date, Reviewer Notes, Actions (Accept/Reject/Waive/Request Reupload).
- **Communications tab**: messages sent/received for this application.
- **Reviewer notes tab**: internal notes with author + timestamp.
- **Audit history tab**: full state transition log.
- **Decision action panel** (sticky bottom): buttons for Admit, Deny, Waitlist, Defer, Withdraw, Request Documents. Each requires confirmation modal with reason/note field.
- **RAG Reviewer Assistant panel** (collapsible right sidebar): scoped to this applicant's documents for AI-assisted review.

### 8. Communications Center

Route: `/admin/admissions/communications`  
Roles: `ADMISSIONS_ADMIN`, `RECRUITER`, `MARKETING_MANAGER`

Purpose: view all communications and compose new messages.

Rendered screen:

![Communications center wireframe](diagrams/admissions-crm-screens-png/08-communications.png)

Key components:

- Tab bar: All Messages, Email, SMS, Phone Calls, Notifications.
- Filter card: date range, direction (inbound/outbound), constituent, application, campaign.
- Table columns: Recipient/Sender, Channel, Subject/Preview, Status (Sent/Delivered/Failed/Opened), Date, Actions.
- Compose button → modal: select channel, select recipient(s) or segment, template picker, body editor, schedule option.
- Phone call logging: manual entry form (constituent, duration, notes, outcome).

### 9. Campaign Builder

Route: `/admin/admissions/campaigns`  
Roles: `MARKETING_MANAGER`, `ADMISSIONS_ADMIN`

Purpose: create and manage segmented outreach campaigns.

Rendered screen:

![Campaign builder wireframe](diagrams/admissions-crm-screens-png/09-campaigns.png)

Key components:

- Campaign list table: Name, Status (Draft/Active/Paused/Completed), Channel, Audience Size, Sent, Opened, Dates, Actions.
- Create/Edit campaign flow (stepper):
  - Step 1: Name, channel (email/SMS), schedule.
  - Step 2: Audience segment builder (filters: application type, status, term, program, event attendance, lead score, recruiter, geography).
  - Step 3: Template selection or compose message.
  - Step 4: Review and activate.
- Campaign detail: metrics cards (sent, delivered, opened, clicked, bounced), GA4-attributed landing page/session conversions where configured, recipient list, and send audit trail.

Implementation note: campaign delivery metrics come from email/SMS providers and CRM send audit. GA4 metrics are displayed separately as aggregate attribution and funnel widgets; applicant PII must not appear in GA event parameters or report tables.

### 10. Event Management

Route: `/admin/admissions/events`  
Roles: `EVENT_MANAGER`, `ADMISSIONS_ADMIN`

Purpose: manage admissions events (open houses, info sessions, tours).

Rendered screen:

![Event management wireframe](diagrams/admissions-crm-screens-png/10-events.png)

Key components:

- Event list table: Name, Type, Date, Capacity, Registered, Checked-In, Status (Upcoming/In Progress/Completed/Cancelled), Actions.
- Create/Edit event form: name, type, description, date/time, location (virtual/physical), capacity, registration deadline, auto-confirmation toggle, auto-reminder toggle.
- Event detail view:
  - Summary card: event info, capacity meter (registered/capacity).
  - Registrations tab: table of registrants with check-in status and actions (Check In, Cancel, Send Reminder).
  - Check-in mode: simplified view for event-day staff with search + one-click check-in.
  - Communications tab: messages sent for this event.

### 11. Lead Management / Pipeline

Route: `/admin/admissions/leads`  
Roles: `ADMISSIONS_ADMIN`, `RECRUITER`

Purpose: view lead pipeline and manage assignment/scoring.

Rendered screen:

![Lead management pipeline wireframe](diagrams/admissions-crm-screens-png/11-leads.png)

Key components:

- Pipeline view toggle: Table view / Kanban board (columns = lead score bands: Cold, Warm, Hot, Ready to Apply).
- Table columns: Constituent, Score Band, Source, Program Interest, Assigned Recruiter, Last Activity, Days Since Contact, Actions.
- Assignment actions: Reassign, Update Score.
- Bulk actions: Reassign, Add to Campaign.

### 12. Lead Assignment Rules (Settings)

Route: `/admin/admissions/settings/lead-rules`  
Roles: `ADMISSIONS_ADMIN`

Purpose: configure automatic lead assignment rules.

Rendered screen:

![Lead assignment rules within settings wireframe](diagrams/admissions-crm-screens-png/14-settings.png)

Note: lead assignment rules are a settings sub-screen, so they are rendered in the Settings & Configuration PNG rather than a separate `12-lead-rules.png` file.

Key components:

- Rules list (ordered, first-match): Priority, Rule Name, Criteria Summary, Assigned Recruiter, Status (Active/Inactive), Actions.
- Create/Edit rule form: name, criteria (territory, program, app type, term, source, event attendance, last-name alpha range), target recruiter, priority order.
- Manual override note.

### 13. Reports Center

Route: `/admin/admissions/reports`  
Roles: `REPORT_VIEWER`, `ADMISSIONS_ADMIN`

Purpose: operational dashboards and exports.

Rendered screen:

![Reports center wireframe](diagrams/admissions-crm-screens-png/12-reports.png)

Key components:

- Report cards grid: Funnel Report, Application Status Report, Recruiter Performance, Campaign Analytics, Event Analytics, GA4 Source Attribution, Document Bottleneck Report, SIS Sync Status.
- Each card: title, description, last generated date, Generate/Download buttons.
- Export history table: Report Name, Format (CSV/PDF), Generated By, Date, Download link.
- Role-based visibility (some reports hidden for non-admin roles).

GA4-backed report widgets show aggregate source/medium/campaign performance, public-page funnel steps, event-registration funnel activity, and inquiry/application conversion trends. CRM-owned report widgets remain authoritative for application status, document bottlenecks, reviewer workload, SIS sync status, and audited exports.

### 14. Audit Log

Route: `/admin/admissions/audit`  
Roles: `ADMISSIONS_ADMIN`, `SUPER_ADMIN`

Purpose: view security-sensitive and workflow audit events.

Rendered screen:

![Audit log wireframe](diagrams/admissions-crm-screens-png/13-audit.png)

Key components:

- Filter card: date range, actor, action type (login, role change, export, application transition, decision, merge, SIS sync, RAG query, custom field change).
- Table columns: Timestamp, Actor, Action, Target Entity, Details/Reason, IP Address.
- Row expand: full detail JSON.
- Export audit log button.

### 15. Settings & Configuration

Route: `/admin/admissions/settings`  
Roles: `ADMISSIONS_ADMIN`

Purpose: hub for admissions configuration.

Rendered screen:

![Settings and configuration wireframe](diagrams/admissions-crm-screens-png/14-settings.png)

Sub-sections (tab bar or cards):

- **Custom Fields**: manage custom field definitions for Constituent and Application. Table: Field Name, Entity, Type, Required, Searchable, Actions.
- **Checklist Templates**: define per-program document checklist requirements. Table: Program, Document Name, Required (Y/N), Description, Actions.
- **Lead Assignment Rules** (links to screen 12).
- **Lead Scoring Rules**: rule list with criteria and score point assignments.
- **SIS Sync Status**: last sync time, sync errors table, manual sync trigger button, approved field mapping display.
- **Import/Export Config**: import templates (CSV/Excel) for constituents, inquiries, applications, event registrations, campaign audiences.

### 16. RAG Staff Assistant

Route: `/admin/admissions/rag-assistant`  
Roles: `ADMISSIONS_ADMIN`, `RECRUITER`, `APPLICATION_REVIEWER`

Purpose: staff-facing policy/knowledge Q&A chatbot.

Rendered screen:

![Admissions policy assistant wireframe](diagrams/admissions-crm-screens-png/15-assistant.png)

Key components:

- Reuses existing RAG Chat pattern (`rag-chat.component`).
- Purple-to-cyan hero: "Admissions Policy Assistant".
- Chat panel with message bubbles, citations panel, source document links.
- Scoped to admissions policy/knowledge documents only.
- Disclaimer banner: "Responses are AI-generated. Verify against official policy."

---

## Screen Set — Applicant Portal Screens

### 17. Applicant Portal — Landing / Account

Route: `/applicant/dashboard`  
Roles: `APPLICANT` (authenticated)

Purpose: applicant sees their applications and next steps.

Rendered screen:

![Applicant portal dashboard wireframe](diagrams/admissions-crm-screens-png/16-portal-home.png)

Key components:

- Simplified header (no admin sidebar): logo, applicant name, logout.
- Welcome card with applicant name.
- Active applications cards: one card per application showing program, term, status badge, checklist progress bar, next action button.
- Start New Application button (if eligible).
- Upcoming events card: registered events with dates.
- Notifications/messages card: recent communications.

### 18. Applicant Portal — Application Form

Route: `/applicant/applications/:id/edit`  
Roles: `APPLICANT`

Purpose: multi-step application form.

Rendered screen:

![Applicant application form wireframe](diagrams/admissions-crm-screens-png/17-portal-apply.png)

Key components:

- Form stepper: Personal Info → Academic History → Program Selection → Documents → Review & Submit.
- Each step is a card with Material outline form fields.
- Save Draft button (always visible).
- Progress indicator showing completed/remaining steps.
- Validation errors shown inline.
- Final Review step: summary of all entered data with Edit links per section.
- Submit button with confirmation dialog.
- Application fee payment (if applicable) before submission.

### 19. Applicant Portal — Application Status

Route: `/applicant/applications/:id/status`  
Roles: `APPLICANT`

Purpose: applicant tracks application progress and uploads documents.

Rendered screen:

![Applicant application status wireframe](diagrams/admissions-crm-screens-png/18-portal-status.png)

Key components:

- Status hero card: current state badge, next expected action, estimated timeline.
- State stepper (visual): showing where they are in the process.
- Document checklist cards: each document shows status (Uploaded/Pending/Accepted/Rejected), upload button, rejection reason (if applicable), reupload action.
- File upload dropzone per document item.
- Decision card (when decision made): Admitted/Denied/Waitlisted with next steps.
- RAG applicant assistant (collapsible panel): scoped to checklist/status questions only.

### 20. Applicant Portal — Event Registration

Route: `/applicant/events`  
Roles: `APPLICANT` (or anonymous for public events)

Purpose: browse and register for admissions events.

Rendered screen:

![Applicant event registration wireframe](diagrams/admissions-crm-screens-png/19-portal-events.png)

Key components:

- Event cards grid: event name, date, type, location, capacity remaining, Register button.
- Registered events section: showing confirmed registrations with calendar add link.
- Registration confirmation modal with auto-email confirmation.

### 21. Public — Inquiry Form

Route: `/inquiry` (anonymous, no auth required)

Purpose: anonymous prospects submit interest.

Rendered screen:

![Public inquiry form wireframe](diagrams/admissions-crm-screens-png/20-portal-inquiry.png)

Key components:

- Clean public-facing page (no admin chrome).
- Form card: First Name, Last Name, Email, Phone, Program of Interest (dropdown), Term of Interest, Source (how did you hear about us), Message/Questions.
- Submit button → confirmation page with "What happens next" info.
- Optional: link to event registration, link to create applicant account.

---

## Workflow Coverage Matrix

| PRD Workflow | Screen(s) |
|---|---|
| Constituent identity & lifecycle | 2, 3 |
| Duplicate detection & merge | 4 |
| Anonymous inquiry | 21 |
| Authenticated application | 18, 19 |
| Application state machine | 6, 7 |
| Document checklist & upload | 7, 19 |
| Application review & decision | 7 |
| Communications (email/SMS/phone) | 8 |
| Segmented campaigns | 9 |
| Events (registration, check-in) | 10, 20 |
| Lead assignment & scoring | 11, 12 |
| Reporting & exports | 13 |
| Audit trail | 14 |
| Custom fields & config | 15 |
| SIS sync monitoring | 15 |
| Import/export management | 15 |
| RAG staff policy Q&A | 16 |
| RAG applicant assistant | 19 |
| RAG reviewer assistant | 7 |
| Applicant portal (status/uploads) | 17, 18, 19 |
| Public event registration | 20 |
| Dashboard / funnel metrics | 1 |

## Implementation Notes

- Keep `Constituent` and `Application` as first-class labels in detail screens.
- Preserve the existing admin template's Roboto typography, white card surfaces, 10px base card radius, and subtle `rgba(183, 192, 206, 0.2)` shadows.
- Use existing status colors consistently:
  - Green: accepted/admitted/complete/enrolled.
  - Orange: awaiting documents/attention needed/warm lead.
  - Cyan/blue: ready/in review/cool.
  - Red: rejected/denied/error/hot lead.
  - Purple: primary admissions/RAG emphasis.
- Design tables mobile-first by stacking filters above data and keeping row actions compact.
- Preserve WCAG AA basics: visible focus states, high-contrast status text, non-color-only status labels, keyboard-accessible actions.
- Applicant portal uses simplified layout (no admin sidebar) but same card/typography system.
- Sidebar nav structure for admissions:
  ```
  ADMISSIONS
  ├── Dashboard
  ├── Constituents
  ├── Inquiries
  ├── Applications
  ├── Duplicates (badge count)
  ├── Communications
  ├── Campaigns
  ├── Events
  ├── Leads
  ├── Reports
  ├── Audit Log
  ├── RAG Assistant
  └── Settings
  ```
- Angular routing module: `/admin/admissions/` prefix for all staff screens; `/applicant/` prefix for portal screens.
- Each screen maps to a lazy-loaded Angular module with its own route children.
