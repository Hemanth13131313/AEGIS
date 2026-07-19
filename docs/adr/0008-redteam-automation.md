# ADR 0008: Continuous Automated Red-Team as CI Gate

## Status
Accepted

## Context
Manual security testing of the AEGIS detection pipeline is infrequent and risks missing regressions when policy engines or LLM models are updated. Without continuous validation, an update might inadvertently disable prompt injection filters or allow excessive agency.

## Decision
We will build a versioned red-team test case library and an automated CLI runner. This runner will execute against the AEGIS `/scan` endpoint on every push to `main` and on a 6-hour schedule.
All test cases will be versioned and immutable—any change requires a version bump and new ID.

## Consequences
- Security regressions will be caught within 6 hours.
- Requires a live AEGIS instance in CI, or we perform a dry-run for simple pull request gates.
- Developers have immediate feedback on detection capability changes.
