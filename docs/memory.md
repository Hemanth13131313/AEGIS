# AEGIS — Project Memory

**Purpose:** This is AEGIS's persistent long-term memory. Any AI coding agent (or human) starting a new session should read this file first. It summarizes what must be remembered across sessions so decisions stay consistent. When something here changes, update this file in the same change set (see `rules.md §13`).

**Full detail lives in:** `prd.md` (product), `architecture.md` (system design), `rules.md` (engineering standards), `phases.md` (build sequence), `design.md` (UI/design system). This file is the *summary index*, not a replacement for those documents.

---

## 1. Project Overview

**AEGIS (AI Enforcement Gateway & Intelligence for Security)** is a cloud-native, runtime security gateway that sits in front of LLM providers, RAG backends, and agentic AI frameworks. It discovers AI endpoints, enforces policy in-line (block/redact/tag-and-allow), detects AI-specific abuse (prompt injection, jailbreaks, RAG poisoning, data exfiltration, unsafe tool calls), continuously validates itself via an adversarial red-team harness, and integrates with existing SIEM/SOAR/compliance workflows.

It is differentiated from static model scanners by operating on **live traffic** with real **enforcement**, not just alerting.

## 2. Core Objectives

1. Real-time discovery of AI/LLM/agent endpoints.
2. In-line policy enforcement at the API layer.
3. Multi-signal detection (heuristic + dual-model LLM ensemble) mapped to OWASP LLM Top 10 / MITRE ATLAS.
4. Continuous adversarial validation via automated red-teaming.
5. Native integration into existing SOC tooling (SIEM/SOAR) and compliance reporting (EU AI Act, ISO 42001).

## 3. Business Goals

- Establish credibility as an open-source-first AI security control.
- Build a foundation extensible into a managed SaaS ("AEGIS Cloud") and policy marketplace.
- Serve as a strong technical portfolio piece (security engineering depth) and potential research basis.

## 4. Critical Assumptions (flagged, not sourced from original brief)

These were introduced because the source material didn't specify them. Revisit if real constraints emerge:

| # | Assumption | Where documented |
|---|---|---|
| 1 | Kubernetes (EKS/GKE/AKS) is the primary deployment target; non-K8s production deployment is out of scope for v1 | `prd.md §15, §18`, `architecture.md §19` |
| 2 | MVP provider support = OpenAI-compatible API + Anthropic API + vLLM (self-hosted); others deferred | `prd.md §16` |
| 3 | Default fail-mode: **fail-open with alert** for detection/scanner path, **fail-closed** for auth/identity failures (configurable per policy) | `prd.md §16`, `rules.md §15.5` |
| 4 | Numeric KPI targets in `prd.md §5` are estimates pending real benchmarking, not sourced figures | `prd.md §5` |
| 5 | MVP scope = Phases 0–4 in `phases.md`; Phase 5+ is post-MVP growth | `prd.md §16`, `phases.md` |
| 6 | Monorepo structure chosen over polyrepo for small-team simplicity | `architecture.md §3` |
| 7 | React Query + Zustand for frontend state (no Redux) | `architecture.md §5` |
| 8 | No existing brand assets — `design.md` palette/typography are an industry-informed default, swappable via design tokens | `design.md §4` |
| 9 | Phase 0/1 target AWS first; GCP/Azure Terraform modules follow in Phase 7 | `phases.md` Phase 0, Phase 7 |
| 10 | Env var naming convention: `AEGIS_<COMPONENT>_<SETTING>` | `architecture.md §20` |

## 5. Architecture Decisions (index — see `architecture.md` for full rationale)

- **Four-plane architecture:** data plane (Envoy + Gateway sidecar), control plane (Policy Engine/OPA), analysis plane (Scanner/RAG Monitor/Baseline, event-driven via Kafka/NATS), observability/integration plane.
- **Synchronous vs. asynchronous split:** fast checks (auth, sanitization, structural guardrails, tool-allowlist) are synchronous in the Gateway (≤10ms budget); deep LLM-based scanning is asynchronous by default, escalatable to synchronous for the highest-risk policy tiers.
- **Dual-model ensemble scanner** with disagreement-escalation, to reduce both false positives and single-model blind spots; scanner isolated from meta-injection via structural prompt separation.
- **Storage split:** PostgreSQL (policy/tenancy — strong consistency), ClickHouse (traces/detections/red-team results — high-volume append), Redis (cache/rate-limit/session-pointer, ephemeral).
- **Policy-as-code via OPA/Rego**, hierarchical resolution: environment → model → application → organization (most specific wins).
- **Identity:** OAuth2/OIDC (Keycloak/Hydra) for users; SPIFFE/SPIRE + mTLS for service-to-service; secret-zero bootstrap via Vault/cloud IAM federation — no static baked-in secrets.
- **Supply chain:** distroless + signed (Cosign) images, SBOM (CycloneDX), SLSA L2+ provenance, admission-controller signature enforcement.

## 6. Tech Stack (quick index — full table in `architecture.md §4`)

