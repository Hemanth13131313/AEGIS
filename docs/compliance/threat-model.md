# AEGIS — Threat Model

**Method**: STRIDE  
**Scope**: AEGIS AI Security Monitoring System  
**Date**: 2025

## System Components in Scope
- Envoy TLS proxy
- AEGIS Gateway
- Policy Engine
- Scanner (dual-model ensemble)
- RAG Monitor
- Kafka event bus
- ClickHouse audit log

## STRIDE Analysis

### Spoofing
| Threat | Component | Mitigation |
|--------|-----------|------------|
| Attacker spoofs AEGIS gateway | Policy Engine | SPIFFE mTLS SVID validation |
| JWT token forged | Gateway | JWKS signature verification, RS256 only |
| Kafka producer spoofed | Kafka | SASL/SSL (Phase 9) |

### Tampering
| Threat | Component | Mitigation |
|--------|-----------|------------|
| Policy modified in transit | Policy Engine API | mTLS, audit log (policy_versions table) |
| Event log tampered | ClickHouse | Append-only MergeTree engine |
| Container image modified | All services | Cosign attestation, digest pinning |

### Repudiation
| Threat | Component | Mitigation |
|--------|-----------|------------|
| Policy change not audited | Policy Engine | policy_versions table with actor field |
| Detection event missing | Scanner | Kafka at-least-once + ClickHouse |

### Information Disclosure
| Threat | Component | Mitigation |
|--------|-----------|------------|
| Payload content in logs | All services | SHA-256 hash only, no raw content |
| Secrets in container images | All services | Distroless + Vault, no env file in image |
| Model exfiltration via probing | Scanner | DataExfiltrationEvaluator blocks probes |

### Denial of Service
| Threat | Component | Mitigation |
|--------|-----------|------------|
| Token flooding | Gateway/Sanitizer | 8192 token limit enforced by Rego + sanitizer |
| Scanner overload | Scanner | Kafka backpressure + concurrency limits |
| Policy engine overload | Policy Engine | Redis LKG cache, fail-open/closed configurable |

### Elevation of Privilege
| Threat | Component | Mitigation |
|--------|-----------|------------|
| Container breakout | All services | Distroless nonroot (uid 65532), read-only FS |
| K8s privilege escalation | All services | Pod Security Standards (restricted) + NetworkPolicy |
| Tool invocation without auth | Gateway | Tool allowlist Rego policy |
