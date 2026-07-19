# AEGIS — Implementation Phases

**Related docs:** `prd.md`, `architecture.md`, `rules.md`, `design.md`, `memory.md`
**Purpose:** a step-by-step build order so an AI coding agent (or a human team) can implement AEGIS incrementally without ambiguity, with each phase independently shippable/testable.

**Recommended implementation order:** Phase 0 → 1 → 2 → 3 → 4 (MVP boundary) → 5 → 6 → 7 → 8.

---

## Phase 0 — Foundations & Scaffolding

**Objective:** Stand up the monorepo, tooling, and empty-but-deployable skeleton so every later phase has a consistent base.

**Features included:**
- Monorepo scaffold per `architecture.md §3`.
- CI pipeline skeleton (lint, unit test, build, SBOM, scan, sign — `architecture.md §21`).
- Local dev environment (docker-compose for Postgres, Redis, Kafka/NATS, ClickHouse) for fast iteration before full k8s.
- Base Helm chart skeletons (empty deployments) and Terraform module skeletons for one cloud (AWS first — **[ASSUMPTION]**).
- Base Envoy config with mTLS termination stub.

**Dependencies:** none (starting point).

**Deliverables:**
- Buildable, empty services for gateway, policy-engine, scanner, rag-monitor, ui, redteam.
- Passing CI on an empty/skeleton commit.
- `docker-compose.dev.yml` bringing up all infra dependencies locally.

**Milestones:**
- M0.1: Repo scaffold merged, CI green.
- M0.2: Local docker-compose stack starts cleanly.
- M0.3: Skeleton Helm install succeeds on a local kind/minikube cluster.

**Testing requirements:** CI pipeline itself is the test artifact — verify each stage (lint/build/scan/sign) fails correctly on an intentionally broken sample commit, then passes on the real skeleton.

**Definition of Done:** A developer (or agent) can clone the repo, run one command to bring up local infra, and run `helm install` against a local cluster with all services reaching a healthy (even if no-op) state.

**Expected outputs:** Working monorepo, CI config, docker-compose file, empty Helm charts, Terraform skeleton.

**Risk assessment:** Low risk, high leverage — mistakes here compound into every later phase, so review folder/naming conventions against `rules.md` carefully before proceeding.

**Estimated complexity:** Low-Medium (mostly configuration, not logic).

---

## Phase 1 — Identity, Ingress & Reverse Proxy (Core Data Plane Skeleton)

**Objective:** Get a real request flowing end-to-end through Envoy → Gateway → a real AI backend, with authentication enforced, before adding any detection intelligence.

**Features included:** (from PRD feature list) #2 Reverse proxy with mTLS/JWT/OAuth2/workload identity; #22 Secret zero bootstrap.

**Dependencies:** Phase 0 complete.

**Deliverables:**
- Envoy configured to terminate TLS and route to the Gateway sidecar.
- Gateway sidecar validates JWT/OAuth2 tokens (Keycloak/Hydra integration) and forwards authenticated requests to a real target (start with one OpenAI-compatible endpoint).
- Workload identity (SPIFFE/SPIRE or cloud-native equivalent) issuing service identities; secret-zero bootstrap working against Vault or cloud secret manager.
- Structured logging + correlation ID propagation (`rules.md §16.4`).

**Milestones:**
- M1.1: Unauthenticated request correctly rejected (401).
- M1.2: Authenticated request successfully proxied to a real LLM provider and response returned unmodified.
- M1.3: mTLS enforced between Envoy↔Gateway and Gateway↔backend where applicable.

**Testing requirements:** Integration test hitting the live proxy path with valid/invalid/expired tokens; latency benchmark to establish the baseline before any enforcement logic is added (this becomes the reference point for the `architecture.md §18` budget).

**Definition of Done:** A client can send a chat-completion request through AEGIS to a real backend and get a correct response, with full authn/identity enforcement and no detection/policy logic yet (pure transparent, authenticated proxy).

**Expected outputs:** Working reverse proxy with authn, base latency benchmark numbers, secret-zero bootstrap docs.

**Risk assessment:** Medium — identity/mTLS misconfiguration is a common source of hard-to-debug outages; allocate explicit time for failure-mode testing (expired cert, revoked token, IdP unavailable).

**Estimated complexity:** Medium.

---

## Phase 2 — Policy Engine & Enforcement Pipeline (Synchronous Path)

**Objective:** Introduce policy-as-code enforcement for the fast, synchronous checks — input sanitization, prompt-structure guardrails, tool-use allowlisting, hierarchical policy resolution — without yet involving the async LLM scanner.

