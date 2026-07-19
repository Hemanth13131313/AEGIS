# ADR 0017: FAISS RAG Anomaly Detection (Phase 8 upgrade)

## Status
Accepted (Phase 8 online, replacing Phase 4 statistical)

## Context
The statistical score analysis introduced in Phase 4 has a high false-negative rate for sophisticated RAG poisoning, as it only looks at retrieval relevance scores rather than the semantic clustering of queries and retrieved chunks.

## Decision
Upgrade to a FAISS embedding-neighborhood approach. Maintain an online FAISS index of legitimate query embeddings. Evaluate new queries via L2 nearest-neighbor distance against the index with a dynamic threshold. Update the index with new valid embeddings dynamically.

## Consequences
- Requires configuring an embedding dimension.
- FAISS becomes an optional dependency for the RAG monitor.
- Implemented with graceful fallback to the Phase 4 statistical detector if FAISS is not installed or the index is unpopulated.
