# ADR 0004: OPA/Rego Policy Engine

## Title
Use OPA/Rego as the policy-as-code engine

## Status
Accepted

## Context
AEGIS requires a robust mechanism for enforcing rules and policies over AI interactions (input sanitization, prompt structure, tool-use allowlisting). We need a declarative, auditable, testable policy engine that can be hot-reloaded dynamically across tenants in a multi-tenant gateway.

## Decision
We will use Open Policy Agent (OPA) with Rego embedded inside the `policy-engine` service. Policies will be distributed to the gateway via HTTP (for Phase 2) and gRPC (Phase 7).

## Consequences
- **Learning curve:** Teams writing policies must learn Rego.
- **Strong testability:** We can write and execute unit tests for policies (`opa test`).
- **Policy-as-code discipline:** All rules are explicit and auditable; no inline if-statements for enforcement in the Gateway.

## Rules
- All Go code must have `context.Context` as the first param.
- Error wrapping with `fmt.Errorf("%w", err)`.
- No global mutable state.
- Rego policies must `import rego.v1`.
- Use `contains` instead of set comprehensions for deny rules.
- Every Rego package must include a corresponding `_test.rego` file.
- No hard-coded secrets.
- Use standard AEGIS error shape everywhere.
- Use cursor-based pagination on all list endpoints (never offset).
