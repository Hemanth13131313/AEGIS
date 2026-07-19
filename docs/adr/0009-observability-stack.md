# ADR 0009: Observability Stack

## Status
Accepted

## Context
AEGIS requires unified metrics, traces, and logs for SLO tracking and security investigation. As a security gateway, identifying performance bottlenecks and correlating security events across microservices (Gateway, Scanner, Policy Engine, RAG Monitor) is critical.

## Decision
We will use Prometheus for metrics, Grafana for dashboards/SLOs, and OTel Collector for trace aggregation and tail sampling.

## Alternatives
- **Datadog**: Rejected due to vendor lock-in and high cost at scale.
- **Elastic Stack**: Deferred due to operational complexity for initial deployment.

## Consequences
- Requires self-hosting and managing Prometheus, Grafana, and OTel Collector.
- Tail sampling at the OTel Collector reduces storage costs by only keeping 100% of error traces and 10% of successful traces.
- All SLOs are queryable via Prometheus using PromQL, simplifying alerting and dashboarding.
