# AEGIS — Product Requirements Document (PRD)

**Product:** AEGIS (AI Enforcement Gateway & Intelligence for Security)
**Document type:** Product Requirements Document
**Status:** Draft v1.0 — foundation for AI-agent-assisted implementation
**Related docs:** `architecture.md`, `rules.md`, `phases.md`, `design.md`, `memory.md`

---

## 1. Executive Summary

AEGIS is a cloud-native runtime security gateway for AI and agentic AI systems. It sits between consumers (apps, users, agents) and LLM/RAG/agent backends (OpenAI, Anthropic, Azure OpenAI, vLLM, TGI, LangChain/CrewAI/AutoGen agents), providing **discovery, enforcement, detection, and adversarial testing** in a single extensible platform.

Unlike static model scanners or one-time red-team tools, AEGIS operates on **live traffic**, enforcing policy in-line (block, redact, tag-and-allow) and continuously detecting runtime abuse (prompt injection, jailbreaks, RAG poisoning, data exfiltration, unsafe tool invocation). Findings are mapped to industry-standard taxonomies (OWASP LLM Top 10, MITRE ATLAS) and exported to existing security operations tooling (SIEM/SOAR).

## 2. Vision

> **To become the default runtime control plane for AI traffic** — the "API gateway + WAF + CASB" equivalent for LLMs and agentic AI — so that any organization can adopt AI capabilities without operating them blind.

AEGIS should be:
- **Vendor-agnostic** — works with any OpenAI-compatible, Anthropic, or self-hosted model endpoint.
- **Enforcement-first** — not just alerting, but actively able to block/redact/rewrite unsafe traffic.
- **Open and extensible** — policies-as-code, pluggable detectors, open detection registry.
- **Operationally native** — fits into existing SOC/DevSecOps workflows (SIEM, SOAR, GitOps) rather than requiring a new console silo.

## 3. Problem Statement

Enterprises are shipping LLM-enabled and agentic AI applications faster than security teams can gain visibility into them. Specific gaps this product addresses:

| Gap | Evidence / Rationale |
|---|---|
| No runtime visibility into AI-specific abuse (prompt injection, indirect/RAG injection, jailbreaks, unsafe tool calls, data exfiltration via AI channels) | OWASP Top 10 for LLM Applications (2025) codifies these as top risks (e.g., LLM01 Prompt Injection, LLM03 Supply Chain) |
| Adversaries have a catalogued, expanding technique set against AI systems | MITRE ATLAS matrix documents 16 tactics and 173+ techniques |
| "Shadow AI" — unsanctioned AI tool usage — is widespread and largely invisible to security teams | Up to 88% of security leaders report employees using unapproved AI tools |
| API-related vulnerabilities, many AI/LLM-linked, are a growing share of actively exploited vulnerabilities | 43% of 2025 CISA KEV additions were API-related |
| No unified product combines discovery + enforcement + behavioral detection + adversarial testing for AI | Market has point solutions (static model scanners, prompt-linting libraries) but no runtime, enforcing gateway |

**Core problem statement:** *Security teams cannot see, control, or continuously test the AI traffic flowing through their organization, leaving a fast-growing, high-impact attack surface effectively unmonitored and unenforced.*

## 4. Goals

### 4.1 Product Goals
1. Provide real-time discovery of all AI/LLM/agent endpoints in an environment.
2. Enforce security policy at the API layer for AI traffic (input, output, tool-use).
3. Detect and classify AI-specific abuse using layered heuristic + LLM-based detection.
4. Continuously validate defenses via an automated adversarial red-team harness.
5. Integrate cleanly into existing enterprise security operations (SIEM/SOAR/compliance reporting).

### 4.2 Business Goals
1. Establish AEGIS as a credible open-source-first security control for AI infrastructure (portfolio/community adoption).
2. Create a foundation extensible into a managed SaaS ("AEGIS Cloud") and/or commercial policy marketplace.
3. Demonstrate technical depth suitable for security engineering hiring portfolios and potential research publication.

## 5. Success Metrics & KPIs

