# AEGIS — Adaptive Enforcement & Guardrail Intelligence System

[![CI](https://github.com/aegis-security/aegis/actions/workflows/ci.yml/badge.svg)](https://github.com/aegis-security/aegis/actions)
[![Red Team](https://github.com/aegis-security/aegis/actions/workflows/redteam.yml/badge.svg)](https://github.com/aegis-security/aegis/actions)
[![Supply Chain](https://github.com/aegis-security/aegis/actions/workflows/supply-chain.yml/badge.svg)](https://github.com/aegis-security/aegis/actions)

AEGIS is an open-source AI security monitoring and guardrail system that protects LLM/AI deployments from prompt injection, jailbreaks, RAG poisoning, data exfiltration, and supply-chain attacks.

## Features

| Feature | Description |
|---------|-------------|
| ⚡ **Data-plane enforcement** | P95 ≤10ms synchronous policy check on every AI request |
| 🛡️ **OPA/Rego policy-as-code** | Hierarchical policies (org→app→model→env), hot-reload, version history |
| 🧠 **Dual-model ensemble scanner** | Two independent LLM classifiers + disagreement escalation |
| 🕷️ **Pre-filter pipeline** | Rule-based <1ms detection for 15+ injection patterns |
| 📁 **RAG poisoning detection** | Statistical anomaly detection on retrieval score distributions |
| 🎯 **Automated red team** | 25-case adversarial test library, CI gate every 6h |
| 📊 **Full observability** | 5 Grafana dashboards, SLO tracking, distributed tracing |
| 🔐 **Zero-trust architecture** | SPIFFE/SPIRE mTLS, Vault secrets, distroless containers |
| 📋 **Compliance ready** | EU AI Act, ISO 42001, NIST AI RMF, OWASP LLM Top 10 |

## Architecture
```
Client → Envoy TLS → Gateway (JWT + OPA) → AI Backend
                          ↓ (async, SHA-256 only)
                     Kafka event bus
              ├─── Scanner (dual-model)
              ├─── RAG Monitor
              └─── ClickHouse (audit log)
```

## Quick Start
See [docs/getting-started.md](docs/getting-started.md)

## Documentation
- [Architecture Reference](docs/architecture.md)
- [EU AI Act Compliance](docs/compliance/eu-ai-act-mapping.md)
- [ISO 42001 Compliance](docs/compliance/iso-42001-mapping.md)
- [OWASP LLM Top 10 Coverage](docs/compliance/owasp-llm-coverage.md)
- [Threat Model](docs/compliance/threat-model.md)
- [ADR Archive](docs/adr/)
- [Runbooks](docs/runbooks/)

## Testing
```bash
make test          # unit tests (Go + Python)
make opa-test      # OPA policy tests
make redteam-dry   # red team dry run
make redteam-validate  # validate test case registry
```

## License
Apache 2.0 — see [LICENSE](LICENSE)

## Security
See [SECURITY.md](SECURITY.md) for vulnerability reporting.
