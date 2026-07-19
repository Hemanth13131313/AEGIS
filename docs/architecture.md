# AEGIS — Architecture Reference

## System Overview
AEGIS (Adaptive Enforcement & Guardrail Intelligence System) is a 4-plane AI security monitoring system providing synchronous policy enforcement and asynchronous behavioral analysis for AI/LLM deployments.

## Architecture Planes

### Data Plane (P95 ≤10ms)
- Envoy → Gateway → AI Backend
- Synchronous: JWT auth + input sanitization + OPA policy check
- Fire-and-forget: Kafka event emit (non-blocking)

### Control Plane
- Policy Engine: policy CRUD, OPA evaluation, hierarchy resolution
- Policy hot-reload via aegis.control.policy-reload Kafka topic

### Analysis Plane (Async)
- Scanner: Kafka consumer → dual-model ensemble → detection emit
- RAG Monitor: Kafka consumer → statistical anomaly detection
- Red Team: scheduled adversarial testing

### Observability Plane
- Metrics: Prometheus + Grafana (5 dashboards)
- Traces: OTel Collector with tail sampling → Jaeger/Tempo
- Logs: structlog (Python) / zap (Go) with request_id correlation
- Alerts: Alertmanager with SLO breach + security event rules

## Service Inventory
| Service | Language | Port | Role |
|---------|----------|------|------|
| gateway | Go | 8080 | Data-plane enforcement |
| policy-engine | Go | 8081/9090 | REST API + gRPC policy server |
| scanner | Python | 8000 | Async ML detection |
| rag-monitor | Python | 8001 | RAG anomaly detection |
| redteam | Python | CLI | Adversarial test runner |
| ui | TypeScript | 3001 | Web dashboard |

## Security Architecture
- **Identity**: SPIFFE/SPIRE (X.509 SVIDs, 1h rotation)
- **Transport**: mTLS everywhere (Envoy → Gateway → PolicyEngine gRPC)
- **Secrets**: HashiCorp Vault (K8s ServiceAccount JWT auth)
- **Policy**: OPA/Rego (hierarchical: org → app → model → env)
- **Containers**: gcr.io/distroless nonroot, Cosign-signed, SBOM-attested
- **Network**: Kubernetes NetworkPolicy default-deny + explicit allow-lists

## Data Privacy
- Raw payload content NEVER stored or logged
- SHA-256 hash only in event records
- PII: user_ref pseudonymized in session records
- Retention: events 7d, detections 30d, redteam_runs 90d

## Compliance
- EU AI Act 2024/1689: ✅ (see docs/compliance/eu-ai-act-mapping.md)
- ISO 42001:2023: ✅ (see docs/compliance/iso-42001-mapping.md)
- NIST AI RMF 1.0: ✅ (see docs/compliance/nist-ai-rmf-mapping.md)
- OWASP LLM Top 10: ✅ (see docs/compliance/owasp-llm-coverage.md)
