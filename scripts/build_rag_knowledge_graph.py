"""Build source-scoped admissions graphs. No model calls or database writes."""

import argparse
import csv
import hashlib
import json
import math
import re
import xml.etree.ElementTree as ET
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DATA = ROOT / "docs/RAGData"
CSV_NAME = "maseno_university_all_programmes_expanded.csv"
PDF_NAME = "CUT-OFF-POINTS-FOR-2023-2024-COHORT.pdf"
BLUEPRINT_NAME = "EduPlan Studio Architecture Blueprint.pdf"


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def identity(kind, *parts):
    value = json.dumps(parts, ensure_ascii=False, separators=(",", ":"))
    return kind + ":" + hashlib.sha256(value.encode()).hexdigest()[:24]


class Graph:
    def __init__(self):
        self.nodes = {}
        self.edges = set()
        self.issues = []

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
        return self.node(
            "SourceArtifact",
            key,
            path="docs/RAGData/" + path.name,
            sha256=checksum,
            authority_status="user_supplied",
        ), checksum

    def document(self):
        for source, _, target in self.edges:
            if source not in self.nodes or target not in self.nodes:
                raise ValueError("Dangling relationship endpoint")
        return {
            "schema_version": 1,
            "nodes": sorted(self.nodes.values(), key=lambda n: n["id"]),
            "edges": [{"source": s, "type": r, "target": t} for s, r, t in sorted(self.edges)],
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


def cutoff(raw):
    if raw == "N/A":
        return None, "not_available"
    if not re.fullmatch(r"\d{1,2}\.\d{3}", raw or ""):
        return None, "unreadable_or_nonstandard"
    value = float(raw)
    if not 0 <= value <= 48:
        return None, "outside_cluster_points_range"
    return value, "ocr_candidate"


def import_cutoffs(graph, pdf, ocr_path):
    source, checksum = graph.source(pdf)
    ocr = json.loads(ocr_path.read_text())
    if ocr["source_sha256"] != checksum:
        raise ValueError("OCR source checksum does not match PDF; regenerate OCR")
    if [p["page"] for p in ocr["pages"]] != [1, 2, 3, 4]:
        raise ValueError("Expected all four cutoff PDF pages")
    school_name = ""
    count = 0
    seen_codes = set()
    for page in ocr["pages"]:
        lines = page["lines"]
        for line in lines:
            if not isinstance(line["text"], str):
                raise ValueError("Invalid OCR text")
            for field in ("x", "y", "width", "height", "confidence"):
                if not math.isfinite(line[field]) or not 0 <= line[field] <= 1:
                    raise ValueError(f"Invalid OCR {field}")
        anchors = []
        for line in lines:
            match = re.search(r"\b[0-9B]{7}\b", line["text"])
            if line["x"] < 0.24 and match:
                anchors.append((line, match.group()))
        anchors.sort(key=lambda a: -a[0]["y"])
        headings = sorted(
            [line for line in lines if line["text"].startswith("SCHOOL OF")],
            key=lambda line: -line["y"],
        )
        for index, (anchor, code) in enumerate(anchors):
            for heading in headings:
                if heading["y"] > anchor["y"]:
                    school_name = heading["text"]
            upper = (anchors[index - 1][0]["y"] + anchor["y"]) / 2 if index else 0.94
            lower = (
                (anchors[index + 1][0]["y"] + anchor["y"]) / 2 if index + 1 < len(anchors) else 0.07
            )
            upper = min(upper, anchor["y"] + 0.02)
            lower = max(lower, anchor["y"] - 0.025)
            cells = [
                line
                for line in lines
                if lower < line["y"] <= upper and not line["text"].startswith("SCHOOL OF")
            ]
            names = sorted(
                [line for line in cells if 0.24 <= line["x"] < 0.63], key=lambda line: -line["y"]
            )
            name = " ".join(line["text"] for line in names)
            if not name or code in seen_codes:
                raise ValueError(f"Missing name or duplicate PDF code: {code}")
            seen_codes.add(code)
            count += 1
            location = f"PDF page {page['page']}, programme code OCR {code}"
            programme = identity("pdf_programme", source, code)
            graph.claim(
                "ProgrammeRecord",
                programme,
                source,
                location,
                name=name,
                programme_code_raw=code,
                record_type="cutoff_table",
                institution_name="Maseno University",
                review_required=True,
                extraction_method=ocr["engine"],
                region_top=upper,
                region_bottom=lower,
            )
            school = identity("school_record", source, school_name)
            graph.claim(
                "SchoolRecord",
                school,
                source,
                "PDF school heading",
                name=school_name,
                review_required=True,
            )
            graph.edge(school, "LISTS_PROGRAMME", programme)
            cluster_cells = [line for line in cells if 0.63 <= line["x"] < 0.72]
            cluster_raw = " ".join(line["text"] for line in cluster_cells)
            if re.fullmatch(r"\d{1,2}", cluster_raw):
                cluster = identity("cluster", source, cluster_raw)
                graph.claim(
                    "ClusterRecord",
                    cluster,
                    source,
                    "PDF cluster column",
                    number=int(cluster_raw),
                    review_required=True,
                )
                graph.edge(programme, "ASSIGNED_TO_CLUSTER", cluster)
            else:
                graph.issues.append({"location": location, "field": "cluster", "raw": cluster_raw})
            if not code.isdigit():
                graph.issues.append({"location": location, "field": "programme_code", "raw": code})
            for cohort, left, right in (("2022/2023", 0.72, 0.82), ("2023/2024", 0.82, 0.95)):
                column = [line for line in cells if left <= line["x"] < right]
                raw = " ".join(line["text"] for line in column)
                value, status = cutoff(raw)
                claim = identity("cutoff", programme, cohort)
                properties = {
                    "cohort": cohort,
                    "raw_value": raw,
                    "extraction_status": status,
                    "review_required": True,
                    "usable_for_eligibility": False,
                }
                if value is not None:
                    properties["points_candidate"] = value
                graph.claim("CutoffClaim", claim, source, location, **properties)
                graph.edge(programme, "STATES_CUTOFF", claim)
                if status not in ("ocr_candidate", "not_available"):
                    graph.issues.append({"location": location, "field": cohort, "raw": raw})
    if count != 87:
        raise ValueError(f"Expected 87 PDF programme rows, extracted {count}")


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
    graph = Graph()
    import_catalog(graph, data_dir / CSV_NAME)
    import_cutoffs(graph, data_dir / PDF_NAME, data_dir / "extracted/cutoffs.ocr.json")
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
    report.extend(
        [
            "",
            "## Trust and coverage",
            "",
            "77 CSV programme records and 87 PDF programme records are source-scoped.",
            "No cross-source programme equivalence has been inferred.",
            "All requirements are unverified text claims, not executable eligibility rules.",
            "Cutoff PDF records need visual review. OCR confidence is not factual certainty.",
            "Both historical cohorts are retained. N/A is unavailable, never zero.",
            "CSV cluster points have no cohort and are not current thresholds.",
            "The blueprint contributes quoted proposals; keep them outside curriculum retrieval.",
            "Blueprint components are not installed or implemented by this graph build.",
            "Selected architecture: existing Go/Angular/Python RAG with LlamaIndex in Python.",
            "See ../architecture-adaptation.md for user decisions and blueprint reconciliation.",
            "These inputs lack cluster formulas, teaching curriculum and actual lesson plans.",
            "",
            "## Cells needing particular attention",
            "",
        ]
    )
    report.extend(f"- {i['location']}: {i['field']} = {i['raw']!r}" for i in graph["review_issues"])
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
