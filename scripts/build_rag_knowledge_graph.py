"""Build source-scoped admissions graphs. No model calls or database writes."""

import argparse
import csv
import hashlib
import json
import re
import xml.etree.ElementTree as ET
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DATA = ROOT / "docs/RAGData"
CSV_NAME = "maseno_university_all_programmes_expanded.csv"
PDF_NAME = "CUT-OFF-POINTS-FOR-2023-2024-COHORT.pdf"
BLUEPRINT_NAME = "EduPlan Studio Architecture Blueprint.pdf"
CUTOFF_REVIEW = "review/cutoff-table-review.json"
PROVENANCE = "review/source-provenance.json"
COHORTS = ("2022/2023", "2023/2024")
UNRESOLVED = ("uncertain", "disputed_source_value")


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def identity(kind, *parts):
    value = json.dumps(parts, ensure_ascii=False, separators=(",", ":"))
    return kind + ":" + hashlib.sha256(value.encode()).hexdigest()[:24]


class Graph:
    def __init__(self, provenance=()):
        self.nodes = {}
        self.edges = set()
        self.issues = []
        self.reviews = {}
        self.provenance = {entry["file"]: entry for entry in provenance}

    def node(self, label, key, **properties):
        value = {"id": key, "label": label, "properties": properties}
        if key in self.nodes and self.nodes[key] != value:
            raise ValueError(f"Conflicting node: {key}")
        self.nodes[key] = value
        return key

    def edge(self, source, relation, target):
        self.edges.add((source, relation, target))

    def claim(self, label, key, source, location, **properties):
        self.node(
            label,
            key,
            source_id=source,
            location=location,
            verification_status="unverified",
            **properties,
        )
        self.edge(key, "SUPPORTED_BY", source)
        return key

    def source(self, path):
        checksum = digest(path)
        key = identity("source", path.name, checksum)
        # Source identity verification is pinned separately from per-cell transcription review.
        verified = {"authority_status": "user_supplied"}
        if path.name in self.provenance:
            verified = dict(self.provenance[path.name])
            if verified.pop("sha256") != checksum:
                raise ValueError(f"Source provenance is pinned to different bytes: {path.name}")
            del verified["file"]
        return self.node(
            "SourceArtifact",
            key,
            path="docs/RAGData/" + path.name,
            sha256=checksum,
            **verified,
        ), checksum

    def document(self):
        for source, _, target in self.edges:
            if source not in self.nodes or target not in self.nodes:
                raise ValueError("Dangling relationship endpoint")
        return {
            "schema_version": 1,
            "nodes": sorted(self.nodes.values(), key=lambda n: n["id"]),
            "edges": [{"source": s, "type": r, "target": t} for s, r, t in sorted(self.edges)],
            "transcription_reviews": self.reviews,
            "review_issues": self.issues,
        }


def import_catalog(graph, path):
    source, _ = graph.source(path)
    seen = set()
    with path.open(encoding="utf-8-sig", newline="") as stream:
        reader = csv.DictReader(stream)
        expected = {
            "category",
            "school",
            "program_name",
            "duration_years",
            "min_kcse_mean_grade",
            "subject_requirements",
            "alternate_entry",
            "cluster_points",
        }
        if set(reader.fieldnames or []) != expected:
            raise ValueError("Unexpected programme CSV columns")
        for row_number, row in enumerate(reader, 2):
            if any(not isinstance(v, str) or not v.strip() for v in row.values()):
                raise ValueError(f"Missing CSV value at row {row_number}")
            row = {k: v.strip() for k, v in row.items()}
            name = row["program_name"]
            if name in seen:
                raise ValueError(f"Duplicate CSV programme: {name}")
            seen.add(name)
            location = f"CSV record {row_number} (header is record 1)"
            programme = identity("catalog_programme", source, name)
            graph.claim(
                "ProgrammeRecord",
                programme,
                source,
                location,
                name=name,
                category=row["category"],
                duration_text=row["duration_years"],
                record_type="catalogue",
                institution_name="Maseno University",
            )
            school = identity("school_record", source, row["school"])
            graph.claim("SchoolRecord", school, source, "CSV school column", name=row["school"])
            graph.edge(school, "LISTS_PROGRAMME", programme)
            for field in ("min_kcse_mean_grade", "subject_requirements", "alternate_entry"):
                claim = identity("requirement", programme, field)
                graph.claim(
                    "AdmissionRequirement",
                    claim,
                    source,
                    location,
                    kind=field,
                    text=row[field],
                    executable=False,
                )
                graph.edge(programme, "STATES_REQUIREMENT", claim)
            if row["cluster_points"] != "N/A":
                claim = identity("undated_cluster_points", programme)
                graph.claim(
                    "CutoffClaim",
                    claim,
                    source,
                    location,
                    raw_value=row["cluster_points"],
                    cohort="unspecified",
                    interpretation="undated_csv_cluster_points",
                    usable_for_eligibility=False,
                )
                graph.edge(programme, "STATES_CUTOFF", claim)


