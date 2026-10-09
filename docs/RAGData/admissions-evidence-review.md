# Maseno admissions evidence review

Bead `schoolCRM-o58.5`; review revision `maseno-admissions-evidence-2026-10-09-v1`.

## Source identity and scope

On 9 October 2026, the supplied four-page cutoff PDF was downloaded from [Maseno University's official site](https://www.maseno.ac.ke/sites/default/files/2024-02/CUT-OFF-POINTS-FOR-2023-2024-COHORT.pdf) with ordinary HTTPS certificate verification. Its bytes match the supplied artifact: SHA-256 `418598585c0a344496f547c43f810e3dfe1bb2469963649db612bf6b48bb0843`. The university registrar heading states academic year 2023/2024; its two columns cover 2022/2023 and 2023/2024. The scan bears a `14 SEP 2023` stamp. These are historical published claims, not current admissions policy, subject requirements or cluster calculation rules.

The supplied CSV still has no verified publisher, revision or applicable cycle. Its undated cluster points cannot be assigned to either PDF cohort by numerical coincidence. Explicit programme equivalences remain unapproved until their source authority and identity are established. Source metadata is pinned in [review/source-provenance.json](review/source-provenance.json).

## Cut-off table scan review

Review revision `maseno-cutoff-scan-review-2026-10-09-v1` is recorded in [review/cutoff-table-review.json](review/cutoff-table-review.json). It is pinned to the PDF checksum above and to the OCR evidence checksum. The parent agent inspected all four page renders and the eight originally flagged cells. An independent Claude visual review then compared all 87 rows against the four page renders: order, code, name, school, cluster and both cut-offs. Every code and school assignment matches; rows 71–72 continue under School of Medicine from page 3. The page 3 stamp is clipped, so its day is unreadable; pages 1, 2 and 4 read `14 SEP 2023`. The printed school heading `MATHEMEMATICS, STATITICS` is kept as printed.

This review checks the transcription only. It does not mean Maseno approved our interpretation, and it does not make any cut-off current or executable.

### Resolved corrections

Each correction keeps the original OCR text in the reviewed table and in the graph.

| S.No | Code | Field | Original extraction | Reviewed reading |
| --- | --- | --- | --- | --- |
| 14 | 1229179 | Cluster | Empty | `3` |
| 15 | 1229509 | Cluster | Empty | `3` |
| 42 | 1229539 | Cluster | Empty | `10` |
| 56 | 1229137 | Cluster | Empty | `19` |
| 57 | 1229155 | Cluster | Empty | `19` |
| 14 | 1229179 | Name | `…Theatre Studies, With Il)` | `…Theatre Studies, With IT)` |
| 20 | 1229621 | Name | `Bachelor of Arts (Language and Communication,` | `Bachelor of Arts (Language and Communication, With IT)` |
| 30 | 1229337 | Name | `Bachelor.of Arts (…)` | `Bachelor of Arts (Counseling Psychology, with IT)` |
| 39 | 1229299 | Name | `Barhelor of Rusiness Entrepreneurship, with IT` | `Bachelor of Business Entrepreneurship, with IT` |
| 52 | 1229612 | Name | `…(Information Systems) •` | `Bachelor of Science (Information Systems)` |
| 82 | 1229765 | Name | `Bachelor or Arts (Urban and Regionai Planning, With IT)` | `Bachelor of Arts (Urban and Regional Planning, With IT)` |

### Unresolved: disputed or uncertain

These cells keep their evidence but never produce a numeric candidate. They need a better copy or confirmation from Maseno.

| S.No | Code | Field | Original extraction | Scan observation | Treatment |
| --- | --- | --- | --- | --- | --- |
| 7 | 1229185 | 2023/2024 cut-off | `17043` | `17043`, no decimal point | Disputed source value; no decimal inserted. |
| 11 | 1229210 | 2023/2024 cut-off | Empty | Damaged characters, roughly `N:/\` | Uncertain; not empty, not `N/A`, not zero. |
| 11 | 1229210 | Name | `…Earth Science, With !!)` | Probably `With IT)`, medium confidence | Uncertain; extraction kept, no correction accepted. |
| 39 | 1229299 | 2022/2023 cut-off | `22.544` | `22 544`, no visible decimal point | Uncertain; not copied from neighbouring rows. |
| 47 | 1229B52 | Programme code | `1229B52` | `1229B52`, the B is clear | Disputed source value; B is not assumed to be a digit. |
| 82 | 1229765 | 2023/2024 cut-off | `18.692` | `18.69` then a damaged final character | Uncertain. |

The other 170 cut-off cells and 82 extracted clusters match the scan, including every `N/A`. Minor notes that do not change values (rows 40, 43, 53, 63 and 76) are recorded in the reviewed table.

## Remaining evidence

CSV publisher/cycle verification and authoritative cross-source programme mappings remain outstanding. Until they are resolved, every CSV requirement stays non-executable, undated CSV cluster points stay unverified, and no programme mapping is approved. None of these inputs supplies an official cluster formula. Architecture blueprint proposals stay separate from student and admissions evidence. This Bead remains in progress; all qualification use stays disabled.

## CSV spot checks against the university catalogue

Maseno's [official programme catalogue](https://maseno.ac.ke/all-programmes), checked 9 October 2026, contains material differences from the supplied CSV. The catalogue does not identify an admissions cycle for these statements:

- CSV record 2 permits English or Kiswahili for Certificate in Information Technology; the university catalogue specifies English alongside Mathematics.
- CSV record 12 presents a standalone C-minus route for Diploma in Business Administration with subject conditions. The catalogue gives a C-plain route or C-minus with a recognised certificate.

These checks establish reasons to withhold CSV rule approval. They do not validate the other records, date the CSV, or approve programme equivalence. A complete row-by-row source review and authoritative cycle remain required.
