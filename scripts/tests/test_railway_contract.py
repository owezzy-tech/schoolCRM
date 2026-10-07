"""Static deployment contracts; image and HTTP proof are separate release gates."""

import fnmatch
import tomllib
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]


class RailwayContractTests(unittest.TestCase):
    def test_service_build_and_health_configuration(self):
        for name in ("auth", "schoolcrm", "rag", "web-admin"):
            with self.subTest(service=name):
                config = tomllib.loads((ROOT / f"zarf/railway/{name}.toml").read_text())
                self.assertEqual(config["build"]["builder"], "DOCKERFILE")
                self.assertTrue((ROOT / config["build"]["dockerfilePath"]).is_file())
                self.assertEqual(
                    config["deploy"]["healthcheckPath"],
                    "/" if name == "web-admin" else "/v1/readiness",
                )
                self.assertEqual(config["deploy"]["restartPolicyMaxRetries"], 5)
                self.assertNotIn(
                    "--mount=type=cache",
                    (ROOT / config["build"]["dockerfilePath"]).read_text(),
                )

    def test_upload_excludes_secrets_and_local_state(self):
        patterns = (ROOT / ".railwayignore").read_text().splitlines()
        for path in (
            ".env",
            "api/services/RAG/.env",
            "api/services/RAG/.env.staging",
            "zarf/keys/development.pem",
            ".beads/dolt/data",
            "api/services/RAG/var/files/private.pdf",
            "zarf/compose/database-data/data/PG_VERSION",
        ):
            with self.subTest(path=path):
                self.assertTrue(
                    any(
                        fnmatch.fnmatch(path, p + "*" if p.endswith("/") else p)
                        for p in patterns
                    )
                )

    def test_image_context_excludes_environment_secrets(self):
        patterns = (ROOT / ".dockerignore").read_text().splitlines()
        for path in (".env", "api/services/RAG/.env", "api/services/RAG/.env.staging"):
            self.assertTrue(any(fnmatch.fnmatch(path, p) for p in patterns))

    def test_auth_default_has_no_development_keys(self):
        dockerfile = (ROOT / "zarf/docker/dockerfile.auth").read_text()
        runtime = dockerfile.split("FROM alpine:3.23 AS auth_runtime", 1)[1].split(
            "FROM auth_runtime AS auth_development", 1
        )[0]
        self.assertNotIn("COPY --from=build_auth --chown=auth:auth /service/zarf/keys", runtime)
        self.assertEqual(dockerfile.strip().splitlines()[-1], "FROM auth_runtime AS production")
        self.assertIn("--target auth_development", (ROOT / "makefile").read_text())

    def test_rag_install_is_locked_and_not_editable(self):
        dockerfile = (ROOT / "zarf/docker/dockerfile.rag").read_text()
        self.assertIn("uv sync --frozen --no-dev --no-editable", dockerfile)
        self.assertNotIn("pip install -e", dockerfile)
        self.assertIn("UV_PROJECT_ENVIRONMENT=/opt/venv", dockerfile)

    def test_frontend_routes_rag_before_general_api(self):
        caddy = (ROOT / "api/frontends/web-admin/Caddyfile").read_text()
        self.assertLess(caddy.index("handle /v1/rag/*"), caddy.index("handle /v1/*"))
        self.assertIn("{$RAG_API_UPSTREAM:rag:7000}", caddy)
        self.assertIn(":{$PORT:80}", caddy)
        self.assertNotIn("handle_path", caddy)


if __name__ == "__main__":
    unittest.main()
