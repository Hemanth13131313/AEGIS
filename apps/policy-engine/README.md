# AEGIS Policy Engine

Purpose: OPA/Rego policy engine with hierarchical scope resolution.

## Local Dev Setup

Required dependencies:
- PostgreSQL (PGX DSN: `AEGIS_POLICY_DB_DSN`)
- Redis (`AEGIS_POLICY_REDIS_ADDR`)
- OPA binary (optional for direct tests)

## Running Tests

- Go Unit/Integration: `go test ./...`
- Rego Policy Tests: `opa test rego/ -v`

## API Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/api/v1/policies` | Cursor-paginated list of policies |
| POST | `/api/v1/policies` | Create a new policy |
| PUT | `/api/v1/policies/{id}` | Update an existing policy (creates version) |
| GET | `/api/v1/policies/{id}` | Get policy by ID |
| GET | `/api/v1/policies/{id}/versions` | Get all versions for audit |
| POST | `/api/v1/check` | Synchronous policy evaluation |
| GET | `/api/v1/health` | Health check |

## Policy Hierarchy Diagram

Organization -> Application -> Model Endpoint -> Environment (Most specific)

## Milestones

- M2.1: Policy CRUD & Versioning
- M2.2: Tool Block & Evaluation Core
- M2.3: Redis Cache Failover

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `AEGIS_POLICY_LISTEN_ADDR` | `:8081` | REST API address |
| `AEGIS_POLICY_GRPC_ADDR` | `:9090` | gRPC address |
| `AEGIS_POLICY_DB_DSN` | - | Postgres connection string |
| `AEGIS_POLICY_REDIS_ADDR` | - | Redis cache address |
| `AEGIS_OTEL_EXPORTER_ENDPOINT` | - | OTel tracing endpoint |