**Features included:** #3 Multi-tenant policy engine; #4 Input sanitization; #6 Prompt structure enforcement; #7 Tool-use allowlisting.

**Dependencies:** Phase 1 complete (authenticated proxy path exists).

**Deliverables:**
- OPA/Rego policy engine service with hierarchical scope resolution (org → app → model → environment) per `architecture.md §9`.
- PostgreSQL schema for orgs/applications/model-endpoints/policies/policy-versions (`architecture.md §7`).
- Gateway integration: policy check step inserted into the request pipeline; local policy cache with fallback-to-last-known-good.
- Input sanitization module (prompt template injection neutralization, max-token budget, charset enforcement).
- Prompt structure guardrails (role segregation, max-turn enforcement).
- Tool-use allowlist with parameter validation.
- Standard error response shape (`rules.md §8`) for policy blocks.

**Milestones:**
- M2.1: Policy CRUD API functional; policy versions audit-logged.
- M2.2: A test policy blocking a specific tool-call parameter correctly blocks matching requests end-to-end.
- M2.3: Policy-engine outage correctly falls back to last-known-good cached policy on the Gateway (no full outage).

**Testing requirements:** `opa test` suite covering allow/deny cases per policy rule (`rules.md §12`); chaos test killing the policy engine mid-traffic to verify graceful degradation.

**Definition of Done:** Requests are enforced according to versioned, hierarchical Rego policy; policy changes propagate to Gateway replicas within the target window (`prd.md §5`); all blocks return the standardized error shape.

**Expected outputs:** Working policy engine + enforcement pipeline; baseline policy bundle checked into `infra/policies/`.

**Risk assessment:** Medium-High — this is the core value proposition; incorrect policy resolution order or caching bugs directly cause either security gaps (under-blocking) or outages (over-blocking). Extra review attention warranted.

**Estimated complexity:** High.

---

## Phase 3 — Event Bus & Async Analysis Plane Skeleton

**Objective:** Introduce the message bus and get events flowing from the Gateway to consumer services, without yet implementing real detection intelligence (stub scanner returning a fixed verdict first).

**Features included:** Infrastructure underpinning #8 (scanner), #9 (behavioral baselines), #10 (RAG monitor) — this phase builds the plumbing; intelligence comes in Phase 4.

**Dependencies:** Phase 2 complete (there are real enforcement events worth analyzing).

**Deliverables:**
- Kafka/NATS JetStream deployed; event schema defined (including a forward-looking `modality` field per `architecture.md §22`).
- Gateway emits request/response events to the bus.
- Stub Scanner, stub RAG Monitor, stub Baseline service consuming from the bus and writing a fixed/simple verdict to ClickHouse.
- Verdict feedback path: async verdict can be attached to a session/event and surfaced (not yet acted upon) in a minimal read-only UI view.

**Milestones:**
- M3.1: Events reliably flow Gateway → Bus → stub consumers with no message loss under load test.
- M3.2: ClickHouse ingesting event/detection records at the target rate.
- M3.3: A minimal UI page lists recent events with stub verdicts.

**Testing requirements:** Throughput/load test of the bus + ClickHouse ingestion path; consumer-lag alerting verified (KEDA autoscaling trigger fires correctly under synthetic load spike).

**Definition of Done:** The full async plumbing exists and is horizontally scalable, verified under load, even though the "intelligence" inside each consumer is still a stub.

**Expected outputs:** Working event bus integration, ClickHouse schema live, KEDA autoscaling configured, minimal trace-listing UI.

**Risk assessment:** Medium — getting the schema right early avoids painful migrations later; treat the event schema as a semi-frozen contract once Phase 4 consumers depend on it.

**Estimated complexity:** Medium.

---

## Phase 4 — Detection Intelligence (Scanner, RAG Monitor, Behavioral Baseline) — **MVP boundary**

**Objective:** Replace the stubs from Phase 3 with real detection logic. This phase completes the MVP as scoped in `prd.md §16 Assumptions` (Phases 1–4 = MVP).

**Features included:** #8 Dual-model LLM scanner; #9 Behavioral baselining; #10 RAG poisoning detection; #11 Data exfiltration heuristics; #12 Supply-chain verification; #13 OWASP/ATLAS mapping; #14 Session trace explorer/timeline replay.

**Dependencies:** Phase 3 complete (bus + storage plumbing proven at scale).

