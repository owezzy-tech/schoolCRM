"""Verify Git/Beads coordination without changing real issues or contacting a remote."""

import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
HELPER = ROOT / "scripts/git-hooks/beads.sh"


class BeadsHookTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.root = Path(self.directory.name)
        (self.root / ".beads").mkdir()
        self.bin = self.root / "bin"
        self.bin.mkdir()
        self.log = self.root / "calls.jsonl"
        fake = self.bin / "bd"
        fake.write_text(
            f"#!{sys.executable}\n"
            "import json, os, sys\n"
            "with open(os.environ['TEST_BEADS_LOG'], 'a') as log:\n"
            "    log.write(json.dumps(sys.argv[1:]) + '\\n')\n"
            "if ' '.join(sys.argv[1:]) == os.environ.get('TEST_BEADS_FAIL'):\n"
            "    sys.exit(17)\n"
        )
        fake.chmod(0o755)
        self.env = {
            **os.environ,
            "PATH": str(self.bin),
            "TEST_BEADS_LOG": str(self.log),
        }
        self.env.pop("BEADS_GIT_SYNC_ACTIVE", None)

    def run_hook(self, name, *args):
        return subprocess.run(
            ["/bin/sh", str(HELPER), name, *args],
            cwd=self.root,
            env=self.env,
            text=True,
            capture_output=True,
        )

    def calls(self):
        if not self.log.exists():
            return []
        return [json.loads(line) for line in self.log.read_text().splitlines()]

    def test_commit_preserves_native_callback_and_commits_task_changes(self):
        self.assertEqual(self.run_hook("pre-commit").returncode, 0)
        self.assertEqual(self.calls(), [["hooks", "run", "pre-commit"], ["dolt", "commit"]])

    def test_push_flushes_before_sync_and_preserves_remote_arguments(self):
        self.assertEqual(
            self.run_hook("pre-push", "origin", "https://example.test/repo").returncode, 0
        )
        self.assertEqual(
            self.calls(),
            [
                ["hooks", "run", "pre-push", "origin", "https://example.test/repo"],
                ["dolt", "commit"],
                ["dolt", "push"],
            ],
        )

    def test_sync_failure_blocks_push(self):
        self.env["TEST_BEADS_FAIL"] = "dolt push"
        self.assertEqual(self.run_hook("pre-push").returncode, 17)

    def test_commit_failure_does_not_push_uncommitted_tasks(self):
        self.env["TEST_BEADS_FAIL"] = "dolt commit"
        self.assertEqual(self.run_hook("pre-push").returncode, 17)
        self.assertNotIn(["dolt", "push"], self.calls())

    def test_pull_preserves_local_changes_first(self):
        self.assertEqual(self.run_hook("post-merge", "0").returncode, 0)
        self.assertEqual(
            self.calls(),
            [["hooks", "run", "post-merge", "0"], ["dolt", "commit"], ["dolt", "pull"]],
        )

    def test_rebase_refreshes_beads_but_amend_does_not(self):
        self.assertEqual(self.run_hook("post-rewrite", "amend").returncode, 0)
        self.assertEqual(self.calls(), [])
        self.assertEqual(self.run_hook("post-rewrite", "rebase").returncode, 0)
        self.assertEqual(self.calls(), [["dolt", "commit"], ["dolt", "pull"]])

    def test_native_hooks_preserve_argument_boundaries(self):
        self.assertEqual(
            self.run_hook("prepare-commit-msg", "message with spaces", "message").returncode, 0
        )
        self.assertEqual(
            self.calls(), [["hooks", "run", "prepare-commit-msg", "message with spaces", "message"]]
        )

    def test_nested_git_transport_does_not_sync_recursively(self):
        self.env["BEADS_GIT_SYNC_ACTIVE"] = "1"
        self.assertEqual(self.run_hook("pre-push").returncode, 0)
        self.assertEqual(self.calls(), [])

    def test_missing_cli_fails_with_actionable_message(self):
        self.env["PATH"] = str(self.root / "missing")
        result = self.run_hook("pre-push")
        self.assertEqual(result.returncode, 1)
        self.assertIn("requires bd on PATH", result.stderr)

    def test_quality_gate_failure_prevents_beads_and_git_work(self):
        scripts = self.root / "scripts/git-hooks"
        scripts.mkdir(parents=True)
        shutil.copyfile(HELPER, scripts / "beads.sh")
        for name in ("pre-commit", "pre-push"):
            with self.subTest(name=name):
                (scripts / f"{name}.sh").write_text("exit 23\n")
                result = subprocess.run(
                    ["/bin/sh", str(ROOT / ".husky" / name)],
                    cwd=self.root,
                    env={**self.env, "PATH": "/bin:" + str(self.bin)},
                    capture_output=True,
                )
                self.assertEqual(result.returncode, 23)
                self.assertEqual(self.calls(), [])


if __name__ == "__main__":
    unittest.main()
