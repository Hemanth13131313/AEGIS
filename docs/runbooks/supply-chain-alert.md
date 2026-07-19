# Runbook: Supply Chain Security Alert

## Alert Description
A new CRITICAL CVE has been detected in a deployed image, or an unauthorized image modification was detected.

## Immediate Actions
1. **Identify Affected Image:** Check the CI/CD logs, Trivy scan results, or Kubernetes admission controller logs to find the image tag and digest.
2. **Verify Exploitability:** Check if the CVE is exploitable in a distroless context (e.g., does it rely on a shell or package manager?).

## Remediation Steps
1. **Update Base Image / Dependencies:** Bump the base image tag or library dependency in `go.mod` / `pyproject.toml`.
2. **Rebuild & Attest:** Run the CI pipeline to rebuild the image, generate a new SBOM, and re-attest with Cosign.
3. **Redeploy:** Roll out the updated image to the cluster via a rolling update.

## Verification
1. Run a Trivy scan on the newly built image locally or via CI to confirm the CVE is resolved.
2. Ensure the Kubernetes admission controller allows the new image and that it successfully starts.
