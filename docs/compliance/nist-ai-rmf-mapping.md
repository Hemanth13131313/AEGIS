# AEGIS — NIST AI Risk Management Framework 1.0 Mapping

## Core Functions

| Function | Category | AEGIS Implementation | Status |
|----------|----------|---------------------|--------|
| GOVERN | 1.1 Policies | infra/policies/ Rego, OPA engine | ✅ |
| GOVERN | 1.7 Processes documented | ADRs 0001-0012, runbooks | ✅ |
| MAP | 1.1 Risk identification | OWASP LLM Top 10 + ATLAS mapping in scanner | ✅ |
| MAP | 2.3 Impact categories | owasp_category + atlas_technique on every detection | ✅ |
| MEASURE | 1.1 AI risk measurement | Red team pass rate, FP rate, SLO metrics | ✅ |
| MEASURE | 2.5 Explainability | policy_ref, reason fields on every decision | ✅ |
| MEASURE | 2.6 Robustness | Dual-model ensemble, pre-filters, disagreement escalation | ✅ |
| MANAGE | 1.1 Risk prioritization | Severity-tiered alerts (critical/high/medium/low) | ✅ |
| MANAGE | 2.2 Incident response | Alertmanager + runbooks | ⚠ Partial |
| MANAGE | 3.1 Risk treatment | block/redact/tag actions with policy_ref | ✅ |
