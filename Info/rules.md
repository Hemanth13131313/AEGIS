# AEGIS — Engineering Rules

**Related docs:** `prd.md`, `architecture.md`, `phases.md`, `design.md`, `memory.md`
**Audience:** human engineers and AI coding agents (e.g., Google Antigravity) implementing AEGIS.

---

## 1. Coding Standards

### General
- Every service is independently buildable and testable; no hidden cross-service imports outside `packages/`.
- Prefer explicit over clever. Optimize for the next engineer (human or AI) reading the code six months from now.
- No commented-out code committed to the main branch. Delete it; Git history is the backup.
- Public functions/exported symbols require a doc comment describing purpose, parameters, return, and error conditions.

### Go (Gateway, Policy Engine)
- Follow standard `gofmt`/`goimports` formatting; enforced in CI, not manually.
- Errors are values: always check and wrap with context (`fmt.Errorf("...: %w", err)`), never silently discard (`_ = err` requires an inline comment justifying it).
- No global mutable state except explicitly documented, mutex-guarded caches (e.g., local policy cache).
- Context (`context.Context`) is the first parameter of any function that does I/O or can be canceled/timed-out.

### Python (Scanner, RAG Monitor, Red-team harness)
- Follow PEP 8, enforced via `ruff`/`black` in CI.
- Type hints required on all public function signatures (`mypy` or `pyright` checked in CI).
- No bare `except:`; always catch specific exceptions and log with context.
- Async I/O (`asyncio`/FastAPI) for all network-bound service code; no blocking calls in async request handlers.

### TypeScript/React (UI)
- Strict mode enabled (`"strict": true` in `tsconfig.json`). No `any` without an inline justification comment.
- Function components + hooks only; no class components.
- No inline business logic in JSX — extract to hooks/`lib/` functions.

### Rego (Policy Engine)
- Every policy package includes accompanying `_test.rego` tests (see §Testing Standards).
- No policy rule may reference an undocumented input field — all `input.*` fields consumed must be documented in the policy's header comment.

## 2. Naming Conventions

| Element | Convention | Example |
|---|---|---|
| Go packages | short, lowercase, no underscores | `sanitize`, `policyclient` |
| Go exported types/functions | PascalCase | `PolicyClient`, `ValidateToken()` |
| Python modules/functions | snake_case | `ensemble_scanner.py`, `run_scan()` |
| Python classes | PascalCase | `EnsembleVerdict` |
| TypeScript components | PascalCase, filename matches component | `SessionTraceExplorer.tsx` |
| TypeScript hooks | camelCase, `use` prefix | `useSessionTrace()` |
| REST endpoints | kebab-case, plural nouns | `/api/v1/policy-versions` |
| Kafka/NATS topics | dot-namespaced | `aegis.events.request`, `aegis.detections.verdict` |
| Environment variables | `AEGIS_<COMPONENT>_<SETTING>`, upper snake case | `AEGIS_GATEWAY_FAIL_MODE` |
| Database tables | snake_case, plural | `policy_versions`, `redteam_cases` |
| Rego packages | `aegis.policy.<scope>` | `aegis.policy.tool_allowlist` |

## 3. Folder Conventions

- Monorepo layout as defined in `architecture.md §3`; do not introduce new top-level directories without updating that document.
- Each `apps/<service>/` is self-contained: its own dependency manifest, its own tests directory, its own Dockerfile.
- Shared code only lives in `packages/`; if two services need the same logic, extract to a package rather than copy-paste.
- Generated code (protobuf stubs, OpenAPI clients) lives in a clearly marked `generated/` subfolder and is excluded from manual edits and from strict lint rules.

## 4. Component Rules (Frontend)

- One component per file; filename == component name.
- Components over ~200 lines should be decomposed into subcomponents or extracted hooks.
- No component fetches data directly with `fetch`/`axios` inline — data access goes through a hook backed by React Query.
- Presentational components (pure UI) are kept separate from container components (data-fetching/orchestration) where a component exceeds trivial complexity.
- All interactive components must be keyboard-navigable and expose correct ARIA roles (see §Accessibility).

## 5. Accessibility Requirements

- WCAG 2.1 AA is the minimum bar for the SOC dashboard UI.
- All interactive elements reachable via keyboard (tab order logical, no keyboard traps).
- Color is never the sole means of conveying detection severity — always paired with an icon/label (see `design.md` for the severity system).
- Minimum contrast ratio 4.5:1 for body text, 3:1 for large text/UI components.
- All data visualizations (charts, trace timelines) have an accessible text/table alternative.
- Form inputs have associated `<label>`s; error states are announced via `aria-live` regions.

## 6. Performance Rules

