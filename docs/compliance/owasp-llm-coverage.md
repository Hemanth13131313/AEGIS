# AEGIS — OWASP LLM Top 10 Detection Coverage

| ID | Vulnerability | Detection Component | Test Cases | Coverage |
|----|---------------|--------------------|-----------|---------|
| LLM01 | Prompt Injection | PromptInjectionEvaluator (pre-filter) + ensemble | RT-LLM01-001 through 011 (11 cases) | ✅ High |
| LLM02 | Insecure Output Handling | Gateway response scanning (Phase 4, output path) | RT-LLM02-001 (1 case) | ⚠ Low |
| LLM03 | Training Data Poisoning | RAG Monitor (score anomaly detection) | RT-LLM03-001 (1 case) | ⚠ Medium |
| LLM04 | Model Denial of Service | Sanitizer token limit (8192) + RT-LLM04-001 | RT-LLM04-001 (1 case) | ✅ Medium |
| LLM05 | Supply Chain Vulnerabilities | supply_chain.rego + Trivy + SBOM attestation | RT-LLM05-001 (1 case) | ✅ High |
| LLM06 | Sensitive Info Disclosure | DataExfiltrationEvaluator + ensemble | RT-LLM06-001 through 004 (4 cases) | ✅ High |
| LLM07 | Insecure Plugin Design | Tool allowlist Rego policy | RT-LLM07-001 (0 cases) | ⚠ Low |
| LLM08 | Excessive Agency | ExcessiveAgencyEvaluator + tool_allowlist.rego | RT-LLM08-001 through 003 (3 cases) | ✅ Medium |
| LLM09 | Overreliance | Out of scope for AEGIS (UI layer concern) | — | ❌ N/A |
| LLM10 | Model Theft | Architecture probe detection via ensemble | RT-LLM10-001 (1 case) | ⚠ Low |

## Coverage Summary
- **High coverage**: LLM01, LLM05, LLM06
- **Medium coverage**: LLM03, LLM04, LLM08
- **Low coverage (improve in future phases)**: LLM02, LLM07, LLM10
- **Out of scope**: LLM09
