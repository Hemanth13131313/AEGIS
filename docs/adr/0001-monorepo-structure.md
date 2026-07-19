# 0001: Use monorepo structure over polyrepo

**Status:** Accepted

## Context
We are beginning the AEGIS project, building multiple services including a gateway (Go), scanner (Python), control-plane (Go), and UI (TS/React). Managing multiple repositories adds overhead to code sharing (SDKs, protobufs), CI/CD setups, and cross-service refactoring. We have a small core team working towards an MVP.

## Decision
We will use a monorepo structure, separating the codebase into `apps/` (services/UI) and `packages/` (shared libraries, protos).

## Consequences
- **Pros:** Easier cross-service refactoring, single CI pipeline, unified developer experience, simpler dependency management across shared libraries.
- **Cons:** The repository size will grow over time, potentially requiring more advanced tooling (e.g., Bazel or advanced NX setups) in the future to keep build/test times low. Initial CI setup requires filtering to only test what changed to be efficient.
