"""Acceptance checks against the actual supplied corpus."""

import json
import shutil
import tempfile
import unittest
from collections import Counter
from pathlib import Path

from build_rag_knowledge_graph import (
    CUTOFF_REVIEW,
    DATA,
    PDF_NAME,
    PROVENANCE,
    build,
    cutoff_points,
)


class CorpusTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.graph = build()
        cls.nodes = {n["id"]: n for n in cls.graph["nodes"]}
        cls.programmes = {
            n["properties"]["programme_code"]: n
            for n in cls.graph["nodes"]
            if n["properties"].get("record_type") == "cutoff_table"
        }

    def linked(self, programme, relation):
        return [
            self.nodes[e["target"]]["properties"]
            for e in self.graph["edges"]
            if e["source"] == self.programmes[programme]["id"] and e["type"] == relation
        ]

    def cutoffs(self, programme):
        return {p["cohort"]: p for p in self.linked(programme, "STATES_CUTOFF")}

    def test_coverage_and_integrity(self):
        nodes = self.graph["nodes"]
        self.assertEqual(len(self.nodes), len(nodes))
        counts = Counter(n["label"] for n in nodes)
        self.assertEqual(counts["SourceArtifact"], 3)
        self.assertEqual(counts["ProgrammeRecord"], 164)
        self.assertEqual(counts["AdmissionRequirement"], 231)
        self.assertEqual(counts["CutoffClaim"], 209)
        self.assertEqual(len(self.programmes), 87)
        for edge in self.graph["edges"]:
            self.assertIn(edge["source"], self.nodes)
            self.assertIn(edge["target"], self.nodes)
        for node in nodes:
            if node["label"] != "SourceArtifact":
                self.assertIn(node["properties"]["source_id"], self.nodes)
                self.assertTrue(node["properties"]["location"])
        (review,) = self.graph["transcription_reviews"].values()
        self.assertEqual(review["rows_reviewed"], 87)
        self.assertEqual(review["review_revision"], "maseno-cutoff-scan-review-2026-10-09-v1")
        self.assertEqual(
            self.programmes["1229179"]["properties"]["location"],
            "PDF page 1, S.No 14, code 1229179",
        )

    def test_historical_example_from_pdf(self):
        values = {c: p["points_candidate"] for c, p in self.cutoffs("1229115").items()}
        self.assertEqual(values, {"2022/2023": 36.045, "2023/2024": 37.840})

    def test_verified_scan_corrections_retain_extraction(self):
        for code, number in (
            ("1229179", 3),
            ("1229509", 3),
            ("1229539", 10),
            ("1229137", 19),
            ("1229155", 19),
        ):
            self.assertEqual(
                [c["number"] for c in self.linked(code, "ASSIGNED_TO_CLUSTER")], [number]
            )
            review = self.programmes[code]["properties"]["field_reviews"]["cluster"]
            self.assertEqual(review, {"status": "corrected", "extracted_text": ""})
        names = {
            "1229179": ("Bachelor of Arts (Drama and Theatre Studies, With IT)", "With Il)"),
            "1229621": ("Bachelor of Arts (Language and Communication, With IT)", "Communication,"),
            "1229337": ("Bachelor of Arts (Counseling Psychology, with IT)", "Bachelor.of"),
            "1229299": ("Bachelor of Business Entrepreneurship, with IT", "Barhelor of Rusiness"),
            "1229612": ("Bachelor of Science (Information Systems)", "Systems) •"),
            "1229765": ("Bachelor of Arts (Urban and Regional Planning, With IT)", "Regionai"),
        }
        for code, (name, extracted) in names.items():
            properties = self.programmes[code]["properties"]
            self.assertEqual(properties["name"], name)
            self.assertIn(extracted, properties["field_reviews"]["name"]["extracted_text"])
        self.assertEqual(sum(i["status"] == "corrected" for i in self.graph["review_issues"]), 11)

    def test_uncertain_or_disputed_cells_never_become_numbers(self):
        for code, cohort, extracted in (
            ("1229210", "2023/2024", ""),
            ("1229299", "2022/2023", "22.544"),
            ("1229765", "2023/2024", "18.692"),
        ):
            claim = self.cutoffs(code)[cohort]
            self.assertEqual(claim["transcription_review"], "uncertain")
            self.assertIsNone(claim["transcription"])
            self.assertEqual(claim["extracted_text"], extracted)
            self.assertTrue(claim["scan_observation"])
            self.assertNotIn("points_candidate", claim)
            self.assertTrue(claim["review_required"])
        horticulture = self.cutoffs("1229185")["2023/2024"]
        self.assertEqual(horticulture["transcription"], "17043")
        self.assertEqual(horticulture["transcription_review"], "disputed_source_value")
        self.assertNotIn("points_candidate", horticulture)
        ecotourism = self.programmes["1229B52"]["properties"]
        self.assertEqual(ecotourism["field_reviews"]["code"]["status"], "disputed_source_value")
        self.assertTrue(ecotourism["review_required"])
        earth = self.programmes["1229210"]["properties"]
        self.assertEqual(earth["name"], "Bachelor of Science (Earth Science, With !!)")
        self.assertEqual(earth["field_reviews"]["name"]["status"], "uncertain")
        # 174 PDF cells minus four unresolved cells and 24 printed N/A values.
        numeric = [n for n in self.graph["nodes"] if "points_candidate" in n["properties"]]
        self.assertEqual(len(numeric), 146)

    def test_unknown_cutoffs_are_not_zero_or_eligible(self):
        self.assertIsNone(cutoff_points("N/A"))
        for text in ("17043", "22 544", "49.000", ""):
            with self.assertRaises(ValueError):
                cutoff_points(text)
        for node in self.graph["nodes"]:
            if node["label"] == "CutoffClaim":
                self.assertFalse(node["properties"]["usable_for_eligibility"])
                self.assertNotEqual(node["properties"].get("points_candidate"), 0)
            if node["label"] == "AdmissionRequirement":
                self.assertFalse(node["properties"]["executable"])

    def test_cross_source_claims_stay_unverified_and_unmapped(self):
        sources = {
            n["properties"]["path"].rsplit("/", 1)[1]: n["properties"]
            for n in self.graph["nodes"]
            if n["label"] == "SourceArtifact"
        }
        self.assertEqual(sources[PDF_NAME]["authority_status"], "official_historical_source")
        self.assertFalse(sources[PDF_NAME]["current_policy"])
        csv = sources["maseno_university_all_programmes_expanded.csv"]
        self.assertEqual(csv["authority_status"], "unverified_user_supplied")
        self.assertIsNone(csv["applicable_cycle"])
        self.assertFalse(csv["programme_mappings_approved"])
        for node in self.graph["nodes"]:
            if node["properties"].get("interpretation") == "undated_csv_cluster_points":
                self.assertEqual(node["properties"]["cohort"], "unspecified")
                self.assertNotIn("points_candidate", node["properties"])
        programme_ids = {n["id"] for n in self.graph["nodes"] if n["label"] == "ProgrammeRecord"}
        for edge in self.graph["edges"]:
            self.assertFalse(edge["source"] in programme_ids and edge["target"] in programme_ids)

    def test_changed_pdf_ocr_review_or_provenance_fail_closed(self):
        def rejected(edit, pattern, target=CUTOFF_REVIEW):
            with tempfile.TemporaryDirectory() as directory:
                data = Path(directory)
                shutil.copytree(DATA, data, dirs_exist_ok=True)
                document = json.loads((data / target).read_text())
                edit(document)
                (data / target).write_text(json.dumps(document))
                with self.assertRaisesRegex(ValueError, pattern):
                    build(data)

        def pin(field, value):
            return lambda document: document.__setitem__(field, value)

        def edit_row(serial, update):
            return lambda document: document["rows"][serial - 1].update(update)

        rejected(pin("source_sha256", "0" * 64), "different PDF")
        rejected(pin("extraction_sha256", "0" * 64), "OCR evidence checksum")
        rejected(
            pin("source_sha256", "incorrect"), "OCR source checksum", "extracted/cutoffs.ocr.json"
        )
        rejected(edit_row(39, {"2022/2023": "22.544"}), "Invalid 2022/2023 review")
        rejected(edit_row(7, {"review": {}}), "not a cluster-points value")
        rejected(edit_row(1, {"cluster": None}), "Unreviewed empty cluster")
        rejected(lambda d: d["rows"].pop(), "S.No 1-87")
        rejected(
            lambda d: d["sources"][0].__setitem__("sha256", "0" * 64), "provenance", PROVENANCE
        )

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
