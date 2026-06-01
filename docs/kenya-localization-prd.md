# PRD: Kenya Localization of Admissions CRM Data Model and Integrations

## Status

This PRD is the active localization layer for the Admissions CRM when deployed for Kenyan institutions. It refines the broader Admissions CRM v1 PRD and replaces US-style assumptions such as PeopleSoft-centric SIS sync, freshman/transfer/graduate application types, freeform addresses, and email-first communications.

Implementation is tracked under Beads issue `schoolCRM-r9k`.

Completed implementation slices:

- `schoolCRM-r9k.1` — Kenya identifier and address value objects.

Open implementation slices:

- `schoolCRM-r9k.2` — Kenya reference data catalog.
- `schoolCRM-r9k.3` — KCSE results and cluster weighting.
- `schoolCRM-r9k.4` — Kenya external adapter ports.
- `schoolCRM-r9k.5` — Kenya identity migration.
- `schoolCRM-r9k.6` — Kenya admissions application localization.
- `schoolCRM-r9k.7` — Kenya integration job tracking.
- `schoolCRM-r9k.8` — Kenya communications preferences.
- `schoolCRM-r9k.9` — Kenya localization documentation.

RAG-related admissions features remain out of scope for the current Kenya implementation pass.

## Problem Statement

The original Admissions CRM v1 PRD was derived from a US-style CRM RFP. It assumed:

- PeopleSoft as the monolithic Student Information System.
- US-style application categories such as freshman, transfer, and graduate.
- Freeform address fields.
- Email-first applicant communications.
- Generic payment and programme catalog integrations.

Those assumptions do not fit a Kenyan admissions process. A Kenya-ready CRM must model:

- KUCCPS placement and programme codes.
- IPRS / Maisha Namba national identity verification.
- NEMIS/KEMIS UPI identifiers for learners.
- KNEC KCSE index numbers and verified examination results.
- County → Sub-County → Ward address hierarchy.
- M-Pesa Daraja application fee payments.
- SMS-first and WhatsApp-supported communications.
- CUE, KUCCPS, and KNQF-backed programme catalog data.

Without this localization, the CRM cannot validate Kenyan identifiers, verify KCSE results, compute KUCCPS cluster eligibility, reconcile M-Pesa fees, segment applicants by Kenyan administrative geography, or communicate reliably through channels applicants actually use.

## Goals

- Make Kenyan identifiers first-class domain concepts.
- Normalize applicant geography against Kenyan administrative boundaries.
- Support official programme and qualification catalogues.
- Verify academic eligibility from KCSE results and KUCCPS cluster rules.
- Replace monolithic SIS sync with explicit Kenya-specific integration adapters.
- Prioritize SMS and WhatsApp while preserving email as a tertiary channel.
- Keep pure domain logic testable without live external services.

## Non-goals

- Direct IPRS whitelisting in v1; use an authorized aggregator.
- TVETA programme ingestion in v1.
- KRA PIN verification, although the data field may be modeled later.
- eCitizen SSO.
- Bank gateway integrations beyond M-Pesa.
- KEMIS-specific integration until a stable public integration surface exists.
- KNEC official REST integration until a stable official API exists.
- Geocoding beyond storing optional coordinates or plus codes.
- Migration of already-created US-style application data unless a separate migration PRD is approved.

## User Stories

1. As a Kenyan applicant, I want to register with my national ID number so that my identity can be verified against IPRS without uploading ID scans.
2. As a Kenyan applicant, I want to optionally provide my NEMIS/KEMIS UPI so that my prior learner records can be linked to my application.
3. As a Kenyan applicant, I want to enter my KCSE index number and examination year so that the institution can verify my results.
4. As a Kenyan applicant, I want the system to calculate cluster eligibility for my chosen programme so that I know whether I meet the cutoff.
5. As a Kenyan applicant, I want to select programmes from an official KUCCPS-coded catalogue so that my application is unambiguous.
6. As a Kenyan applicant, I want to pay application fees through M-Pesa STK Push.
7. As a Kenyan applicant, I want application updates by SMS first and WhatsApp when I opt in.
8. As an admissions officer, I want to find applicants by national ID, UPI, or KCSE index number.
9. As an admissions officer, I want verification status tracked independently for every identifier.
10. As an admissions officer, I want to filter applicants by County, Sub-County, and Ward.
11. As a finance officer, I want every M-Pesa transaction tied to a unique application reference.
12. As a programme administrator, I want annual refreshes of KUCCPS, CUE, and KNQF catalogue data.
13. As a reviewer, I want verified KCSE results visible on the application review screen.
14. As an operator, I want each external integration monitored and retried independently.

