"""Validate the SchoolCRM k3s overlay (zarf/k8s/vps) for the shared owezzy-db node.

Run with python3 -m unittest discover -s scripts/tests -p test_k8s_vps_contract.py.
Needs the kustomize binary and PyYAML. No cluster, images or credentials are used.
"""

import json
import re
import shutil
import subprocess
import unittest
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[2]
VPS = ROOT / "zarf/k8s/vps"
# What reaches the cluster: platform/ and cloudflared/ by hand, apps/ and jobs/ through Flux.
OVERLAYS = ("platform", "apps", "jobs", "cloudflared")


def build_path(path):
    result = subprocess.run(
        ["kustomize", "build", str(path)],
        capture_output=True,
        text=True,
        check=True,
    )
    return [doc for doc in yaml.safe_load_all(result.stdout) if doc]


def build(overlay):
    return build_path(VPS / overlay)


def pod_specs(docs):
    for doc in docs:
        if doc["kind"] in ("Deployment", "Job"):
            yield doc["metadata"]["name"], doc["spec"]["template"]["spec"]


def containers(spec):
    return spec.get("initContainers", []) + spec["containers"]


def env_map(container):
    return {e["name"]: e for e in container.get("env", [])}


@unittest.skipUnless(shutil.which("kustomize"), "kustomize is not installed")
class K8sVPSContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.docs = {name: build(name) for name in OVERLAYS}
        cls.all = [doc for docs in cls.docs.values() for doc in docs]

    def test_everything_lives_in_the_schoolcrm_namespace(self):
        namespaces = [d for d in self.all if d["kind"] == "Namespace"]
        self.assertEqual([d["metadata"]["name"] for d in namespaces], ["schoolcrm"])
        labels = namespaces[0]["metadata"]["labels"]
        self.assertEqual(labels["pod-security.kubernetes.io/enforce"], "restricted")
        for doc in self.all:
            if doc["kind"] not in ("Namespace", "StorageClass"):
                self.assertEqual(doc["metadata"]["namespace"], "schoolcrm", doc["metadata"]["name"])

    def test_shared_database_is_referenced_not_managed(self):
        kinds = {d["kind"] for d in self.all}
        self.assertNotIn("StatefulSet", kinds)
        postgres = next(d for d in self.all if d["kind"] == "Service" and d["metadata"]["name"] == "postgres")
        self.assertEqual(postgres["spec"]["type"], "ExternalName")
        for name, spec in pod_specs(self.all):
            for c in spec["containers"]:
                env = env_map(c)
                for prefix in ("AUTH", "SCHOOLCRM"):
                    if f"{prefix}_DB_HOST" in env:
                        self.assertEqual(env[f"{prefix}_DB_HOST"]["value"], "postgres:5432")
                        self.assertEqual(env[f"{prefix}_DB_NAME"]["value"], "schoolcrm")
                        self.assertEqual(env[f"{prefix}_DB_USER"]["value"], "schoolcrm")
                        self.assertEqual(env[f"{prefix}_DB_DISABLE_TLS"]["value"], "false", name)
                        self.assertLessEqual(int(env[f"{prefix}_DB_MAX_OPEN_CONNS"]["value"]), 10)

    def test_nothing_is_exposed_on_the_host(self):
        for doc in self.all:
            if doc["kind"] == "Service":
                self.assertIn(doc["spec"].get("type", "ClusterIP"), ("ClusterIP", "ExternalName"))
        for name, spec in pod_specs(self.all):
            self.assertFalse(spec.get("hostNetwork", False), name)
            for c in containers(spec):
                for port in c.get("ports", []):
                    self.assertNotIn("hostPort", port, name)

    def test_runtime_is_unprivileged_read_only_and_bounded(self):
        for name, spec in pod_specs(self.all):
            self.assertTrue(spec["securityContext"]["runAsNonRoot"], name)
            self.assertFalse(spec.get("automountServiceAccountToken", True), name)
            for c in containers(spec):
                sc = c["securityContext"]
                self.assertFalse(sc["allowPrivilegeEscalation"], name)
                self.assertTrue(sc["readOnlyRootFilesystem"], name)
                self.assertEqual(sc["capabilities"]["drop"], ["ALL"], name)
                limits = c["resources"]["limits"]
                self.assertIn("memory", limits, name)
                self.assertIn("cpu", limits, name)

    def test_migrations_are_explicit_jobs_and_never_seed(self):
        api = dict(pod_specs(self.docs["apps"]))["schoolcrm"]
        self.assertEqual(api.get("initContainers", []), [])
        self.assertEqual([c["name"] for c in api["containers"]], ["schoolcrm"])
        jobs = {d["metadata"]["name"]: d for d in self.docs["jobs"] if d["kind"] == "Job"}
        self.assertEqual(jobs["schoolcrm-migrate"]["spec"]["template"]["spec"]["containers"][0]["command"], ["./admin", "migrate"])
        for job in jobs.values():
            self.assertEqual(job["spec"]["backoffLimit"], 0)
        for doc in self.all:
            text = yaml.safe_dump(doc)
            self.assertNotIn("migrate-seed", text)
            self.assertNotRegex(text, r"\bseed\b")

    def test_secrets_come_only_from_secret_objects(self):
        self.assertFalse([d for d in self.all if d["kind"] in ("Secret", "ConfigMap")])
        for name, spec in pod_specs(self.all):
            for c in containers(spec):
                for e in c.get("env", []):
                    if any(word in e["name"] for word in ("PASSWORD", "TOKEN", "API_KEY", "DATABASE_URL")):
                        self.assertIn("secretKeyRef", e.get("valueFrom", {}), f"{name}: {e['name']}")
                    self.assertNotIn("configMapKeyRef", e.get("valueFrom", {}), f"{name}: {e['name']}")

    def test_every_service_runs_the_release_version(self):
        version = json.loads((ROOT / "package.json").read_text())["version"]
        for name, spec in pod_specs(self.docs["apps"] + self.docs["jobs"]):
            for c in containers(spec):
                self.assertRegex(c["image"], rf"^ghcr\.io/owezzy-tech/schoolcrm-[a-z-]+:v{re.escape(version)}$", name)
                self.assertEqual(c["imagePullPolicy"], "IfNotPresent", name)

    def test_flux_deploys_only_inside_the_namespace_and_migrates_first(self):
        flux = build_path(ROOT / "zarf/k8s/flux")
        repo = next(d for d in flux if d["kind"] == "GitRepository")
        self.assertIn("semver", repo["spec"]["ref"])
        syncs = {d["metadata"]["name"]: d["spec"] for d in flux if d["kind"] == "Kustomization" and "sourceRef" in d.get("spec", {})}
        self.assertEqual(set(syncs), {"schoolcrm-migrations", "schoolcrm-apps"})
        for spec in syncs.values():
            self.assertEqual(spec["serviceAccountName"], "flux-deployer")
            self.assertTrue(spec["wait"])
        self.assertEqual(syncs["schoolcrm-migrations"]["path"], "./zarf/k8s/vps/jobs")
        self.assertTrue(syncs["schoolcrm-migrations"]["force"])
        self.assertEqual(syncs["schoolcrm-apps"]["dependsOn"], [{"name": "schoolcrm-migrations"}])
        binding = next(d for d in flux if d["kind"] == "RoleBinding" and d["metadata"]["name"] == "flux-deployer")
        self.assertEqual(binding["metadata"]["namespace"], "schoolcrm")
        self.assertFalse([d for d in flux if d["kind"] == "ClusterRoleBinding" and "flux-deployer" in yaml.safe_dump(d)])

    def test_user_files_survive_claim_deletion(self):
        storage = next(d for d in self.all if d["kind"] == "StorageClass")
        self.assertEqual(storage["reclaimPolicy"], "Retain")
        pvc = next(d for d in self.all if d["kind"] == "PersistentVolumeClaim")
        self.assertEqual(pvc["spec"]["storageClassName"], storage["metadata"]["name"])

    def test_quota_fits_the_namespace_workloads(self):
        quota = next(d for d in self.all if d["kind"] == "ResourceQuota")["spec"]["hard"]
        mib = lambda v: int(v[:-2]) * (1024 if v.endswith("Gi") else 1)
        limit = lambda spec: sum(mib(c["resources"]["limits"]["memory"]) for c in spec["containers"])
        services = sum(limit(spec) for d in self.all if d["kind"] == "Deployment" for spec in [d["spec"]["template"]["spec"]])
        # Flux starts all migration Jobs together during a release, alongside the services.
        jobs = sum(limit(d["spec"]["template"]["spec"]) for d in self.all if d["kind"] == "Job")
        self.assertLessEqual(services + jobs, mib(quota["limits.memory"]))


if __name__ == "__main__":
    unittest.main()
