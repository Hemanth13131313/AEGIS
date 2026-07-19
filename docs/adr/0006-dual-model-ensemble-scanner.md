# ADR 0006: Dual-model ensemble scanner with structural isolation for meta-injection prevention

## Status
Accepted

## Context
Single-model classifiers have blind spots and are vulnerable to prompt injection targeting the classifier itself (meta-injection). If a classifier is given an adversarial payload that overrides its own classification instructions, it may falsely report "SAFE".

## Decision
- Implement a dual-model ensemble (e.g., Mistral + Llama) running independently and concurrently.
- Enforce strict structural isolation: the user payload must be encapsulated within a JSON string inside the `user` message, while classification instructions are kept exclusively in the `system` message.
- Any disagreement between the models on whether to allow or block will result in a `tag` action, triggering an alert for human review.

## Alternatives Considered
- Single model: Rejected due to significant blind spots and susceptibility to meta-injection.
- Regex only: Rejected as too brittle for adversarial evasion techniques.
- Embedding similarity: Deferred to Phase 7 due to implementation complexity and resource requirements.

## Consequences
- Requires 2x LLM compute cost, which increases operational expenses.
- Higher precision and strong meta-injection mitigation.
- Demands concurrent execution (via `asyncio.gather`) to maintain a P95 latency of ≤500ms.
- Requires vLLM or similar infrastructure for efficient local serving.
