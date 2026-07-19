# 0015 - Kyverno Admission Control

## Status
Accepted

## Context
We need to enforce supply chain security by ensuring only signed container images and those with valid SBOMs are deployed to the cluster.

## Decision
- We will use Kyverno ClusterPolicy to enforce Cosign signatures and CycloneDX attestations on all pods in the `aegis` namespace.
- We will also enforce security standards like disallowing privileged containers.

## Consequences
- Images that are unsigned or lack SBOM attestation are rejected at pod creation.
- Requires installing the Kyverno cluster addon.
- Drastically hardens the runtime security posture against malicious images.
