# Kenya Admissions Localization Data Sources

Research date: 2026-06-01  
Tracking issue: `schoolCRM-r9k.9`

This reference lists the current source systems and provider documentation for the Kenya Admissions CRM localization. It supports [`docs/kenya-localization-prd.md`](../kenya-localization-prd.md) and should be refreshed before any importer hardcodes source URLs or dataset versions.

## Summary

```
┌───────────────────────────────────────┬──────────────────────┬──────────────────────────────┐
│ Source                                │ Type                 │ Public API                   │
├───────────────────────────────────────┼──────────────────────┼──────────────────────────────┤
│ HDX COD-AB Kenya boundaries           │ Static files         │ No — SHP/GeoJSON/XLSX        │
│ CUE approved programmes               │ PDF + portal         │ No — dated PDF snapshots     │
│ KUCCPS programme and cluster lists    │ Auth portal files    │ No — Excel/PDF per cycle     │
│ KNQA / KNQF levels                    │ PDF + portal         │ No — publications only       │
│ Safaricom M-Pesa Daraja               │ REST API             │ Yes                          │
│ Celcom Africa SMS                     │ REST API/SMPP        │ Yes                          │
│ WhatsApp Cloud API                    │ REST API             │ Yes                          │
│ Smile ID IPRS aggregation             │ REST API             │ Yes, commercial agreement    │
│ MetaMap IPRS aggregation              │ REST API             │ Yes, commercial agreement    │
│ Didit identity verification           │ REST API             │ Yes, Kenya IPRS unconfirmed  │
└───────────────────────────────────────┴──────────────────────┴──────────────────────────────┘
```

## Administrative Boundaries

### HDX COD-AB Kenya

- Dataset page: <https://data.humdata.org/dataset/cod-ab-ken>
- Formats: GeoJSON, SHP/geodatabase, XLSX.
- Levels: Admin 0 country, Admin 1 counties, Admin 2 sub-counties.
- Current notes from research: COD-AB v01; boundaries created 2018-07-03; IMWG-endorsed in October 2016.
- License: CC BY-IGO.

Implementation notes:

- Treat this as a static file import, not an API integration.
- Re-check the dataset page before pinning a file URL; HDX download URLs can change.
- HDX is an operational humanitarian layer, not necessarily the government-authoritative Survey of Kenya layer.
- If legal/government reporting requires authoritative boundaries, verify against KNBS or Survey of Kenya before launch.

## University and Programme Catalogues

### CUE — Commission for University Education

- Official site: <https://www.cue.or.ke>
- Accredited universities portal: <https://imis.cue.or.ke/RecognitionAndEquationforQualifications/AccreditedUniversities>
- Approved programmes PDF, December 2022: <https://www.cue.or.ke/documents/2023/Approved_Academic_Programmes_Offered_Universities_in_Kenya_December_2022.pdf>
- Approved programmes PDF, July 2021: <https://www.cue.or.ke/documents/Approved_Academic_Programmes_July2021.pdf>
- Accredited universities PDF, August 2022: <https://www.cue.or.ke/documents/Accredited_Universities_Kenya_August_2022.pdf>

Implementation notes:

- CUE publishes dated PDF snapshots and portal pages, not a public JSON API.
- Importers should preserve the source date and document URL used for each catalogue version.
- A newer stable PDF URL may appear under `cue.or.ke/documents/`; verify before each annual refresh.

### KUCCPS — Programmes, Clusters, and Placement Cycles

- Student portal: <https://students.kuccps.net>
- Degree application cycle page: <https://kuccps.net/degree-application-20252026-placement-cycle>
- Official domain: <https://www.kuccps.ac.ke> redirects to `kuccps.net`.

Implementation notes:

- No public REST API was found.
- Programme lists and cluster documents are normally downloaded from the authenticated student portal under the "Programme Lists" tab.
- Treat KUCCPS files as placement-cycle-specific inputs, such as `2025-2026`.
- Do not hardcode cluster subjects or cutoffs; they can change by cycle.
- Import jobs should be idempotent and should record source file hash, placement cycle, import timestamp, and row counts.

### KNQA / KNQF Levels