- Gateway synchronous request path budget: **≤10ms** added latency for non-scanner checks (auth, sanitization, structural guardrails) — see `architecture.md §18`.
- Any new synchronous check added to the hot path must include a benchmark in its PR showing it stays within budget.
- No N+1 query patterns in the Policy Engine or reporting API; use batched loads/joins.
- Frontend: route-based code splitting; no single JS bundle chunk over 250KB gzipped without justification.
- All list/table views in the UI must paginate or virtualize beyond 200 rows.

## 7. Security Rules

- **Never** log raw prompt/response bodies at INFO level or below; body logging requires an explicit, tenant-scoped debug flag and is redacted by default (see `architecture.md §14`).
- **Never** hard-code secrets, API keys, or credentials in source, config files, or Helm values — use the secret-zero pattern (`architecture.md §8, §20`).
- All new inter-service communication must use mTLS; plaintext internal HTTP is never acceptable, including in local dev (use self-signed certs locally).
- Any new dependency is checked for known CVEs (Trivy/Grype) before merge; a new dependency introducing a high/critical CVE without a documented mitigation blocks the PR.
- User-supplied input is never interpolated directly into the scanner's own system prompt without structural delimiting (mitigates meta-injection against the scanner itself — see `architecture.md §16`).
- All container images are distroless and signed (Cosign) before deployment; unsigned images must be rejected by the admission controller in every environment including staging.

## 8. API Standards

- REST APIs are versioned via URI path (`/api/v1/...`); breaking changes require a new version, not an in-place change.
- Every error response follows a consistent shape:
  ```json
  {
    "error": {
      "code": "POLICY_BLOCKED",
      "category": "policy" ,
      "message": "Request blocked by policy aegis.policy.tool_allowlist (rule: disallow-shell-exec)",
      "policy_ref": "org/app/model/env policy id + version",
      "request_id": "uuid"
    }
  }
  ```
- `category` is one of: `policy` (AEGIS blocked it), `upstream` (AI provider error), `internal` (AEGIS fault) — per `architecture.md §15`. Clients must be able to branch on this field.
- All list endpoints support cursor-based pagination (`?cursor=...&limit=...`), never offset-based, to remain stable under concurrent writes.
- gRPC services use protobuf with backward-compatible field evolution only (no field renumbering/reuse).

## 9. Database Conventions

- Every table has `id` (UUID, primary key), `created_at`, and `updated_at` columns as a baseline.
- Foreign keys are always indexed.
- Migrations are forward-only, versioned, and reviewed like code (no manual schema edits against any environment).
- PostgreSQL is the source of truth for policy/tenancy; ClickHouse is append-only for trace/detection/red-team data — no service treats ClickHouse as mutable source-of-truth state.
- PII/sensitive fields (e.g., raw prompt text in a trace) must be tagged in the schema/migration comments so redaction tooling and access control can target them explicitly.

## 10. Git Workflow

- **Trunk-based development** with short-lived feature branches (`feature/<phase>-<short-desc>`, `fix/<short-desc>`).
- Direct pushes to `main` are disabled; all changes go through pull request with at least one required review.
- CI (lint, unit tests, policy tests, security scan) must pass before merge is allowed.
- Squash-merge into `main` to keep history readable; branch is deleted after merge.
- Release branches (`release/vX.Y`) are cut for tagged releases; hotfixes are cherry-picked back.

## 11. Commit Message Format

Conventional Commits, enforced via commit-lint in CI:

```
<type>(<scope>): <short summary>

[optional body]

[optional footer(s)]
```

