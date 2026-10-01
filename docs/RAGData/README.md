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

`extracted/cutoffs.ocr.json` preserves all four pages' OCR text, confidence and normalised
bounding boxes. Coordinates use the bottom-left origin. Its checksum must match the PDF.
PDF table rows are reconstructed using programme-code anchors and column positions.
This importer intentionally targets this four-page layout and rejects incomplete coverage.

All PDF claims require review against the scan, including apparently clean numeric cells.
OCR confidence does not establish correctness. Particular extraction problems appear in
the report; missing cells and `N/A` never become zero. The original document appears to
print `17043` for Horticulture in 2023/2024. This is preserved and flagged, not corrected
to an assumed value. The PDF contains both 2022/2023 and 2023/2024 cut-offs.

CSV requirements remain unverified natural-language claims. Subject alternatives, combined
conditions and postgraduate entry paths are not converted into executable rules. CSV
cluster points have no stated cohort. No cutoff claim is authorised for automated eligibility.
The PDF's institutional header does not authenticate the supplied copy, and the CSV's
original authority and date still need confirmation.

To regenerate OCR on macOS using PDFKit and Apple Vision:

```bash
swift scripts/ocr_rag_cutoffs.swift \
  docs/RAGData/CUT-OFF-POINTS-FOR-2023-2024-COHORT.pdf \
  docs/RAGData/extracted/cutoffs.ocr.json
python3 scripts/build_rag_knowledge_graph.py
```

Review OCR changes before rebuilding or publishing the resulting graph.

The blueprint's extracted text is pinned to its PDF checksum in
`extracted/eduplan.blueprint.json`. Import verifies each proposal's quotation against
that extraction. A changed PDF requires fresh extraction and review. The initial extraction
used pypdf 6.19.0 with Unicode NFKC normalisation; graph rebuilds do not require pypdf.

The user selected the existing Go/Angular/Python stack with LlamaIndex for the Python
RAG integration. [Architecture adaptation](architecture-adaptation.md) maps the blueprint
to that stack and identifies schema/security gaps. No runtime integration is claimed.

The broader school assistant requirements, including curriculum retrieval, semantic search,
lesson-plan generation/versioning/publishing, admin management and access control, are
tracked in Beads epic `schoolCRM-o58`. Curriculum files and publishing permissions remain
unresolved. These artifacts implement admissions data preparation `schoolCRM-5pp`
and blueprint incorporation `schoolCRM-4tv`.