## Domain Modules

### 1. Reference Data Catalog

The Reference Data Catalog is a read-mostly module containing:

- Counties.
- Sub-counties.
- Wards.
- Universities.
- Programmes.
- KUCCPS clusters and cutoffs.
- KNQF levels.

Writes occur only through migrations or annual refresh jobs. Application traffic reads from the catalog but must not mutate it.

Canonical sources are tracked in [`docs/research/kenya-data-sources.md`](research/kenya-data-sources.md).

### 2. Kenyan Address

The Kenyan address model contains:

- Postal block: PO Box, 5-digit postal code, postal town.
- Physical block: county code/name, sub-county, ward, estate or locality, street, building, house number.
- Geo block: optional latitude, longitude, and plus code.

At least one address block must be populated. County/sub-county/ward combinations must resolve through the Reference Data Catalog once that catalog exists.

### 3. Identifier Value Objects

Identifier value objects validate format only:

- National ID.
- UPI.
- KCSE index number.

Verification status belongs to the Constituent aggregate, not the value object.

### 4. KCSE Result and Cluster Weighting

`KCSEResult` stores:

- KCSE index number.
- Examination year.
- Subject codes.
- Subject grades.
- Mean grade.

Cluster weighting is deterministic pure domain logic calculated against cluster definitions from the Reference Data Catalog. Verification of results belongs to the KNEC adapter.

## External Adapters

The Kenya implementation replaces monolithic SIS sync with explicit adapters:

- `KuccpsImporter` — imports programmes, clusters, cutoffs, and placement-cycle data.
- `KnecVerifier` — verifies KCSE result slips or portal/QR outputs.
- `IprsLookup` — verifies national ID data through an authorized aggregator.
- `MpesaDarajaGateway` — initiates STK Push and handles payment callbacks.
- `AfricasTalkingSmsGateway` — sends SMS and records delivery reports.
- `WhatsAppCloudGateway` — sends approved WhatsApp templates and session messages.

All adapters should share common retry, circuit-breaker, configuration, observability, and fixture-based contract-test infrastructure. CI must not call live vendor APIs.

## Existing Aggregate Changes

### Constituent

Replace `external_sis_id` as the primary external identity with three nullable, independently indexed identifiers:

- National ID plus verification status, verified timestamp, and adapter receipt.
- UPI plus verification status, verified timestamp, and adapter receipt.
- KCSE index number plus verification status, verified timestamp, and adapter receipt.

Keep `external_sis_id` as a deprecated field for one release, backfill where unambiguous, and flag ambiguous values for manual review.

### Application

Use Kenya-appropriate application types:

- `KUCCPS_PLACEMENT`
- `SELF_SPONSORED_UNDERGRAD`
- `DIPLOMA`
- `MASTERS`
- `PHD`
- `TVET`
- `BRIDGING`
- `CERTIFICATE`

Applications may include KUCCPS placement data and KCSE results. Programme selection should reference the Reference Data Catalog, not a freeform string.

### Integration Jobs

Split monolithic sync job tracking into per-adapter job records:

- KUCCPS import jobs.
- KNEC verification jobs.
- IPRS lookup jobs.
- M-Pesa payment jobs.
- SMS delivery jobs.
- WhatsApp delivery jobs.

Each job type owns its own retry policy, external reference, failure detail, and observability surface.

### Notification Preferences

Kenya defaults to:

1. SMS.
2. WhatsApp.
3. Email.

Existing email-only preferences should migrate safely, with SMS enabled by default where a phone number exists and WhatsApp disabled until explicit opt-in.

## Testing Strategy

- Value-object tests must be table-driven and verify external behavior only.
- KCSE and cluster-weight tests must cover grade boundaries, missing subjects, and deterministic cluster calculations.
- Catalog tests should verify seed integrity and lookup behavior against canonical counts where stable.
- Adapter tests must use recorded fixtures or stubs, never live vendor calls in CI.
- Migration tests must prove no identity data is silently lost.

## Acceptance Criteria

- All child Beads under `schoolCRM-r9k` are closed.
- Kenya PRD and research source docs are committed under `docs/`.
- Existing Admissions CRM docs link to this PRD as the active Kenya localization layer.
- Backend domain, app, store, migration, and API layers expose the Kenya model.
- Web-admin and applicant-facing screens can display and collect Kenya-specific fields, using static data first where API integration is deferred.
- External adapter implementations are independently configurable and testable.
- Validation gates pass for every implementation slice.