- `type` ∈ `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `chore`, `security`.
- `scope` is the affected component (e.g., `gateway`, `scanner`, `policy-engine`, `ui`, `redteam`).
- Example: `feat(scanner): add dual-model disagreement escalation path`
- Example: `security(gateway): reject unsigned images at admission`

## 12. Testing Standards

| Test type | Requirement |
|---|---|
| Unit tests | Required for all non-trivial functions; minimum 80% line coverage per service, enforced in CI |
| Integration tests | Required for any change touching cross-service contracts (Gateway↔Policy Engine, Gateway↔Bus↔Scanner) |
| Policy tests (`opa test`) | Required for every Rego policy change; both allow and deny cases must be covered |
| Contract tests | Protobuf/OpenAPI schemas validated against consumers in CI to prevent breaking changes |
| Load/latency tests | Required before merging any change to the synchronous Gateway path; must stay within the budget in `architecture.md §18` |
| Red-team regression | New red-team test cases added for any newly discovered detection gap; run nightly against `main` |
| E2E tests | Golden-path scenarios (discovery → policy apply → detection → trace replay) run against an ephemeral k8s/kind cluster in CI |

## 13. Documentation Standards

- Every service has a `README.md` covering: purpose, local dev setup, how to run tests, how to deploy.
- Public APIs documented via OpenAPI (REST) / `.proto` comments (gRPC), generated into `docs/api/` — not hand-maintained separately.
- Architectural decisions with lasting consequence are recorded as ADRs (`docs/adr/NNNN-title.md`), referencing `architecture.md` sections they affect.
- `memory.md` must be updated whenever a decision changes something it documents (see `memory.md` itself for the update protocol).

## 14. Code Review Checklist

- [ ] Does this change stay within its service boundary (no unauthorized cross-service coupling)?
- [ ] Are new synchronous hot-path additions within the latency budget (with benchmark evidence)?
- [ ] Are secrets/credentials absent from the diff (no accidental commits)?
- [ ] Do error responses follow the standard shape and correct `category`?
- [ ] Are new Rego rules accompanied by `opa test` cases (allow + deny)?
- [ ] Is logging redaction-safe (no raw prompt/response bodies at default log level)?
- [ ] Are new dependencies scanned and free of unmitigated high/critical CVEs?
- [ ] Is test coverage adequate per §12?
- [ ] Does this change require an update to `architecture.md`, `prd.md`, or `memory.md`? If so, is that update included?
- [ ] Are new UI components accessible (keyboard nav, ARIA, contrast)?

## 15. Things the AI Must NEVER Do

1. **Never** invent or assume an API contract for OpenAI/Anthropic/Azure OpenAI/vLLM/TGI — verify against the actual provider documentation/spec before implementing an adapter.
2. **Never** log or persist raw prompt/response content outside the explicitly designed, access-controlled trace store, and never at default log verbosity.
3. **Never** hard-code secrets, tokens, or connection strings anywhere in source or config committed to Git.
4. **Never** bypass the policy engine for "just this one" enforcement decision — all enforcement decisions flow through the policy engine, even hard-coded defaults must be expressed as Rego, not inline `if` statements in the Gateway.
5. **Never** silently downgrade a `fail-closed` security path (auth/identity) to `fail-open` to "fix" an availability issue — that trade-off requires an explicit policy/config change, documented and reviewed.
6. **Never** ship a container image that is unsigned or has unresolved high/critical CVEs to any environment beyond local dev.
7. **Never** remove or weaken an existing red-team regression test to make a build pass — fix the underlying detection or explicitly and visibly mark the test as a known accepted risk with sign-off.
8. **Never** add a new top-level architectural component without first updating `architecture.md` (documentation precedes implementation for structural changes).
9. **Never** concatenate untrusted user input directly into the scanner's own instruction/system context without structural separation.
10. **Never** assume a numeric target (latency, coverage, throughput) beyond what's documented in `prd.md`/`architecture.md` — flag it as an open question instead of guessing silently.

## 16. Things the AI Should ALWAYS Do

1. **Always** check `memory.md` for prior decisions before proposing a new architectural approach.
2. **Always** map new detection logic to an OWASP LLM Top 10 category and/or MITRE ATLAS technique explicitly.
3. **Always** write both the allow-path and deny-path test for any new policy rule.
4. **Always** propagate a `request_id`/correlation ID through every hop for traceability.
5. **Always** prefer extending the existing pluggable detector/adapter interfaces (`architecture.md §22`) over one-off special cases.
6. **Always** document assumptions explicitly (in code comments and/or `memory.md`) when information is genuinely unavailable, per the standard used throughout this documentation set.
7. **Always** run policy changes through `opa test` locally before pushing.
8. **Always** consider streaming (SSE) behavior explicitly when touching the enforcement pipeline — do not assume request/response is always fully buffered.
9. **Always** keep this documentation set (`prd.md`, `architecture.md`, `rules.md`, `phases.md`, `design.md`, `memory.md`) internally consistent — update cross-references when one changes.

## 17. Anti-Patterns to Avoid

- **God-service Gateway:** stuffing scanner/ML logic directly into the Go Gateway "for speed" — breaks the async analysis-plane separation and the latency budget model; use the message bus.
- **Silent fallback verdicts:** scanner returning a default "safe" verdict on internal error instead of surfacing an explicit "scan unavailable" state to the policy engine.
- **Shadow policy logic:** implementing enforcement conditionals in application code that duplicate/bypass what Rego policy already expresses — creates two sources of truth.
- **Chatty synchronous fan-out:** Gateway synchronously calling multiple analysis services in sequence on the hot path instead of publishing once to the bus.
- **Config drift between environments:** hand-editing Helm values per environment instead of parameterizing via environment-specific values files under version control.
- **Alert fatigue by design:** shipping detection logic with unmeasured/unbounded false-positive rates instead of validating against the benchmark corpus referenced in `prd.md §5, §12`.
- **Over-fetching in the UI:** dashboard components pulling full trace payloads when only summary/aggregate data is needed for the current view.