def cutoff_points(text):
    """Return reviewed historical points, or None for a printed N/A. Never zero."""
    if text == "N/A":
        return None
    if not re.fullmatch(r"\d{1,2}\.\d{3}", text) or not 0 <= float(text) <= 48:
        raise ValueError(f"Reviewed cutoff is not a cluster-points value: {text!r}")
    return float(text)


def field_review(row, field):
    """Validate one reviewed cell; unlisted cells were reviewed and accepted as extracted."""
    review = row.get("review", {}).get(field)
    value = row[field]
    if review is None:
        if value is None:
            raise ValueError(f"Unreviewed empty {field} in S.No {row['serial']}")
        return "reviewed", {}
    status = review["status"]
    if status == "corrected" and value is not None and review["extracted"] != value:
        return status, {"extracted_text": review["extracted"]}
    if status == "uncertain" and review["observed"] and (value is None or field == "name"):
        return status, {
            "extracted_text": review["extracted"],
            "scan_observation": review["observed"],
        }
    if status == "disputed_source_value" and value is not None and review["observed"]:
        return status, {"scan_observation": review["observed"]}
    raise ValueError(f"Invalid {field} review in S.No {row['serial']}")


def import_cutoffs(graph, pdf, review_path):
    """Import the reviewed transcription; OCR stays pinned audit evidence, not a build input."""
    source, checksum = graph.source(pdf)
    review = json.loads(review_path.read_text())
    extraction = pdf.parent / review["extraction_file"]
    if review["source_sha256"] != checksum or review["source_file"] != pdf.name:
        raise ValueError("Cutoff review is pinned to a different PDF; review the new scan")
    if json.loads(extraction.read_text())["source_sha256"] != checksum:
        raise ValueError("OCR source checksum does not match PDF; regenerate OCR")
    if review["extraction_sha256"] != digest(extraction):
        raise ValueError("OCR evidence checksum differs from the reviewed extraction")
    rows = review["rows"]
    if [row["serial"] for row in rows] != list(range(1, 88)):
        raise ValueError("Expected reviewed S.No 1-87 in printed order")
    if len({row["code"] for row in rows}) != len(rows):
        raise ValueError("Duplicate PDF programme code")
    pinned = {"review_revision": review["review_revision"], "reviewed_on": review["reviewed_on"]}
    graph.reviews[source] = {
        **pinned,
        "reviewers": review["reviewers"],
        "scope": review["scope"],
        "rows_reviewed": len(rows),
    }
    for row in rows:
        location = f"PDF page {row['page']}, S.No {row['serial']}, code {row['code']}"
        reviews = {}
        for field in ("code", "name", "cluster", *COHORTS):
            status, evidence = field_review(row, field)
            reviews[field] = {"status": status, **evidence}
            if status != "reviewed":
                graph.issues.append(
                    {"location": location, "field": field, "reviewed": row[field], **reviews[field]}
                )
        exceptions = {
            field: review
            for field, review in reviews.items()
            if field not in COHORTS and review["status"] != "reviewed"
        }
        programme = identity("pdf_programme", source, row["code"])
        graph.claim(
            "ProgrammeRecord",
            programme,
            source,
            location,
            name=row["name"],
            programme_code=row["code"],
            record_type="cutoff_table",
            institution_name="Maseno University",
            field_reviews=exceptions,
            review_note=row.get("note"),
            review_required=any(r["status"] in UNRESOLVED for r in exceptions.values()),
            **pinned,
        )
        school = identity("school_record", source, row["school"])
        graph.claim("SchoolRecord", school, source, "PDF school heading", name=row["school"])
        graph.edge(school, "LISTS_PROGRAMME", programme)
        if row["cluster"] is not None:
            cluster = identity("cluster", source, str(row["cluster"]))
            graph.claim(
                "ClusterRecord", cluster, source, "PDF cluster column", number=row["cluster"]
            )
            graph.edge(programme, "ASSIGNED_TO_CLUSTER", cluster)
        for cohort in COHORTS:
            evidence = dict(reviews[cohort])
            status = evidence.pop("status")
            properties = {
                "cohort": cohort,
                "transcription": row[cohort],
                "transcription_review": status,
                "review_required": status in UNRESOLVED,
                "usable_for_eligibility": False,
                **evidence,
                **pinned,
            }
            # Uncertain and disputed cells keep their evidence but never become numbers.
            if status not in UNRESOLVED:
                points = cutoff_points(row[cohort])
                properties["availability"] = "printed" if points is not None else "not_available"
                if points is not None:
                    properties["points_candidate"] = points
            claim = identity("cutoff", programme, cohort)
            graph.claim("CutoffClaim", claim, source, location, **properties)
            graph.edge(programme, "STATES_CUTOFF", claim)


