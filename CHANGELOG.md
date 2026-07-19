# Changelog

All notable changes to AEGIS will be documented in this file.
Format: [Keep a Changelog](https://keepachangelog.com/) — Semantic Versioning.

## [Unreleased]

## [0.8.0] — 2025 (Phase 8: Compliance & Evidence)
### Added
- EU AI Act, ISO 42001, NIST AI RMF compliance mapping documents
- OWASP LLM Top 10 detection coverage matrix
- STRIDE threat model document
- Evidence bundle generator (ZIP with manifest + checksum)
- eu_ai_act.rego compliance policy with test coverage
- Getting started guide

## [0.7.0] — 2025 (Phase 7: Hardening)
### Added
- SPIFFE/SPIRE full integration (workloadapi.X509Source)
- mTLS on all gRPC connections
- Distroless containers (nonroot) with binary hardening flags
- Cosign keyless image signing + SBOM attestation
- Trivy CRITICAL CVE CI gate
- CodeQL SAST (Go/Python/JS)
- Kubernetes NetworkPolicy default-deny
- Pod Security Standards (restricted)
- Vault K8s ServiceAccount auth

## [0.6.0] — 2025 (Phase 6: Observability)
### Added
- 5 Grafana dashboards (overview, SLO, scanner, red team, policy)
- Prometheus alerting rules (9 rules across SLO/security/Kafka)
- OTel tail sampling
- Alertmanager integration
- Runbooks for key alerts

## [0.5.0] — 2025 (Phase 5: Red Team)
### Added
- 25-case adversarial test registry
- 5 test case generators
- Red team CLI with Rich output
- CI gate: scheduled every 6h

## [0.4.0] — 2025 (Phase 4: Detection Engine)
### Added
- Dual-model ensemble scanner (concurrent asyncio.gather)
- Structural isolation for meta-injection prevention
- 3 rule-based pre-filters (<1ms)
- RAG anomaly detection (statistical)

## [0.3.0] — 2025 (Phase 3: Event Bus)
### Added
- Kafka 5-topic topology
- events.proto schema
- Fire-and-forget gateway producer (SHA-256 hash only)
- Scanner + RAG Monitor Kafka consumers

## [0.2.0] — 2025 (Phase 2: Policy Engine)
### Added
- OPA/Rego policy engine with hierarchy resolver
- PostgreSQL + ClickHouse schemas
- /api/v1/check verdict endpoint

## [0.1.0] — 2025 (Phase 1: Identity & Ingress)
### Added
- JWT/JWKS authentication middleware
- Envoy TLS proxy configuration
- Vault secret bootstrap

## [0.0.1] — 2025 (Phase 0: Foundations)
### Added
- Monorepo structure
- 5 service skeletons
- CI/CD pipeline
- Helm charts
