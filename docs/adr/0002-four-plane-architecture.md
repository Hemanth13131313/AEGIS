# ADR 0002: Adopt four-plane architecture (data/control/analysis/observability)

**Status:** Accepted

## Context
AEGIS needs to serve as a runtime security gateway for AI systems. We need to evaluate fast, deterministic rules (like blocklists or rate limits) without slowing down the hot path of LLM requests. However, we also need to perform deep ML-based scanning (like injection detection) which can be latency-heavy. 
We need a way to separate latency-sensitive enforcement from deep ML scanning.

## Decision
We will adopt a four-plane architecture for the system:
1. **Data Plane:** Synchronous traffic enforcement (Envoy + Go Gateway).
2. **Control Plane:** Policy management and synchronous evaluation (OPA/Rego).
3. **Analysis Plane:** Asynchronous deep scanning via a message bus (Kafka/NATS + Python Scanner/RAGMon).
4. **Observability Plane:** Telemetry, traces, and metrics (ClickHouse, OpenTelemetry).

## Consequences
- **Lower p95 latency on the hot path:** By doing deep scanning asynchronously or out-of-band for some traffic, we can keep the synchronous path fast.
- **Resilience:** Scanner failures don't block traffic (depending on policy).
- **Complexity:** Requires careful fail-mode configuration (fail-open vs fail-closed) and handling eventual consistency for detection events.
