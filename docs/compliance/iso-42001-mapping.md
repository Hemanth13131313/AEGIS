# AEGIS — ISO/IEC 42001:2023 Compliance Mapping

**Standard**: ISO/IEC 42001:2023 — Artificial Intelligence Management System (AIMS)

## Clause Coverage

| Clause | Requirement | AEGIS Implementation | Status |
|--------|-------------|---------------------|--------|
| 4.1 | Understanding the organization | Architecture.md, PRD.md define context | ✅ |
| 4.2 | Interested parties | Multi-tenant model (org/app/model scope) | ✅ |
| 5.1 | Leadership commitment | SECURITY.md, CONTRIBUTING.md, ADRs | ✅ |
| 5.2 | AI Policy | infra/policies/ Rego policies | ✅ |
| 6.1 | Risk assessment | Red team automation (Phase 5), threat model | ✅ |
| 6.2 | Objectives | SLOs defined and tracked (Phase 6) | ✅ |
| 7.1 | Resources | Helm resource limits, HPA, runbooks | ✅ |
| 7.5 | Documented information | All ADRs, compliance docs, runbooks | ✅ |
| 8.1 | Operational planning | deploy/ directory, CI/CD pipelines | ✅ |
| 8.4 | AI system lifecycle | Phase-by-phase development, versioned policies | ✅ |
| 9.1 | Monitoring/measurement | Prometheus, Grafana SLO dashboard | ✅ |
| 9.2 | Internal audit | Red team CI gate (every 6h), CodeQL SAST | ✅ |
| 10.1 | Nonconformity/corrective action | Alert runbooks, on-call procedures | ⚠ Partial |
| A.2.2 | AI risk assessment | OWASP LLM Top 10 + MITRE ATLAS coverage | ✅ |
| A.2.5 | AI transparency | Trace Explorer, policy_ref in responses | ✅ |
| A.2.6 | Accountability | Policy version history, actor field on all changes | ✅ |
| A.3.3 | Data for AI | RAG Monitor, input sanitization, no raw storage | ✅ |
| A.4.1 | System performance | P95 SLOs tracked continuously | ✅ |
| A.6.1 | AI incident response | Runbooks, Alertmanager, trace explorer | ⚠ Partial |
