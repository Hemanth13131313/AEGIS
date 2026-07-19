# ADR 0007: RAG Anomaly Detection

## Status
Partially accepted (Phase 4 statistical; Phase 7 embedding)

## Context
RAG poisoning (LLM03 / AML.T0048) and data exfiltration via retrieval manipulation (LLM06) are emerging attack vectors. We need a way to detect anomalous retrieval patterns without adding significant latency or disrupting normal retrieval flows.

## Decision
- **Phase 4**: Implement statistical anomaly detection based on retrieval score distributions. This approach is fast (<5ms) and can detect obvious anomalies such as uniformly high scores (adversarial injection) or sudden score drop-offs.
- **Phase 7**: Implement embedding neighborhood analysis using FAISS or hnswlib to detect semantically anomalous queries compared to known safe patterns.

## Consequences
- Phase 4 will have false negatives for highly sophisticated, targeted attacks that craft context-aware prompts.
- Phase 7 will significantly improve recall but introduces infrastructure dependencies (vector database/FAISS) and memory overhead.