| Metric | Target (v1 / MVP) | Target (v2 / Growth) |
|---|---|---|
| Mean detection latency (request → verdict) | < 250ms (non-streaming), < 100ms added latency for streaming | < 100ms all paths |
| False positive rate on legitimate creative/technical prompts (benchmark set) | < 8% | < 3% |
| Red-team technique coverage (MITRE ATLAS techniques with automated test cases) | ≥ 30 techniques | ≥ 80 techniques |
| Endpoint discovery accuracy (known test environment) | ≥ 90% of AI endpoints identified | ≥ 98% |
| Policy update propagation time (control plane → data plane) | < 5s | < 1s |
| Concurrent session tracking (per gateway replica) | 10K sessions | 1M+ sessions (horizontally scaled) |
| SIEM/SOAR integration coverage | 1 SIEM export format (CEF/Syslog) + webhook | STIX/TAXII + Sigma + native Splunk/Sentinel apps |
| Uptime of enforcement path | 99.9% | 99.95% |

**Assumption:** exact numeric targets above are industry-informed estimates (not sourced from user-provided material) and should be revisited once real traffic/benchmark data exists. Marked explicitly as **[ASSUMPTION]**.

## 6. User Personas

### 6.1 Security Engineer / AI Security Analyst — "Priya"
- **Role:** Operates the SOC, triages detections, tunes policies.
- **Goals:** Fast, low-noise detection; clear evidence per alert; ability to simulate policy changes before rollout.
- **Pain points today:** No visibility into AI traffic at all; relies on manual log grep or nothing.

### 6.2 CISO / DPO / Compliance Lead — "Marcus"
- **Role:** Owns AI governance risk and regulatory exposure (EU AI Act, ISO 42001).
- **Goals:** Auditable evidence of controls, compliance reports, risk dashboards for the board.
- **Pain points today:** Cannot demonstrate AI-specific controls to auditors; Shadow AI is an unknown liability.

### 6.3 DevSecOps / Platform Engineer — "Jordan"
- **Role:** Deploys and operates AEGIS as infrastructure (Kubernetes, Terraform, Helm).
- **Goals:** GitOps-friendly deployment, minimal operational overhead, clear SLOs, safe rollout/rollback of policy.
- **Pain points today:** Security tooling that requires bespoke, non-declarative configuration.

### 6.4 AI Product Team Engineer — "Sam"
- **Role:** Builds LLM-powered features and embeds AI/agents into a SaaS product.
- **Goals:** Ship fast without becoming a security bottleneck; get clear, actionable guidance when blocked (not just an opaque 403).
- **Pain points today:** Security controls that break legitimate use cases with no explanation.

## 7. User Journey (Representative — Priya, Security Engineer)

