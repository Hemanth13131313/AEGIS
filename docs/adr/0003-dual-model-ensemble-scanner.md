# ADR 0003: Use dual-model LLM ensemble for injection/jailbreak detection

**Status:** Accepted

## Context
Detecting sophisticated prompt injections and jailbreaks is challenging. A single detection model often has blind spots and vulnerabilities to specific adversarial patterns. Over-reliance on one model increases the risk of false positives and bypasses.

## Decision
We will use a dual-model LLM ensemble approach for detecting injections and jailbreaks.
- We will deploy two independent open-weight models via vLLM.
- If there is a disagreement between the two models, we will escalate the event rather than silently picking one.

## Consequences
- **Higher compute cost:** Running two models requires more GPU resources.
- **Better precision:** Ensemble models generally provide higher accuracy and robustness against adversarial attacks.
- **Meta-injection mitigation:** Better defense through structural isolation and diverse model architectures.
