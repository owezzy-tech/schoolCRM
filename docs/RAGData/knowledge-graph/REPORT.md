# School RAG knowledge graph build report

Nodes: 678
Relationships: 1386

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

## Trust and coverage

77 CSV programme records and 87 PDF programme records are source-scoped.
No cross-source programme equivalence has been inferred.
All requirements are unverified text claims, not executable eligibility rules.
Cutoff PDF records need visual review. OCR confidence is not factual certainty.
Both historical cohorts are retained. N/A is unavailable, never zero.
CSV cluster points have no cohort and are not current thresholds.
The blueprint contributes quoted proposals; keep them outside curriculum retrieval.
Blueprint components are not installed or implemented by this graph build.
Selected architecture: existing Go/Angular/Python RAG with LlamaIndex in Python.
See ../architecture-adaptation.md for user decisions and blueprint reconciliation.
These inputs lack cluster formulas, teaching curriculum and actual lesson plans.

## Cells needing particular attention

- PDF page 1, programme code OCR 1229185: 2023/2024 = '17043'
- PDF page 1, programme code OCR 1229210: 2023/2024 = ''
- PDF page 1, programme code OCR 1229179: cluster = ''
- PDF page 1, programme code OCR 1229509: cluster = ''
- PDF page 2, programme code OCR 1229539: cluster = ''
- PDF page 2, programme code OCR 1229B52: programme_code = '1229B52'
- PDF page 3, programme code OCR 1229137: cluster = ''
- PDF page 3, programme code OCR 1229155: cluster = ''
