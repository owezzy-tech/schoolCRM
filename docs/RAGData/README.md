# School RAG knowledge graphs

The supplied CSV and cut-off PDF describe Maseno admissions. The EduPlan Studio PDF
describes proposed architecture and features. None contains teaching curriculum,
lesson resources, a cluster calculation formula or current verified admissions policy.

Build the source-scoped graphs into a single portable graph:

```bash
python3 scripts/build_rag_knowledge_graph.py
python3 -m unittest discover -s scripts -p 'test_rag_knowledge_graph.py'
```

Outputs are `knowledge-graph/graph.json`, `knowledge-graph/graph.graphml` and
`knowledge-graph/REPORT.md`. JSON preserves property types; GraphML can be opened in
graph-analysis tools. These are local artifacts, not a deployed graph database or
an integrated RAG endpoint. The build uses Python's standard library.

## Model

- `SourceArtifact` records the supplied filename and SHA-256 checksum.
- `SchoolRecord` lists `ProgrammeRecord` nodes through `LISTS_PROGRAMME`.
- Programmes link to `AdmissionRequirement` through `STATES_REQUIREMENT`.
- PDF programmes link to `ClusterRecord` through `ASSIGNED_TO_CLUSTER` when readable.
- Programmes link to individual `CutoffClaim` nodes through `STATES_CUTOFF`.
- Every sourced node links to its source through `SUPPORTED_BY` and retains a location.
- Blueprint feature/component/data/target proposals have verbatim quotations and PDF page
  references. Its source links to them through `PROPOSES`; they are not installed components.

Identifiers derive from source identity and programme name/code or claim kind.
Different source revisions produce separate records, preserving historical claims.
CSV and PDF programme names are not automatically treated as equivalent. Reviewed
cross-source mappings are needed before combining their claims for qualification advice.

## Evidence and review

Source identity, per-cell transcription review and qualification authority are separate:

- `review/source-provenance.json` pins each source's checksum, publisher and authority.
  The build copies it onto `SourceArtifact` and fails if the bytes differ.
- `review/cutoff-table-review.json` is the reviewed transcription of all 87 printed PDF rows.
  It is pinned to the PDF and OCR checksums and records the review revision and reviewers.
  It is the build input for PDF rows; the build fails if either checksum changes.
- No cutoff claim is usable for eligibility, whatever its review status.

`extracted/cutoffs.ocr.json` preserves all four pages' OCR text, confidence and normalised
bounding boxes as audit evidence. Coordinates use the bottom-left origin. Fields not listed
in a row's `review` were compared with the scan and accept the extracted text. Listed fields
are `corrected` (scan clearly readable; `extracted` keeps the OCR text), `uncertain` (not
reliably readable; no reviewed value) or `disputed_source_value` (clearly printed but
anomalous, such as Horticulture's `17043` and code `1229B52`). Uncertain and disputed cells
never produce `points_candidate`. Printed `N/A` is unavailable, never zero. The PDF holds
historical 2022/2023 and 2023/2024 cut-offs, not current policy or a calculation rule.

CSV requirements remain unverified natural-language claims. Subject alternatives, combined
conditions and postgraduate entry paths are not converted into executable rules. CSV
cluster points have no stated cohort. The CSV's original authority and cycle still need
confirmation, so no CSV-to-PDF programme mapping is approved.

To regenerate OCR on macOS using PDFKit and Apple Vision:

```bash
swift scripts/ocr_rag_cutoffs.swift \
  docs/RAGData/CUT-OFF-POINTS-FOR-2023-2024-COHORT.pdf \
  docs/RAGData/extracted/cutoffs.ocr.json
python3 scripts/build_rag_knowledge_graph.py
```

New OCR output changes the extraction checksum, so the build fails until the reviewed table
is compared with the scan again and given a new review revision.

The blueprint's extracted text is pinned to its PDF checksum in
`extracted/eduplan.blueprint.json`. Import verifies each proposal's quotation against
that extraction. A changed PDF requires fresh extraction and review. The initial extraction
used pypdf 6.19.0 with Unicode NFKC normalisation; graph rebuilds do not require pypdf.

The user selected the existing Go/Angular/Python stack with LlamaIndex for the Python
RAG integration. [Architecture adaptation](architecture-adaptation.md) maps the blueprint
to that stack and identifies schema/security gaps. No runtime integration is claimed.

The broader School AI requirements and delivery status are recorded in
[the ARD](../school-ai-ard.md) and Beads epic `schoolCRM-o58`. These source artifacts
were originally prepared for admissions data preparation `schoolCRM-5pp` and blueprint
incorporation `schoolCRM-4tv`. Their admission-rule suitability is being reviewed in
`schoolCRM-o58.5`; see [the evidence review](admissions-evidence-review.md). Runtime
curriculum and lesson workflows use their own scoped, reviewed sources and Go APIs.