def import_blueprint(graph, pdf, extraction_path):
    source, checksum = graph.source(pdf)
    graph.nodes[source]["properties"]["source_role"] = "architecture_proposal"
    extraction = json.loads(extraction_path.read_text())
    if extraction["source_sha256"] != checksum:
        raise ValueError("Blueprint extraction checksum does not match PDF")
    pages = {page["page"]: " ".join(page["text"].split()) for page in extraction["pages"]}
    if set(pages) != {1, 2, 3}:
        raise ValueError("Expected all three blueprint pages")
    # Explicit source quotations, not claims about installed components or working features.
    claims = [
        ("FeatureProposal", 1, "Curriculum discovery", "curriculum resource discovery"),
        ("FeatureProposal", 1, "AI lesson generation", "AI-assisted lesson plan generation"),
        ("FeatureProposal", 1, "Approval workflow", "strict approval workflows"),
        ("FeatureProposal", 1, "Immutable version history", "immutable version control"),
        ("FeatureProposal", 1, "PDF export", "PDF Export (Maroto v2)"),
        ("FeatureProposal", 2, "Structured generation", "output strict JSON"),
        ("FeatureProposal", 2, "Publication states", "DRAFT, REVIEW, APPROVED, PUBLISHED"),
        (
            "FeatureProposal",
            3,
            "School and department isolation",
            "prevents unauthorized access across schools, departments, and user roles.",
        ),
        ("FeatureProposal", 3, "Admin user management", "p, admin, users, manage"),
        (
            "FeatureProposal",
            3,
            "Teacher plan creation and reading",
            "p, teacher, lesson_plans, create p, teacher, lesson_plans, read",
        ),
        ("ComponentProposal", 1, "React frontend", "React SPA Client"),
        ("ComponentProposal", 1, "Go router", "go-chi Router"),
        ("ComponentProposal", 1, "Casbin authorisation", "Casbin RBAC Auth"),
        ("ComponentProposal", 1, "LangChainGo lesson engine", "Lesson Engine (LangChainGo)"),
        ("ComponentProposal", 1, "pgvector retrieval", "RAG Search (pgvector)"),
        ("ComponentProposal", 1, "Temporal workflows", "Workflow Engine (Temporal)"),
        ("ComponentProposal", 1, "PostgreSQL 16", "PostgreSQL 16 Storage"),
        ("ComponentProposal", 1, "SQLC data access", "Relational Data (SQLC)"),
        (
            "ComponentProposal",
            3,
            "Redis and Asynq queues",
            "Redis cluster backed background queue isolation via Asynq.",
        ),
        ("DataEntityProposal", 2, "Lesson plan", "CREATE TABLE lesson_plans"),
        ("DataEntityProposal", 2, "Lesson plan version", "CREATE TABLE lesson_plan_versions"),
        ("DataEntityProposal", 2, "Curriculum resource", "CREATE TABLE curriculum_resources"),
        ("ServiceTargetProposal", 3, "Availability", "99.9% Uptime"),
        ("ServiceTargetProposal", 3, "API latency", "Sub-100ms API Latency"),
        (
            "ServiceTargetProposal",
            3,
            "RAG latency",
            "RAG Response: <1.2s embedding & extraction loop.",
        ),
    ]
    for label, page, name, quote in claims:
        if quote not in pages[page]:
            raise ValueError(f"Blueprint quotation not found on page {page}: {quote}")
        key = identity("blueprint_claim", source, label, name)
        graph.claim(
            label,
            key,
            source,
            f"PDF page {page}",
            name=name,
            evidence_quote=quote,
            adoption_status="proposed",
            implemented=False,
        )
        graph.edge(source, "PROPOSES", key)


def build(data_dir=DATA):
    graph = Graph(json.loads((data_dir / PROVENANCE).read_text())["sources"])
    import_catalog(graph, data_dir / CSV_NAME)
    import_cutoffs(graph, data_dir / PDF_NAME, data_dir / CUTOFF_REVIEW)
    import_blueprint(
        graph, data_dir / BLUEPRINT_NAME, data_dir / "extracted/eduplan.blueprint.json"
    )
    return graph.document()


