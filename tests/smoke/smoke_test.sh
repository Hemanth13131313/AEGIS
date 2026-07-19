#!/usr/bin/env bash
# AEGIS Smoke Test
# Usage: AEGIS_URL=http://localhost:8080 bash smoke_test.sh
set -euo pipefail

BASE="${AEGIS_URL:-http://localhost:8080}"
PASS=0
FAIL=0

check() {
  local name="$1" expected="$2" actual="$3"
  if [ "$actual" = "$expected" ]; then
    echo "✅ $name"
    PASS=$((PASS+1))
  else
    echo "❌ $name (expected: $expected, got: $actual)"
    FAIL=$((FAIL+1))
  fi
}

# Health checks
check "gateway health" "200" "$(curl -s -o /dev/null -w '%{http_code}' $BASE/health)"
check "policy-engine health" "200" "$(curl -s -o /dev/null -w '%{http_code}' ${POLICY_ENGINE_URL:-http://localhost:8081}/api/v1/health)"
check "scanner health" "200" "$(curl -s -o /dev/null -w '%{http_code}' ${SCANNER_URL:-http://localhost:8000}/health)"

# Auth check: missing token returns 401
check "auth missing token 401" "401" "$(curl -s -o /dev/null -w '%{http_code}' $BASE/v1/chat/completions)"

# Scanner scan: benign payload returns allow
VERDICT=$(curl -s -X POST ${SCANNER_URL:-http://localhost:8000}/scan \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"smoke-01","request_id":"smoke-req-01","payload":"What is the capital of France?","context":{}}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['action'])" 2>/dev/null || echo "error")
check "scanner benign returns allow" "allow" "$VERDICT"

echo ""
echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -eq 0 ] && exit 0 || exit 1
