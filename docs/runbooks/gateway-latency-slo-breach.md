# Runbook: Gateway Latency P95 SLO Breach

## Alert: GatewayLatencyP95SLOBreach
**Severity**: Critical  
**SLO**: P95 ≤10ms

## Immediate Actions (< 5 minutes)
1. Check `aegis-overview` Grafana dashboard for latency trend
2. Run: `kubectl top pods -l app=aegis-gateway` (check CPU/memory)
3. Check policy engine health: `curl http://policy-engine:8081/api/v1/health`
4. Check Kafka lag: `kafka-consumer-groups.sh --bootstrap-server kafka:9092 --describe --group aegis-gateway`

## Diagnostic Steps
- High policy engine latency? → Check Redis cache hit rate (aegis:policy:* keys)
- High upstream latency? → Issue is downstream AI backend, not AEGIS
- CPU throttling? → Scale gateway replicas or increase resource limits

## Escalation
- > 15ms P95 for > 10 minutes: page on-call engineer
- > 50ms P95: incident declaration
