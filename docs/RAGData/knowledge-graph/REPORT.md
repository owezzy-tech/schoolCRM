# School RAG knowledge graph build report

Nodes: 678
Relationships: 1391

## Node counts

- AdmissionRequirement: 231
- ClusterRecord: 19
- ComponentProposal: 9
- CutoffClaim: 209
- DataEntityProposal: 3
- FeatureProposal: 10
- ProgrammeRecord: 164
- SchoolRecord: 27
- ServiceTargetProposal: 3
- SourceArtifact: 3

## Source authority

- `docs/RAGData/EduPlan Studio Architecture Blueprint.pdf`: user_supplied; no publisher identity check
- `docs/RAGData/CUT-OFF-POINTS-FOR-2023-2024-COHORT.pdf`: official_historical_source; byte_identical_to_official_https_download
- `docs/RAGData/maseno_university_all_programmes_expanded.csv`: unverified_user_supplied; no publisher identity check

## Cut-off table transcription review

Revision `maseno-cutoff-scan-review-2026-10-09-v1` (2026-10-09) covers 87 of 87 printed rows: code, name, school, cluster and both cohort cut-offs.
Reviewers: parent agent; independent Claude visual review.
Transcription review of the printed historical table only. It does not establish Maseno approval of this interpretation, current admissions policy, or any eligibility rule.

### Resolved scan corrections (original extraction retained)

- PDF page 1, S.No 14, code 1229179: name = 'Bachelor of Arts (Drama and Theatre Studies, With IT)'; extracted 'Bachelor of Arts (Drama and Theatre Studies, With Il)'
- PDF page 1, S.No 14, code 1229179: cluster = 3; extracted ''
- PDF page 1, S.No 15, code 1229509: cluster = 3; extracted ''
- PDF page 1, S.No 20, code 1229621: name = 'Bachelor of Arts (Language and Communication, With IT)'; extracted 'Bachelor of Arts (Language and Communication,'
- PDF page 2, S.No 30, code 1229337: name = 'Bachelor of Arts (Counseling Psychology, with IT)'; extracted 'Bachelor.of Arts (Counseling Psychology, with IT)'
- PDF page 2, S.No 39, code 1229299: name = 'Bachelor of Business Entrepreneurship, with IT'; extracted 'Barhelor of Rusiness Entrepreneurship, with IT'
- PDF page 2, S.No 42, code 1229539: cluster = 10; extracted ''
- PDF page 3, S.No 52, code 1229612: name = 'Bachelor of Science (Information Systems)'; extracted 'Bachelor of Science (Information Systems) •'
- PDF page 3, S.No 56, code 1229137: cluster = 19; extracted ''
- PDF page 3, S.No 57, code 1229155: cluster = 19; extracted ''
- PDF page 4, S.No 82, code 1229765: name = 'Bachelor of Arts (Urban and Regional Planning, With IT)'; extracted 'Bachelor or Arts (Urban and Regionai Planning, With IT)'

### Unresolved: uncertain or disputed (never numeric or eligible)

- PDF page 1, S.No 7, code 1229185: 2023/2024 = '17043'; Printed '17043' with no decimal point; no decimal is inserted. [disputed_source_value]
- PDF page 1, S.No 11, code 1229210: name = 'Bachelor of Science (Earth Science, With !!)'; extracted 'Bachelor of Science (Earth Science, With !!)'; Final characters damaged; probable reading 'With IT)' at medium confidence, not accepted as a correction. [uncertain]
- PDF page 1, S.No 11, code 1229210: 2023/2024 = None; extracted ''; Damaged characters, approximately 'N:/\'; the shape suggests N/A but is not readable. [uncertain]
- PDF page 2, S.No 39, code 1229299: 2022/2023 = None; extracted '22.544'; Scan shows '22 544' with no visible decimal point when zoomed. [uncertain]
- PDF page 2, S.No 47, code 1229B52: code = '1229B52'; Printed '1229B52'; the B is clear and is not assumed to be a digit. [disputed_source_value]
- PDF page 4, S.No 82, code 1229765: 2023/2024 = None; extracted '18.692'; Scan shows '18.69' followed by a damaged final character. [uncertain]

## Trust and coverage

77 CSV programme records and 87 PDF programme records are source-scoped.
PDF cut-offs are historical 2022/2023 and 2023/2024 published claims only.
They are not current policy and no cutoff claim is usable for eligibility.
N/A is unavailable, never zero.
The scan review checks transcription; it is not institutional approval of its use.

## Blocked: CSV authority and programme mappings

CSV publisher, revision and applicable cycle are unverified.
All CSV requirements are unverified text claims, not executable eligibility rules.
CSV cluster points have no cohort; they are not historical or current thresholds.
No cross-source programme equivalence is inferred or approved.

## Other inputs

The blueprint contributes quoted proposals; keep them outside curriculum retrieval.
Blueprint components are not installed or implemented by this graph build.
Selected architecture: existing Go/Angular/Python RAG with LlamaIndex in Python.
See ../architecture-adaptation.md for user decisions and blueprint reconciliation.
These inputs lack cluster formulas, teaching curriculum and actual lesson plans.
