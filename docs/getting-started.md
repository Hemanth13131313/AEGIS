# AEGIS — Getting Started

## Prerequisites
- Docker + Docker Compose
- Go 1.22+
- Python 3.12+ with uv
- Node 20+
- OPA CLI (for policy testing)
- kubectl (for K8s deployment)

## Local Development (5 minutes)

1. **Start infrastructure:**
   ```bash
   cp .env.dev.example .env.dev
   docker compose -f docker-compose.dev.yml up -d
   ```

2. **Setup Vault:**
   ```bash
   bash infra/vault/setup-dev.sh
   ```

3. **Init Kafka topics:**
   ```bash
   docker compose -f docker-compose.tools.yml up kafka-init
   ```

4. **Run tests:**
   ```bash
   make test
   make opa-test
   ```

5. **Run red team (dry run):**
   ```bash
   make redteam-dry
   ```

6. **Start UI:**
   ```bash
   cd apps/ui && npm install && npm run dev
   ```

## Verify Installation
- Kafka UI: http://localhost:9080
- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9090
- Keycloak: http://localhost:8080

## Next Steps
- Read the [Architecture Reference](architecture.md)
- Review the [Policy Hierarchy](../infra/policies/)
- Run a real red team: `AEGIS_REDTEAM_TARGET=http://localhost:8080 make redteam`
