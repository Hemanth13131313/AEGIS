"""
FAISS-based RAG anomaly detector.
Phase 8 upgrade from Phase 4 statistical detector.

Approach: Maintain an online FAISS index of legitimate query embeddings.
On each retrieval event:
1. Compute L2 distance from query embedding to its K nearest neighbors in index
2. If distance exceeds dynamic threshold: anomalous (potential RAG poisoning)
3. Update index with new embedding (online learning)

FAISS not imported at module level — optional dep. Graceful fallback to Phase 4 statistical detector.
"""
from dataclasses import dataclass
from typing import Optional
import structlog

logger = structlog.get_logger(__name__)

try:
    import faiss  # type: ignore
    import numpy as np
    FAISS_AVAILABLE = True
except ImportError:
    FAISS_AVAILABLE = False
    logger.warning("FAISS not available — falling back to statistical RAG monitor")

from app.monitor import RAGMonitor, RetrievalEvent, RAGVerdict

class FAISSRAGMonitor:
    """
    FAISS embedding-neighborhood anomaly detector.
    Falls back to Phase 4 RAGMonitor if FAISS not available or index empty.
    """
    
    def __init__(self, embedding_dim: int = 1536, k_neighbors: int = 5,
                 anomaly_threshold: float = 2.0, min_index_size: int = 100):
        self._dim = embedding_dim
        self._k = k_neighbors
        self._threshold = anomaly_threshold
        self._min_index = min_index_size
        self._fallback = RAGMonitor()
        
        if FAISS_AVAILABLE:
            self._index = faiss.IndexFlatL2(embedding_dim)
            logger.info("FAISS RAG monitor initialized", dim=embedding_dim)
        else:
            self._index = None
    
    async def analyze_retrieval_event(self, event: RetrievalEvent) -> RAGVerdict:
        if not FAISS_AVAILABLE or self._index is None or self._index.ntotal < self._min_index:
            # Fallback to statistical detector until index is populated
            return await self._fallback.analyze_retrieval_event(event)
        
        if not event.query_embedding:
            return await self._fallback.analyze_retrieval_event(event)
        
        query = np.array([event.query_embedding], dtype=np.float32)
        distances, _ = self._index.search(query, self._k)
        avg_distance = float(distances[0].mean())
        
        # Update index with new query embedding
        self._index.add(query)
        
        is_anomalous = avg_distance > self._threshold
        anomaly_score = min(avg_distance / (self._threshold * 2), 1.0)
        
        return RAGVerdict(
            anomaly_score=anomaly_score,
            is_anomalous=is_anomalous,
            reason=f"faiss_avg_distance={avg_distance:.4f} threshold={self._threshold}",
            owasp_category="LLM03" if is_anomalous else "none",
            atlas_technique="AML.T0021" if is_anomalous else "none",
            action="tag" if is_anomalous else "allow"
        )
    
    def index_size(self) -> int:
        return self._index.ntotal if self._index else 0