- Ingress: Envoy. Gateway: Go (Rust optional for latency path). Scanner/RAG Monitor: Python/FastAPI + vLLM/TGI (Llama/Mistral). Policy Engine: OPA/Rego. Frontend: React/TypeScript + Grafana panels. DB: PostgreSQL + ClickHouse + Redis. Bus: Kafka/NATS JetStream. Identity: Keycloak/Hydra + SPIFFE/SPIRE. IaC: Terraform + Helm + Kustomize. CI/CD: GitHub Actions/GitLab CI + SLSA. Observability: OpenTelemetry + Prometheus/Grafana + Tempo/Loki.

## 7. Design Principles (index — full system in `design.md`)

- Dark-mode-first, high-signal/low-noise "instrument panel" aesthetic for SOC analysts.
- Severity always communicated via color **+ icon + label**, never color alone (accessibility + operator trust).
- Evidence-first: every detection/block is one click from underlying proof.
- No hype/marketing language in-product ("AI-powered," "next-gen" banned from UI copy).
- All new UI must reuse existing design tokens/primitives before introducing new ones.

## 8. Coding Conventions (index — full rules in `rules.md`)

- Go: gofmt-enforced, explicit error wrapping, context-first function signatures.
- Python: PEP8/ruff/black, type-hinted, async I/O, no bare excepts.
- TypeScript: strict mode, function components + hooks only, no business logic in JSX.
- Rego: every policy package has `_test.rego` covering allow + deny cases.
- Conventional Commits (`type(scope): summary`); trunk-based Git workflow, squash-merge, required CI + review.

## 9. Important APIs / Contracts

- **External-facing:** transparent reverse-proxy pass-through of native provider APIs (`/v1/chat/completions`, `/v1/embeddings`, agent/tool endpoints) — no client SDK changes required for baseline protection.
- **Management API:** REST (`/api/v1/...`, cursor-paginated) for UI/integrations; gRPC (protobuf, backward-compatible evolution only) for control-plane↔data-plane policy distribution.
- **Standard error shape** (must be used everywhere): `{ error: { code, category: policy|upstream|internal, message, policy_ref?, request_id } }` — see `rules.md §8`.
- **Event bus schema:** includes a forward-looking `modality` field from day one (text-only populated in MVP) to avoid future breaking changes when multimodal input is added.

## 10. Business Rules

- Policy resolution order: environment override → model → application → organization default (most specific wins) — never resolved any other way.
- Every enforcement decision must flow through the Policy Engine (Rego) — no inline hard-coded enforcement logic in the Gateway, ever (`rules.md §15.4`).
- Raw prompt/response bodies are never logged at default verbosity; redaction is the default, not opt-in.
- Fail-closed for identity/auth failures; fail-open-with-alert for detection/scanner failures, unless a tenant's policy explicitly overrides this (must be an explicit, reviewed config change, never a silent code change).

## 11. Known Limitations

- Cannot inspect encrypted/tunneled traffic to AI endpoints that bypass proxy termination — documented limitation, not silently ignored (`prd.md §13`).
- No support for non-Kubernetes production deployment in v1.
- No fine-tuning/training of detection models in v1 (off-the-shelf open models via vLLM/TGI only).
- Mobile UI is read-only/summary-only; no policy editing or full trace replay on mobile (`design.md §11`).
- Numeric performance/coverage targets are estimates pending real benchmark validation.

## 12. Future Ideas (see `prd.md §17` for full detail)

- AEGIS Cloud (multi-tenant managed SaaS).
- AI WAF/CASB policy marketplace.
- Formal verification research on scanner precision under adversarial prompting.
- Large-scale empirical measurement study of AI abuse patterns vs. MITRE ATLAS technique prevalence.
- Premium managed UI / red-team packs layered on an open-source core.

## 13. Open Questions

- What are the *real* target latency/throughput/coverage numbers once actual traffic/benchmark data exists? (Currently estimates — `prd.md §5`.)
- Which cloud is the true primary target long-term — is AWS-first (current assumption) correct, or should GCP/Azure be prioritized differently based on target customer base?
- Should API-key-based client integration be promoted to a first-class path, or remain secondary to the transparent-proxy model indefinitely?
- What is the actual EU AI Act / ISO 42001 control checklist to validate the compliance-report generator against (needs real regulatory-text review, not just the summary used in `prd.md`)?
- Is a Redux-free frontend state approach (React Query + Zustand) still the right call once the UI's real complexity (trace explorer, policy simulation) is fully built out?

## 14. Developer Reminders

- Always check this file and `architecture.md` before proposing a new architectural component.
- Always map new detections to OWASP LLM Top 10 / MITRE ATLAS explicitly — this is a product differentiator, not an afterthought.
- Always write both allow and deny Rego tests for any new policy rule.
- Always keep the six documentation files cross-consistent — a change in one that affects another must update both in the same PR.
- MVP boundary is Phases 0–4 in `phases.md` — don't scope-creep growth features into the MVP definition of done.

## 15. Non-Negotiable Requirements

1. Every enforcement decision flows through the Policy Engine (Rego) — no bypasses.
2. No unsigned or CVE-unmitigated container image beyond local dev.
3. No raw prompt/response body logging at default verbosity.
4. No hard-coded secrets anywhere in source or config.
5. Every detection must carry an OWASP LLM Top 10 and/or MITRE ATLAS mapping.
6. Fail-closed for identity/auth; explicit, reviewed config (never silent code change) to alter fail-open/closed behavior elsewhere.
7. Severity in the UI is always color + icon + label — never color alone.
