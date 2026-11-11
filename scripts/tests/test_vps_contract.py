"""Validate VPS isolation and startup contracts using Docker's Compose parser.

Run with python3 -m unittest discover -s scripts/tests -p test_vps_contract.py.
No daemon, image pulls, live credentials or database access are required.
"""

import fnmatch
import json
import os
import subprocess
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
COMPOSE = [
    "docker",
    "compose",
    "--env-file",
    os.devnull,
    "-f",
    str(ROOT / "zarf/vps/compose.yml"),
    "--profile",
    "maintenance",
    "config",
]


class VPSContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.environment = {
            key: value
            for key, value in os.environ.items()
            if not key.startswith(("COMPOSE_", "SCHOOLCRM_", "AUTH_"))
        }
        cls.environment.update(
            {
                "SCHOOLCRM_AUTH_IMAGE": "example.invalid/auth@sha256:" + "0" * 64,
                "SCHOOLCRM_API_IMAGE": "example.invalid/api@sha256:" + "1" * 64,
            }
        )
        result = subprocess.run(
            COMPOSE + ["--no-env-resolution", "--format", "json"],
            env=cls.environment,
            capture_output=True,
            text=True,
            check=True,
        )
        cls.config = json.loads(result.stdout)

    def test_shared_database_is_not_managed_or_recreated(self):
        self.assertEqual(set(self.config["services"]), {"auth", "api", "migrate"})
        self.assertNotIn("volumes", self.config)
        database = self.config["networks"]["database"]
        self.assertTrue(database["external"])
        self.assertEqual(database["name"], "schoolcrm-database")
        for name, service in self.config["services"].items():
            prefix = "AUTH" if name == "auth" else "SCHOOLCRM"
            env = service["environment"]
            self.assertEqual(env[f"{prefix}_DB_NAME"], "schoolcrm")
            self.assertEqual(env[f"{prefix}_DB_USER"], "schoolcrm")
            self.assertEqual(env[f"{prefix}_DB_DISABLE_TLS"], "false")
            self.assertLessEqual(int(env[f"{prefix}_DB_MAX_OPEN_CONNS"]), 10)

    def test_only_http_ports_are_published_on_host_loopback(self):
        expected = {"api": ("6300", 3000), "auth": ("6600", 6000)}
        for name, (published, target) in expected.items():
            ports = self.config["services"][name]["ports"]
            self.assertEqual(len(ports), 1)
            self.assertEqual(ports[0]["host_ip"], "127.0.0.1")
            self.assertEqual(
                (ports[0]["published"], ports[0]["target"]), (published, target)
            )
        self.assertNotIn("ports", self.config["services"]["migrate"])

    def test_runtime_cannot_grow_without_limits_or_write_its_image(self):
        for service in self.config["services"].values():
            self.assertEqual(service["user"], "1000:1000")
            self.assertTrue(service["read_only"])
            self.assertEqual(service["cap_drop"], ["ALL"])
            self.assertIn("no-new-privileges:true", service["security_opt"])
            self.assertLessEqual(int(service["mem_limit"]), 512 * 1024 * 1024)
            self.assertLessEqual(float(service["cpus"]), 1)
            self.assertEqual(service["pids_limit"], 128)
            self.assertEqual(service["logging"]["options"]["max-size"], "10m")

    def test_migrations_are_explicit_and_never_seed(self):
        migration = self.config["services"]["migrate"]
        self.assertEqual(migration["profiles"], ["maintenance"])
        self.assertEqual(migration["entrypoint"], ["./admin"])
        self.assertEqual(migration["command"], ["migrate"])
        self.assertEqual(migration["restart"], "no")
        api = self.config["services"]["api"]
        self.assertEqual(api["depends_on"]["auth"]["condition"], "service_healthy")

    def test_missing_images_fail_without_development_defaults(self):
        for key in ("SCHOOLCRM_AUTH_IMAGE", "SCHOOLCRM_API_IMAGE"):
            env = self.environment.copy()
            del env[key]
            result = subprocess.run(
                COMPOSE + ["--no-env-resolution", "--quiet"],
                env=env,
                capture_output=True,
                text=True,
            )
            self.assertNotEqual(result.returncode, 0)
            self.assertIn(key, result.stderr)

    def test_secrets_are_external_and_excluded_from_build_context(self):
        auth = self.config["services"]["auth"]
        keys = auth["volumes"][0]
        self.assertTrue(keys["read_only"])
        self.assertFalse(keys["bind"]["create_host_path"])
        for service in self.config["services"].values():
            self.assertEqual(len(service["env_file"]), 1)
            self.assertFalse(any("PASSWORD" in k for k in service["environment"]))
        patterns = (ROOT / ".dockerignore").read_text().splitlines()
        for path in (
            "zarf/vps/secrets/auth.env",
            "zarf/vps/backups/source.dump",
            "zarf/vps/deploy.env",
        ):
            self.assertTrue(
                any(
                    fnmatch.fnmatch(path, p + "*" if p.endswith("/") else p)
                    for p in patterns
                )
            )


if __name__ == "__main__":
    unittest.main()
