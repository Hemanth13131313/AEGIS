# ADR 0012: Generate and attest SBOMs with Cosign for supply-chain transparency

## Status
Accepted

## Context
LLM supply chain attacks (LLM05) require verifiable provenance of all deployed artifacts to ensure that our runtime environment has not been tampered with and contains only authorized dependencies.

## Decision
We will use Syft to generate a CycloneDX SBOM per image. We will use Cosign keyless OIDC attestation on the main branch to sign and attest these SBOMs.

## Consequences
- Adds approximately 2 minutes to the CI pipeline.
- Provides a full, verifiable dependency graph essential for rapid incident response and vulnerability management.