**Deliverables:**
- Dual-model ensemble scanner (two independent open-weight models via vLLM/TGI) with disagreement-escalation logic (`architecture.md §6.3`).
- RAG Monitor: embedding-neighborhood anomaly detection + citation/retrieval-drift tracking against the monitored vector DB.
- Behavioral baseline service: online clustering/changepoint detection per user/app.
- Data-exfiltration heuristics (encoding/steganography checks, chunked-transfer pattern detection, anomalous URL detection).
- Supply-chain verification: model digest checks, provenance metadata, MCP/plugin attestation checks.
- OWASP LLM Top 10 / MITRE ATLAS mapping layer applied to every detection.
- Session trace explorer UI: full prompt → tool-call → output timeline replay.
- Benchmark validation against the false-positive/detection-rate targets in `prd.md §5, §12`.

**Milestones:**
- M4.1: Scanner achieves target detection rate on the injection/jailbreak benchmark corpus (`prd.md §12`).
- M4.2: RAG poisoning detector flags injected poisoned documents in a controlled test RAG corpus.
- M4.3: Trace explorer supports full replay of a real flagged session end-to-end in the UI.
- M4.4: All detections carry correct OWASP/ATLAS tags with evidence.

**Testing requirements:** Full benchmark-corpus evaluation (precision/recall); adversarial meta-injection test against the scanner itself (`rules.md §7`); UI e2e test of the trace-replay flow.

**Definition of Done:** AEGIS can discover, enforce, detect (with real intelligence, not stubs), and let an analyst investigate a real incident end-to-end through the trace explorer — this is the shippable MVP.

**Expected outputs:** Working detection intelligence across all three analysis services; benchmark evaluation report; functioning trace explorer.

**Risk assessment:** High — this is the most technically challenging phase (ML/detection quality, false-positive tuning, meta-injection robustness). Budget the most iteration time here; do not treat first-pass thresholds as final.

**Estimated complexity:** High.

---

## Phase 5 — Red-Team Harness & Continuous Validation

**Objective:** Build the automated adversarial testing capability that continuously validates Phase 4's detection quality.

**Features included:** #15 Red-team harness; #20 Self-hosted red-team test-case registry.

**Dependencies:** Phase 4 complete (there must be real detectors to validate against).

**Deliverables:**
- Automated red-team test generators for prompt injection/jailbreak, RAG corruption, and unsafe tool invocation, targeting MITRE ATLAS technique coverage.
- Versioned, self-hosted test-case registry (`apps/redteam/testcases/`).
- Scheduled runner (nightly regression) producing pass/fail results with evidence, stored in ClickHouse.
- UI view for red-team run history and technique-coverage tracking against the `prd.md §5` KPI target.

**Milestones:**
- M5.1: Harness executes ≥30 MITRE ATLAS-mapped test cases automatically (MVP coverage target from `prd.md §5`).
- M5.2: A newly discovered detection gap can be codified as a new test case and added to the regression registry within the same sprint it was found.

**Testing requirements:** Verify the harness itself doesn't generate false "pass" results (i.e., validate the harness's own test cases produce a genuine detection when run against a known-vulnerable configuration).

**Definition of Done:** Nightly automated red-team runs execute against `main`/staging and regress-test all previously discovered techniques, with coverage tracked toward the growth-phase KPI target.

**Expected outputs:** Working red-team harness, versioned test-case registry, coverage dashboard.

**Risk assessment:** Medium — risk is mainly in keeping generated test cases realistic/non-trivial (avoid a harness that only tests things detectors already trivially catch).

**Estimated complexity:** Medium-High.

---

## Phase 6 — Observability, SIEM/SOAR & Compliance Reporting

**Objective:** Complete the operational integration surface so AEGIS fits into existing SOC and compliance workflows.

**Features included:** #17 OTLP/Prometheus/Grafana export + alerting; #18 SIEM (Syslog/CEF) and SOAR (webhook/STIX-TAXII) integration; #16 Compliance reports (EU AI Act, ISO 42001).

**Dependencies:** Phase 4 (real detections to export) and Phase 5 (red-team results to include in compliance evidence) complete.

**Deliverables:**
- OpenTelemetry Collector pipeline wired from all services to Prometheus/Grafana and PagerDuty/Opsgenie alerting.
- CEF/Syslog export of detections to a test SIEM (e.g., a local Elastic/Splunk instance for validation).
- STIX/TAXII and webhook export to a test SOAR.
- Compliance evidence-bundle generator producing EU AI Act / ISO 42001-aligned reports referencing real policy, detection, and red-team data.

**Milestones:**
- M6.1: A detection event appears correctly formatted in a test SIEM within the target latency.
- M6.2: A generated compliance evidence bundle correctly links to underlying policy/detection/red-team records.
- M6.3: PagerDuty/Opsgenie alert fires correctly on a high-severity detection.

**Testing requirements:** End-to-end integration test against real (or realistic sandbox) SIEM/SOAR targets; compliance report content review against actual EU AI Act / ISO 42001 control lists.

