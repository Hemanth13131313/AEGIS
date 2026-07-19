# AEGIS Gateway

Cloud-native runtime security gateway for AI systems (Phase 1).

## Purpose
- JWT/OIDC Authentication middleware
- Transparent reverse proxy to AI backends
- Secret-zero bootstrap stub via Vault
- Observability via OpenTelemetry tracing

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `AEGIS_GATEWAY_LISTEN_ADDR` | Host/Port for Gateway to listen on | `:8080` |
| `AEGIS_GATEWAY_UPSTREAM_URL` | URL of the AI backend / model service | `http://localhost:5000` |
| `AEGIS_GATEWAY_POLICY_ENGINE_ADDR` | Address of the gRPC Policy Engine | `localhost:9090` |
| `AEGIS_GATEWAY_FAIL_MODE` | Behavior when policy engine fails (`open` or `closed`) | `closed` |
| `AEGIS_GATEWAY_JWKS_URL` | URL to fetch JWKS for token validation | (empty) |
| `AEGIS_GATEWAY_AUTH_SKIP_VERIFY` | Skip auth verification (dev only) | `false` |
| `AEGIS_OTEL_EXPORTER_ENDPOINT` | OTLP gRPC endpoint for tracing | (empty) |
| `AEGIS_SANITIZE_MAX_PAYLOAD_BYTES` | Maximum allowed payload size | `1048576` (1MB) |

## Quick Start

1. Start Gateway in dev mode:
```bash
AEGIS_GATEWAY_AUTH_SKIP_VERIFY=true go run ./cmd/main.go
```

2. Test Health Endpoint:
```bash
curl http://localhost:8080/health
```

## Milestone Status
- **M1.1**: Local Gateway Stub — Done
- **M1.2**: Identity & Envoy Stub — Done
- **M1.3**: Telemetry Traces — Done
