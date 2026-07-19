# AEGIS — EU AI Act Compliance Mapping

**Regulation**: Regulation (EU) 2024/1689 (AI Act)
**Date**: 2025
**Status**: Mapped to Phase 8 implementation

## Risk Classification
AEGIS is classified as a **General Purpose AI Safety Monitoring System** — not itself a high-risk AI system, but designed to make high-risk AI deployments compliant.

## Article Coverage

| Article | Requirement | AEGIS Implementation | Status |
|---------|-------------|---------------------|--------|
| Art. 9 | Risk Management System | OPA/Rego policy engine with hierarchical scope, quarterly red-team runs | ✅ Implemented |
| Art. 10 | Data Governance | RAG Monitor detects poisoning; no raw content stored (SHA-256 hash only) | ✅ Implemented |
| Art. 11 | Technical Documentation | ADRs 0001-0012, architecture.md, this document | ✅ Implemented |
| Art. 12 | Record Keeping | ClickHouse immutable event log, policy version history (immutable) | ✅ Implemented |
| Art. 13 | Transparency | Trace Explorer shows full decision chain; policy_ref in every block response | ✅ Implemented |
| Art. 14 | Human Oversight | Ensemble disagreement → tag (human review queue); block requires two-model agreement | ✅ Implemented |
| Art. 15 | Accuracy, Robustness, Cybersecurity | Dual-model ensemble; red-team automation; SPIFFE mTLS; distroless containers | ✅ Implemented |
| Art. 17 | Quality Management System | CI/CD with CodeQL, Trivy, OPA tests, red-team CI gate | ✅ Implemented |
| Art. 26 | Obligations for Deployers | Policy hierarchy allows org-level custom policies; deployment guide in docs/ | ⚠ Partial |
| Art. 72 | Market Surveillance | Compliance page + evidence bundle export | ⚠ Partial |

## Evidence Artifacts
- Policy version history: `apps/policy-engine/migrations/001_initial_schema.sql` (policy_versions table)
- Red team results: `apps/redteam/testcases/registry.json` + ClickHouse aegis.redteam_runs
- SBOM attestations: Generated per-release by `supply-chain.yml`
- ADR archive: `docs/adr/` (0001-0012)
- SLO records: Grafana `aegis-slo.json` dashboard + Prometheus time series
