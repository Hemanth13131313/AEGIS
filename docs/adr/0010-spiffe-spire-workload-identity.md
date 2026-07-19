# ADR 0010: SPIFFE/SPIRE for zero-trust workload identity and mTLS

## Status
Accepted (Phase 7 stub → Phase 7 full implementation)

## Context
Service-to-service auth via static certs is fragile; SPIRE provides automatic cert rotation and strongly binds identity to workloads.

## Decision
We will use SPIFFE workload identity with SPIRE. X.509 SVIDs will be rotated every 1 hour to limit the blast radius of any compromised credentials.

## Consequences
- SPIRE agent is required on each node.
- Slight startup latency for services to fetch SVIDs initially.
- Eliminates the need for long-lived service certificates and manual rotation.
