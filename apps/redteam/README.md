# AEGIS Red Team Automation

Phase 5: Continuous adversarial validation of the AEGIS detection pipeline.

## Purpose
An automated red-team harness that tests the AEGIS detection pipeline with a versioned library of adversarial test cases. Validates that the scanner correctly blocks OWASP LLM Top 10 and MITRE ATLAS techniques, and that policy changes don't introduce regressions.

## Test Case Format
Test case IDs follow the format: `RT-{OWASP_ID}-{3-digit-sequence}` (e.g., RT-LLM01-001).

## Environment Variables
- `AEGIS_REDTEAM_TARGET_URL`: The AEGIS /scan URL (e.g., `http://localhost:8080`)
- `AEGIS_API_KEY`: Auth token if required.

## CLI Usage

```bash
# Run tests
uv run python cli.py run --target http://localhost:8080 --concurrency 5

# Dry run with specific file
uv run python cli.py run --dry-run --cases testcases/registry.json

# List cases
uv run python cli.py list-cases --owasp LLM01

# Validate test cases schema
uv run python cli.py validate --cases testcases/registry.json

# Generate mock cases
uv run python cli.py generate --type prompt_injection --count 5
```

## CI Integration
Exits with code 1 on failures, serving as a CI gate.
Outputs JSON for dashboards.

## Adding New Test Cases
Always bump the version and add a new record. Never mutate existing test cases to maintain historical integrity.
All block cases must pass, and safe cases must never be false-positived.
