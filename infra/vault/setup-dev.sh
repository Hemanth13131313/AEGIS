#!/usr/bin/env bash
# Sets up Vault for AEGIS local development
# Run after: docker compose up vault
set -euo pipefail

VAULT_ADDR="${VAULT_ADDR:-http://localhost:8200}"
VAULT_TOKEN="${VAULT_DEV_ROOT_TOKEN_ID:-dev-root-token}"

export VAULT_ADDR VAULT_TOKEN

echo "Enabling KV v2 secrets engine..."
vault secrets enable -path=secret kv-v2 2>/dev/null || echo "Already enabled"

echo "Writing placeholder secrets..."
vault kv put secret/aegis/gateway/config \
  jwks_url="http://keycloak:8080/realms/aegis/protocol/openid-connect/certs" \
  kafka_password="changeme_dev_only"

vault kv put secret/aegis/policy-engine/config \
  db_password="changeme_dev_only" \
  redis_password="changeme_dev_only"

vault kv put secret/aegis/scanner/config \
  model_api_key="" \
  kafka_password="changeme_dev_only"

echo "Writing AEGIS policy..."
vault policy write aegis infra/vault/aegis-policy.hcl

echo "Vault dev setup complete."
echo "Note: All secrets above are dev placeholders. Replace in production with real values via K8s auth."
