# AEGIS Security Policy

## Reporting Vulnerabilities
Do NOT open a public GitHub issue for security vulnerabilities.

Report to: security@aegis-security.io
PGP Key: [PLACEHOLDER - add before publishing]
Response time: ≤48 hours acknowledgement, ≤7 days initial assessment

## Supported Versions
| Version | Security Support |
|---------|------------------|
| 0.x.x (current) | ✅ Active |

## Scope
- AEGIS gateway, policy engine, scanner, RAG monitor, red team harness
- Infrastructure configs (Helm charts, Kubernetes manifests)
- Rego policies

## Out of Scope
- Third-party AI models behind AEGIS
- Kafka/PostgreSQL/Redis infrastructure (report to vendors)

## Security Controls
- All inter-service communication: mTLS via SPIFFE/SPIRE
- All secrets: Vault (no secrets in code or environment files)
- All images: distroless, Cosign-signed, SBOM-attested
- All policies: OPA/Rego with mandatory test coverage
- All releases: Trivy CRITICAL CVE gate + CodeQL SAST

## Vulnerability Disclosure Timeline
1. Day 0: Report received, acknowledged within 48h
2. Day 7: Initial severity assessment shared with reporter
3. Day 30: Patch developed and tested
4. Day 45: Coordinated disclosure (patch released + CVE published)

## Bug Bounty
No formal bug bounty program at this time. Recognition in CHANGELOG.md and SECURITY.md for verified reports.