- Official authority site: <https://knqa.go.ke>
- KNQF framework page: <https://knqa.go.ke/qualifications-framework>
- Online portal: <https://qa.knqa.go.ke>
- Reference publication: <https://knqa.go.ke/wp-content/uploads/2024/07/The-Harmonizer-Vol.-2.pdf>

Implementation notes:

- KNQF is a 10-level qualifications framework under the KNQF Act 2014.
- KNQF levels are stable enough to seed as reference data, but descriptors should cite the source publication version.
- KNQA provides publications and portals, not a public API for level descriptors.

## Payments

### Safaricom M-Pesa Daraja

- Developer portal: <https://developer.safaricom.co.ke>
- API catalogue: <https://developer.safaricom.co.ke/apis/>

Relevant APIs:

- STK Push / Lipa Na M-PESA Online.
- C2B.
- Transaction Status.
- Reversal.
- Account Balance.

Implementation notes:

- Target Daraja 3.0/current developer portal behavior, not legacy Daraja 2.x assumptions.
- Sandbox and production credentials are separate.
- Production onboarding requires Kenyan business documentation, paybill/till setup, banking, and KRA documentation.
- Every STK push must use a unique application reference and persist callback state for reconciliation.

## Communications

### Celcom Africa SMS

- Website: <https://celcomafrica.com>
- Product surface: Bulk SMS, SMS Gateway/API, SMPP, Sender IDs, and Delivery Reports.

Implementation notes:

- Account registration is required for credentials.
- Sender IDs and shortcodes require registration, Communications Authority compliance, and lead time.
- Delivery callbacks should be treated as first-class integration events.
- Keep SMS credits/billing out of the admissions domain port; reconcile costs through provider reports.

### WhatsApp Cloud API

- Cloud API docs: <https://developers.facebook.com/docs/whatsapp/cloud-api>
- Webhooks: <https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks>
- Message templates: <https://developers.facebook.com/docs/whatsapp/message-templates>
- Pricing: <https://developers.facebook.com/docs/whatsapp/pricing>

Implementation notes:

- Use Cloud API only; the on-premises API is deprecated.
- A verified Meta Business account and WhatsApp Business phone number are required.
- Outbound-initiated messages require pre-approved templates.
- Store template identifiers and approval status in configuration, not in code.

## Identity Verification

There is no general public IPRS REST API suitable for direct integration. The v1 integration should use an authorized aggregator and should log every verification attempt with timestamp, provider, request reference, and result.

### Smile ID

- Developer docs: <https://docs.usesmileid.com>
- Kenya coverage page: <https://usesmileid.com/countries/kenya>

Implementation notes:

- Smile ID documents Kenya national ID checks against IPRS, passport, KRA PIN, biometric KYC, document verification, and AML/sanctions screening.
- Authentication uses partner ID and API key from the Smile ID portal.
- Production requires a commercial agreement.
- Lowest-risk initial choice based on documented Kenya coverage and African market focus.

### MetaMap

- API guide: <https://docs.metamap.com/docs/api-guide>
- Kenya IPRS GovCheck reference: <https://docs.metamap.com/reference/govchecks-kenya-iprs>
- Kenya BRS GovCheck reference: <https://docs.metamap.com/reference/govchecks-kenya-brs>

Implementation notes:

- Kenya IPRS checks validate national ID and name data through a webhook-based flow.
- The Kenya-specific docs appeared older at research time; verify current API version before implementation.
- Production requires a commercial agreement.

### Didit

- Developer docs: <https://docs.didit.me>
- LLM-friendly documentation index: <https://docs.didit.me/llms.txt>

Implementation notes:

- Public docs describe broad KYC/KYB coverage, but Kenya IPRS support was not confirmed in public pages at research time.
- Treat Didit as unconfirmed for Kenya IPRS until vendor confirmation.
- If evaluated, require written confirmation of Kenya national ID/IPRS support before choosing it for v1.

## Importer Design Guidance

- Keep raw downloaded files under controlled fixtures or object storage with checksums.
- Store source URL, source publication date, import timestamp, row counts, and file hash for every catalog refresh.
- Prefer manual-reviewed annual refresh jobs over unattended scraping for CUE and KUCCPS.
- Keep parsing code deterministic and fixture-tested.
- Never call live vendor APIs in CI; use recorded fixtures or provider sandboxes only in explicitly marked integration environments.
