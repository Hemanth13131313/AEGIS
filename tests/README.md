# Tests

This directory structures the test suite across three layers for the AEGIS project:

- **unit/**: Per-service unit tests usually live alongside the source code (`apps/<service>/tests/`). This folder can hold project-wide unit test utilities.
- **integration/**: Cross-service integration tests (e.g., Gateway ↔ PolicyEngine, Gateway ↔ Message Bus ↔ Scanner).
- **e2e/**: Full golden-path End-to-End tests using an ephemeral Kind (Kubernetes-in-Docker) cluster. These are intended to be run in CI.
