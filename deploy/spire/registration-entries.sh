#!/usr/bin/env bash
# Registers SPIFFE workload identity entries in SPIRE
# Run after SPIRE server is up: bash registration-entries.sh
set -euo pipefail

SERVER_SOCK="${SPIRE_SERVER_SOCK:-/tmp/spire-server/private/api.sock}"

register() {
  local spiffe_id="$1"
  local selector="$2"
  echo "Registering: $spiffe_id"
  spire-server entry create \
    -socketPath "$SERVER_SOCK" \
    -spiffeID "$spiffe_id" \
    -parentID "spiffe://aegis.cluster/spire/agent/join_token/$(cat /dev/urandom | tr -dc 'a-f0-9' | head -c 16)" \
    -selector "$selector" \
    -ttl 3600
}

register "spiffe://aegis.cluster/gateway" "unix:uid:1000"
register "spiffe://aegis.cluster/policy-engine" "unix:uid:1001"
register "spiffe://aegis.cluster/scanner" "unix:uid:1002"
register "spiffe://aegis.cluster/rag-monitor" "unix:uid:1003"
register "spiffe://aegis.cluster/redteam" "unix:uid:1004"

echo "All entries registered."
