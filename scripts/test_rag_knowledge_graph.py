"""Acceptance checks against the actual supplied corpus."""

import json
import tempfile
import unittest
from collections import Counter
from pathlib import Path

from build_rag_knowledge_graph import DATA, PDF_NAME, Graph, build, cutoff, import_cutoffs


class CorpusTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.graph = build()

    def test_coverage_and_integrity(self):
        nodes = self.graph["nodes"]
        ids = {n["id"] for n in nodes}
        self.assertEqual(len(ids), len(nodes))
        counts = Counter(n["label"] for n in nodes)
        self.assertEqual(counts["SourceArtifact"], 3)
        self.assertEqual(counts["ProgrammeRecord"], 164)
        self.assertEqual(counts["AdmissionRequirement"], 231)
        self.assertEqual(counts["CutoffClaim"], 209)
        for edge in self.graph["edges"]:
            self.assertIn(edge["source"], ids)
            self.assertIn(edge["target"], ids)
        for node in nodes:
            if node["label"] != "SourceArtifact":
                self.assertIn(node["properties"]["source_id"], ids)
                self.assertTrue(node["properties"]["location"])

    def test_historical_example_from_pdf(self):
        programme = next(
            n for n in self.graph["nodes"] if n["properties"].get("programme_code_raw") == "1229115"
        )
        cutoffs = {
            e["target"]
            for e in self.graph["edges"]
            if e["source"] == programme["id"] and e["type"] == "STATES_CUTOFF"
        }
        values = {
            n["properties"]["cohort"]: n["properties"]["points_candidate"]
            for n in self.graph["nodes"]
            if n["id"] in cutoffs
        }
        self.assertEqual(values, {"2022/2023": 36.045, "2023/2024": 37.840})

    def test_unknown_cutoffs_are_not_zero_or_eligible(self):
        self.assertEqual(cutoff("N/A"), (None, "not_available"))
        self.assertEqual(cutoff("17043"), (None, "unreadable_or_nonstandard"))
        self.assertEqual(cutoff("49.000"), (None, "outside_cluster_points_range"))
        for node in self.graph["nodes"]:
            if node["label"] == "CutoffClaim":
                self.assertFalse(node["properties"]["usable_for_eligibility"])
            if node["label"] == "AdmissionRequirement":
                self.assertFalse(node["properties"]["executable"])

    def test_headers_and_footers_do_not_become_claims(self):
        for node in self.graph["nodes"]:
            if node["label"] == "CutoffClaim":
                raw = node["properties"]["raw_value"]
                self.assertNotIn("Cut off", raw)
                self.assertNotIn("SEP", raw)
            if node["label"] == "ProgrammeRecord":
                self.assertNotIn("Programme Name", node["properties"]["name"])

    def test_changed_pdf_cannot_reuse_old_ocr(self):
        ocr = json.loads((DATA / "extracted/cutoffs.ocr.json").read_text())
        ocr["source_sha256"] = "incorrect"
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "ocr.json"
            path.write_text(json.dumps(ocr))
            with self.assertRaisesRegex(ValueError, "checksum"):
                import_cutoffs(Graph(), DATA / PDF_NAME, path)

    def test_rebuild_is_deterministic(self):
        self.assertEqual(self.graph, build())

    def test_blueprint_is_proposal_evidence(self):
        proposals = [node for node in self.graph["nodes"] if node["label"].endswith("Proposal")]
        self.assertEqual(len(proposals), 25)
        for node in proposals:
            self.assertFalse(node["properties"]["implemented"])
            self.assertEqual(node["properties"]["adoption_status"], "proposed")
            self.assertTrue(node["properties"]["evidence_quote"])
            self.assertRegex(node["properties"]["location"], r"PDF page [123]")
        states = next(
            node for node in proposals if node["properties"]["name"] == "Publication states"
        )
        self.assertEqual(
            states["properties"]["evidence_quote"], "DRAFT, REVIEW, APPROVED, PUBLISHED"
        )


if __name__ == "__main__":
    unittest.main()