**Definition of Done:** A security team using existing SIEM/SOAR tooling and a compliance lead needing audit evidence can both get what they need directly from AEGIS without a bespoke console being the only interface.

**Expected outputs:** Working observability/export pipeline, sample compliance report artifacts.

**Risk assessment:** Low-Medium — mostly integration work against well-documented standards, but format/schema accuracy (CEF, STIX/TAXII) needs careful validation.

**Estimated complexity:** Medium.

---

## Phase 7 — Multi-Cloud Deployment, GitOps & Hardening

**Objective:** Generalize deployment beyond the single-cloud Phase-0 assumption and harden the supply chain end-to-end.

**Features included:** #21 Multi-cloud Helm/Terraform/GitOps; supply-chain hardening across #12 and CI/CD (`architecture.md §21`).

**Dependencies:** Phases 1–6 functionally complete on the initial cloud target.

**Deliverables:**
- Terraform modules for the remaining two clouds (GCP, Azure) mirroring the AWS module from Phase 0/1.
- Full SLSA L2+ provenance and admission-controller signature verification (Kyverno/Sigstore policy-controller) enforced in all environments.
- GitOps pipeline (ArgoCD/Flux) reference implementation for policy and Helm-values promotion across environments.
- Chaos/failure-mode test suite (control-plane outage, message-bus partition, scanner overload) run against a staging multi-cloud-representative environment.

**Milestones:**
- M7.1: Full stack deploys cleanly on GCP and Azure using the same Helm charts with only Terraform/values differences.
- M7.2: Unsigned image is correctly rejected by the admission controller in every environment.
- M7.3: Chaos suite passes with documented, expected degradation behavior (matching fail-open/fail-closed policy from `prd.md §16`).

**Testing requirements:** Cross-cloud deployment smoke tests; chaos engineering scenarios documented with expected vs. actual behavior.

**Definition of Done:** AEGIS is deployable, with equivalent behavior and hardening guarantees, across AWS/GCP/Azure, satisfying the multi-cloud constraint in `prd.md`.

**Expected outputs:** Complete multi-cloud Terraform/Helm set, chaos test report, GitOps reference pipeline.

**Risk assessment:** Medium — cloud-specific IAM/identity federation quirks are the most likely source of friction.

**Estimated complexity:** Medium-High.

---

## Phase 8 — Growth Features (Post-MVP, Toward `prd.md §17` Future Scope)

**Objective:** Begin building toward the longer-term product vision once the core platform (Phases 0–7) is stable in production/portfolio use.

**Features included:** Multi-tenant AEGIS Cloud groundwork; policy marketplace groundwork; expanded provider/framework coverage; expanded MITRE ATLAS technique coverage (target ≥80 per `prd.md §5`); research-oriented scanner precision work.

**Dependencies:** Phases 0–7 complete and stable.

**Deliverables (indicative — to be re-scoped against real usage data before starting):**
- Full multi-tenant data isolation hardening for a hosted "AEGIS Cloud" offering.
- Policy bundle signing/distribution mechanism suitable for a marketplace model.
- Additional provider adapters (beyond the MVP set) using the adapter pattern from `architecture.md §22`.
- Expanded red-team technique coverage toward the growth-phase KPI target.

**Milestones:** To be defined once Phase 0–7 usage data exists — **[ASSUMPTION]** this phase is intentionally left more loosely specified since it depends on real-world feedback the earlier phases will generate.

**Testing requirements:** Same rigor as prior phases; multi-tenancy work specifically requires dedicated isolation/penetration testing before any hosted offering goes live.

**Definition of Done:** Defined per sub-initiative at the time each is scoped; this phase is a backlog container, not a single deliverable.

**Expected outputs:** Varies by sub-initiative; tracked as its own set of phase documents once prioritized.

**Risk assessment:** Variable/TBD — explicitly deferred, not risk-assessed in detail here to avoid speculative planning ahead of real data.

**Estimated complexity:** TBD.

---

## Cross-Phase Notes

- **MVP = Phases 0–4.** Everything from Phase 5 onward is valuable but not required for a functioning, demonstrable product, consistent with `prd.md §16 Assumptions`.
- Each phase's Definition of Done must pass before the next phase begins in earnest, though limited parallel work (e.g., starting UI polish for Phase 4 while Phase 3 plumbing is finishing) is acceptable if it doesn't violate `rules.md` review gates.
- Any phase that reveals a needed change to `architecture.md` or `prd.md` must update those documents in the same PR/change set (see `rules.md §14` checklist), keeping the documentation set internally consistent throughout the build.