1. **Discovery:** AEGIS is deployed as a sidecar/reverse proxy in the Kubernetes cluster; within minutes it surfaces a live inventory of AI endpoints (internal LLM services, external API calls, MCP servers).
2. **Baseline policy:** Priya applies a default policy tier (org → app → model → environment hierarchy) via Rego policy-as-code, version-controlled in Git.
3. **Live monitoring:** Traffic begins flowing through the gateway. The scanner and RAG monitor tag/verdict each session; the token/session trace explorer lets Priya replay a suspicious session.
4. **Detection triage:** An alert fires — indirect prompt injection detected in a RAG-sourced document. AEGIS maps it to OWASP LLM01 / MITRE ATLAS technique, provides evidence (the offending retrieved chunk, the model's response, confidence score).
5. **Response:** Priya adjusts policy (block class of injected pattern), simulates the change against recent traffic, and deploys it — propagated to all gateway replicas within seconds.
6. **Continuous validation:** The red-team harness runs the new adversarial test case nightly as a regression check going forward.
7. **Reporting:** Marcus pulls a compliance evidence bundle for the quarterly AI governance review, sourced automatically from AEGIS's detection and policy audit logs.

## 8. User Stories

| ID | As a... | I want to... | So that... |
|---|---|---|---|
| US-01 | Security Engineer | see all AI endpoints in my environment automatically | I don't rely on manual inventory |
| US-02 | Security Engineer | define policy hierarchically (org/app/model/env) | I can apply broad defaults and narrow overrides |
| US-03 | Security Engineer | see a request replay (prompt → tool calls → output) | I can understand what happened in a flagged session |
| US-04 | CISO | export compliance evidence bundles | I can satisfy EU AI Act / ISO 42001 audit requirements |
| US-05 | DevSecOps engineer | deploy AEGIS via Helm/Terraform | I can manage it like the rest of my infrastructure |
| US-06 | AI product engineer | get a clear reason when a request is blocked | I can fix my prompt/integration instead of guessing |
| US-07 | Security Engineer | run automated red-team tests against my RAG pipeline | I can validate defenses before a real attacker does |
| US-08 | Security Engineer | route detections to my SIEM/SOAR | I don't have to monitor a separate console |
| US-09 | Platform Engineer | roll out a policy change safely (canary/simulate) | I don't break legitimate traffic |
| US-10 | Security Engineer | see detections mapped to OWASP LLM Top 10 / MITRE ATLAS | I can communicate risk using standard taxonomies |

## 9. Functional Requirements

### 9.1 Discovery
- FR-1: System shall passively infer AI endpoints from proxied traffic patterns (request paths, payload shape).
- FR-2: System shall parse agent/MCP manifests where available to enrich discovery.
- FR-3: Discovered endpoints shall be classified by provider type (OpenAI-compatible, Anthropic, self-hosted, agent framework).

### 9.2 Enforcement
- FR-4: System shall provide a reverse proxy enforcing mTLS and JWT/OAuth2 validation before requests reach backend AI services.
- FR-5: System shall support configurable input sanitization (prompt template injection neutralization, max-token budgets, charset enforcement).
- FR-6: System shall support configurable output filtering (PII/secret masking, toxicity flagging, hallucination metadata tagging).
- FR-7: System shall enforce prompt structure guardrails (system/user/assistant role segregation, maximum conversation turns/chain depth).
- FR-8: System shall support tool-use allowlisting with parameter-level validation for agentic tool calls.
- FR-9: Policy engine shall support hierarchical scope resolution: organization → application → model → environment.
- FR-10: Enforcement actions shall support at minimum: **allow**, **allow-and-tag**, **redact/rewrite**, **block**.

### 9.3 Detection
- FR-11: System shall run a dual-model LLM ensemble scanner for prompt injection/jailbreak detection, escalating on model disagreement.
- FR-12: System shall build per-user/per-app behavioral baselines (sequence length, tool-call patterns, time-of-day, error rates) and flag deviations.
- FR-13: System shall detect RAG poisoning via embedding-neighborhood anomaly detection and citation/retrieval-distribution drift.
- FR-14: System shall detect candidate data-exfiltration patterns (steganographic encoding, repeated small chunked transfers, anomalous embedded URLs).
- FR-15: System shall verify LLM/plugin/MCP supply-chain integrity (model digest checks, provenance metadata, attestation).
- FR-16: Every detection shall be mapped to relevant OWASP LLM Top 10 categories and MITRE ATLAS tactics/techniques, with narrative and supporting evidence.

### 9.4 Investigation & Reporting
- FR-17: System shall provide a token/session trace explorer supporting timeline replay of prompt → tool calls → output.
- FR-18: System shall generate compliance reports (EU AI Act, ISO 42001, internal policy attestation) with linked evidence.
- FR-19: System shall export detections/metrics via OTLP to Prometheus/Grafana and alert via PagerDuty/Opsgenie.
- FR-20: System shall integrate with SIEM via Syslog/CEF and with SOAR via webhooks and STIX/TAXII.

### 9.5 Red-Teaming
- FR-21: System shall include an automated red-team harness generating prompt injection/jailbreak/RAG-corruption/unsafe-tool-invocation test cases.
- FR-22: System shall maintain a self-hosted, versioned registry of red-team test cases for regression testing.

### 9.6 Deployment & Operations
- FR-23: System shall be deployable via Helm charts and Terraform modules across AWS/GCP/Azure.
- FR-24: System shall bootstrap credentials via workload identity / secretless backends (e.g., HashiCorp Vault) — "secret zero" pattern.

## 10. Non-Functional Requirements

| Category | Requirement |
|---|---|
| **Performance** | Enforcement path must add ≤100ms latency at p95 for streaming responses (SSE-compatible chunked enforcement). |
| **Scalability** | Must support horizontal scaling of gateway sidecars and scanner microservices independently (KEDA-based autoscaling on queue depth). |
| **Availability** | Enforcement path target 99.9%+ uptime; scanner/RAG-monitor failures must fail to a configurable safe-default (fail-open or fail-closed per policy). |
| **Security** | All inter-service traffic mTLS; distroless, signed container images (Sigstore/Cosign); no plaintext secrets at rest. |
| **Multi-tenancy** | Policy and data isolation must be enforceable per tenant/org. |
| **Auditability** | All policy changes and detections must be immutably logged with actor, timestamp, and diff. |
| **Extensibility** | Detection modules and policy rules must be pluggable without core redeploy (policy-as-code via OPA/Rego). |
| **Portability** | Must run identically across AWS/GCP/Azure Kubernetes distributions (no cloud-specific hard dependencies in core). |
| **Compliance** | Must produce evidence artifacts mappable to EU AI Act and ISO 42001 control requirements. |

## 11. Feature Breakdown

The 22 core features from the project brief are grouped into delivery-oriented capability areas (see `phases.md` for sequencing):

1. **Discovery & Inventory** — endpoint auto-discovery, MCP/agent manifest parsing.
2. **Ingress & Identity** — reverse proxy, mTLS, JWT/OAuth2, workload identity, secret-zero bootstrap.
3. **Policy Engine** — multi-tenant hierarchical policy (OPA/Rego), input sanitization, output filtering, prompt-structure guardrails, tool-use allowlisting.
4. **Detection Engine** — dual-model scanner, behavioral baselining, RAG poisoning detection, exfiltration heuristics, supply-chain verification.
5. **Triage & Investigation** — OWASP/ATLAS mapping, session trace explorer/timeline replay.
6. **Red-Team Harness** — automated adversarial test generation, self-hosted test-case registry.
7. **Compliance & Reporting** — evidence bundles, EU AI Act/ISO 42001 report templates.
8. **Observability & Integration** — OTLP/Prometheus/Grafana, PagerDuty/Opsgenie, SIEM (Syslog/CEF), SOAR (webhook/STIX-TAXII).
9. **Platform & Delivery** — Helm/Terraform/GitOps, multi-cloud support.

## 12. Acceptance Criteria (Representative, MVP scope)

- **Discovery:** Given traffic flowing through the proxy for 15 minutes in a test environment with 5 known AI endpoints, AEGIS lists ≥90% of them with correct provider classification.
- **Enforcement:** Given a policy blocking a specific tool-call parameter pattern, a request matching that pattern is blocked and logged with policy ID and reason within 250ms.
- **Detection:** Given a known prompt-injection test corpus, the dual-model scanner flags ≥90% of injected prompts with an OWASP LLM01 tag and confidence score.
- **Trace explorer:** Given a completed session, a user can replay the full prompt → tool-call → output timeline in the UI.
- **Red-team:** Given a scheduled run, the harness executes ≥30 MITRE ATLAS-mapped test cases and produces pass/fail results with evidence.
- **SIEM export:** Given a detection event, a corresponding CEF-formatted syslog message is emitted within 5 seconds.

## 13. Edge Cases

- Streaming responses that are only partially generated before a policy violation is detected mid-stream — system must be able to truncate/redact in-flight.
- Multi-turn sessions where injected content appears many turns before the exploit triggers (delayed-payload injection).
- Legitimate security-research or red-team traffic that mimics attack patterns (must support an explicit "test mode" tag to avoid false incident escalation).
- Backend AI provider outage or rate-limiting — gateway must surface a clear upstream-error distinct from a policy block.
- Policy engine unavailable (control plane down) — data plane must have a safe fallback (last-known-good policy cached locally).
- Non-UTF8 or binary payloads sent to endpoints expecting text (charset enforcement edge case).
- Very large tool-call parameter payloads used to smuggle exfiltrated data across multiple small calls (must be caught by exfiltration heuristics, not just single-request inspection).
- Encrypted/tunneled traffic to AI endpoints that the proxy cannot terminate (must be documented as an explicit limitation, not silently ignored).

## 14. Risks

| Risk | Impact | Likelihood | Mitigation |
|---|---|---|---|
| Scanner false positives disrupt legitimate creative/technical use cases | High (user trust) | Medium | Dual-model ensemble with disagreement escalation rather than single-model block; tunable confidence thresholds |
| Added latency degrades UX for streaming chat experiences | High | Medium | Chunked/streaming-aware enforcement; async scanning with tag-and-allow fallback under policy |
| Adversarial prompts specifically target the scanner itself (meta-injection) | High | Medium | Isolate scanner prompt context from user-controlled input; independent ensemble voting |
| Multi-cloud/K8s deployment complexity slows adoption | Medium | Medium | Ship opinionated Helm defaults; document a single-cluster "quickstart" path first |
| Data/session traces themselves become a sensitive data store (traces may contain PII) | High | Medium | Field-level redaction in trace storage; role-based access to raw traces |
| Open-source core cannibalizes commercial/SaaS ambitions | Medium | Low | Core gateway/scanner open; premium UI/red-team packs and managed hosting as differentiators |

## 15. Constraints

- Must support both cloud-hosted (OpenAI, Anthropic, Azure OpenAI) and self-hosted (vLLM, TGI) model backends from day one.
- Must operate without requiring source-code changes to the protected AI application (transparent proxy model) as the default integration path; SDK-based integration is optional/secondary.
- Initial supported orchestration platform is Kubernetes; non-Kubernetes (bare VM/docker-compose) support is explicitly deferred (see Out-of-Scope).
- Team/resourcing is assumed to be small (portfolio/early-stage project), so MVP scope must be achievable incrementally — see `phases.md`.

## 16. Assumptions

**[ASSUMPTION]** The following are inferred, industry-standard decisions where the source material did not specify an exact answer:
1. Default deployment target is Kubernetes (EKS/GKE/AKS) — chosen because it's explicitly listed in the tech stack and is the dominant enterprise pattern for sidecar architectures.
2. Default LLM providers supported at MVP: OpenAI-compatible API + Anthropic API + one self-hosted option (vLLM) — full provider matrix expands in later phases.
3. Default enforcement failure mode is **fail-open with alert** for the detection/scanner path (to avoid breaking production AI features) but **fail-closed** for authentication/identity failures. This is a security/availability trade-off that should be revisited with actual stakeholders and is explicitly configurable per tenant policy.
4. Numeric performance/KPI targets in Section 5 are estimates, not sourced from the brief; to be validated with real benchmarking.
5. "v1/MVP" scope is defined as Phases 1–4 in `phases.md`; "v2/Growth" as Phases 5+.

## 17. Future Scope

- **AEGIS Cloud:** multi-tenant managed SaaS offering with per-customer AI endpoint isolation.
- **AI WAF & CASB marketplace:** vendor-agnostic policy marketplace, similar to WAF rule marketplaces today.
- **Formal verification research:** formal guarantees on scanner precision under adversarial prompting (potential thesis-level research direction).
- **Large-scale measurement study:** empirical research paper on AI abuse pattern prevalence mapped to MITRE ATLAS technique frequency.
- **Premium managed UI & red-team packs** on top of an open-source core gateway/scanner.

## 18. Out-of-Scope (v1)

- Non-Kubernetes deployment targets (bare-metal/VM-only, docker-compose-only production use) — quickstart/dev mode may use docker-compose, but it is not a supported production target initially.
- Fine-tuning or training of the detection LLMs themselves (v1 uses off-the-shelf open models, e.g., Llama/Mistral, via vLLM/TGI).
- Native mobile client applications for the SOC dashboard.
- Formal certification against EU AI Act/ISO 42001 (AEGIS produces *evidence*, not certification itself).
- Support for AI providers/frameworks beyond the initially targeted set (expansion tracked as backlog, not MVP).

## 19. Glossary

| Term | Definition |
|---|---|
| **AEGIS** | AI Enforcement Gateway & Intelligence for Security — this product. |
| **LLM** | Large Language Model. |
| **RAG** | Retrieval-Augmented Generation — pattern where an LLM's context is augmented with retrieved documents. |
| **Prompt Injection** | Attack where crafted input causes an LLM to deviate from intended behavior (direct) or is smuggled via retrieved/external content (indirect). |
| **Jailbreak** | Technique used to bypass an LLM's safety/behavioral guardrails. |
| **MCP** | Model Context Protocol — standard for connecting AI agents/models to external tools/data. |
| **OWASP LLM Top 10** | OWASP's ranked list of top security risks specific to LLM applications. |
| **MITRE ATLAS** | Adversarial Threat Landscape for Artificial-Intelligence Systems — MITRE's knowledge base of adversary tactics/techniques against AI systems. |
| **SIEM** | Security Information and Event Management system. |
| **SOAR** | Security Orchestration, Automation and Response platform. |
| **STIX/TAXII** | Standards for structured threat-intelligence representation and exchange. |
| **OPA/Rego** | Open Policy Agent and its policy language, used for policy-as-code. |
| **Secret Zero** | The bootstrapping problem of how a service obtains its first credential securely, typically solved via workload identity. |
| **KEDA** | Kubernetes Event-Driven Autoscaling. |
| **SLSA** | Supply-chain Levels for Software Artifacts — a framework for build integrity. |