def write_outputs(graph, output):
    output.mkdir(parents=True, exist_ok=True)
    (output / "graph.json").write_text(json.dumps(graph, indent=2, ensure_ascii=False) + "\n")
    counts = Counter(n["label"] for n in graph["nodes"])
    report = [
        "# School RAG knowledge graph build report",
        "",
        f"Nodes: {len(graph['nodes'])}",
        f"Relationships: {len(graph['edges'])}",
        "",
        "## Node counts",
        "",
    ]
    report.extend(f"- {kind}: {count}" for kind, count in sorted(counts.items()))
    report.extend(["", "## Source authority", ""])
    for node in graph["nodes"]:
        if node["label"] == "SourceArtifact":
            p = node["properties"]
            check = p.get("identity_check", "no publisher identity check")
            report.append(f"- `{p['path']}`: {p['authority_status']}; {check}")
    report.extend(["", "## Cut-off table transcription review", ""])
    for review in graph["transcription_reviews"].values():
        report.extend(
            [
                f"Revision `{review['review_revision']}` ({review['reviewed_on']}) covers "
                f"{review['rows_reviewed']} of 87 printed rows: code, name, school, cluster and "
                "both cohort cut-offs.",
                "Reviewers: " + "; ".join(r["reviewer"] for r in review["reviewers"]) + ".",
                review["scope"],
            ]
        )

    def describe(issue):
        evidence = [f"extracted {issue['extracted_text']!r}"] if "extracted_text" in issue else []
        evidence += [issue["scan_observation"]] if "scan_observation" in issue else []
        return f"- {issue['location']}: {issue['field']} = {issue['reviewed']!r}; " + "; ".join(
            evidence
        )

    issues = graph["review_issues"]
    report.extend(["", "### Resolved scan corrections (original extraction retained)", ""])
    report.extend(describe(i) for i in issues if i["status"] == "corrected")
    report.extend(["", "### Unresolved: uncertain or disputed (never numeric or eligible)", ""])
    report.extend(f"{describe(i)} [{i['status']}]" for i in issues if i["status"] != "corrected")
    report.extend(
        [
            "",
            "## Trust and coverage",
            "",
            "77 CSV programme records and 87 PDF programme records are source-scoped.",
            "PDF cut-offs are historical 2022/2023 and 2023/2024 published claims only.",
            "They are not current policy and no cutoff claim is usable for eligibility.",
            "N/A is unavailable, never zero.",
            "The scan review checks transcription; it is not institutional approval of its use.",
            "",
            "## Blocked: CSV authority and programme mappings",
            "",
            "CSV publisher, revision and applicable cycle are unverified.",
            "All CSV requirements are unverified text claims, not executable eligibility rules.",
            "CSV cluster points have no cohort; they are not historical or current thresholds.",
            "No cross-source programme equivalence is inferred or approved.",
            "",
            "## Other inputs",
            "",
            "The blueprint contributes quoted proposals; keep them outside curriculum retrieval.",
            "Blueprint components are not installed or implemented by this graph build.",
            "Selected architecture: existing Go/Angular/Python RAG with LlamaIndex in Python.",
            "See ../architecture-adaptation.md for user decisions and blueprint reconciliation.",
            "These inputs lack cluster formulas, teaching curriculum and actual lesson plans.",
        ]
    )
    (output / "REPORT.md").write_text("\n".join(report) + "\n")
    # Portable, dependency-free GraphML for graph tools. JSON preserves typed properties.
    ns = "http://graphml.graphdrawing.org/xmlns"
    ET.register_namespace("", ns)
    root = ET.Element(f"{{{ns}}}graphml")
    for key in ("label", "properties", "type"):
        ET.SubElement(
            root, f"{{{ns}}}key", id=key, **{"for": "all", "attr.name": key, "attr.type": "string"}
        )
    element = ET.SubElement(root, f"{{{ns}}}graph", id="school_rag", edgedefault="directed")
    for node in graph["nodes"]:
        entry = ET.SubElement(element, f"{{{ns}}}node", id=node["id"])
        ET.SubElement(entry, f"{{{ns}}}data", key="label").text = node["label"]
        ET.SubElement(entry, f"{{{ns}}}data", key="properties").text = json.dumps(
            node["properties"]
        )
    for index, edge in enumerate(graph["edges"]):
        entry = ET.SubElement(
            element, f"{{{ns}}}edge", id=f"e{index}", source=edge["source"], target=edge["target"]
        )
        ET.SubElement(entry, f"{{{ns}}}data", key="type").text = edge["type"]
    ET.ElementTree(root).write(output / "graph.graphml", encoding="utf-8", xml_declaration=True)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--data-dir", type=Path, default=DATA)
    parser.add_argument("--output", type=Path, default=DATA / "knowledge-graph")
    args = parser.parse_args()
    write_outputs(build(args.data_dir), args.output)
    print(f"Built graph in {args.output}")
