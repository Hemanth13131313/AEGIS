"""
RAG Monitor — detects anomalous retrieval patterns that may indicate
RAG poisoning (LLM03 / AML.T0048) or data exfiltration via retrieval (LLM06).

Phase 4 approach: statistical anomaly detection on retrieval scores.
Phase 7: Replace with embedding neighborhood analysis using FAISS/hnswlib.
"""
from dataclasses import dataclass

import numpy as np
import structlog

logger = structlog.get_logger(__name__)

@dataclass
class RetrievalEvent:
    session_id: str
    query_embedding: list[float]  
    retrieved_chunks: list[str]
    retrieval_scores: list[float]
    top_k: int = 5

@dataclass
class RAGVerdict:
    anomaly_score: float
    is_anomalous: bool
    reason: str
    owasp_category: str
    atlas_technique: str
    action: str

class RAGMonitor:
    def __init__(self, anomaly_threshold: float = 0.7, score_variance_threshold: float = 2.0):
        self.anomaly_threshold = anomaly_threshold
        self.score_variance_threshold = score_variance_threshold
        
    async def analyze_retrieval_event(self, event: RetrievalEvent) -> RAGVerdict:
        scores = event.retrieval_scores
        chunks = event.retrieved_chunks
        
        anomaly_score, reason = self._compute_anomaly_score(scores, chunks)
        
        is_anomalous = anomaly_score >= self.anomaly_threshold
        action = "allow"
        owasp = "none"
        atlas = "none"
        
        if is_anomalous:
            owasp = "LLM03"
            atlas = "AML.T0048"
            action = "tag"
            
        if anomaly_score >= 0.9:
            action = "block"
            
        return RAGVerdict(
            anomaly_score=anomaly_score,
            is_anomalous=is_anomalous,
            reason=reason,
            owasp_category=owasp,
            atlas_technique=atlas,
            action=action
        )
        
    def _compute_anomaly_score(self, scores: list[float], chunks: list[str]) -> tuple[float, str]:
        if not chunks:
            return 0.3, "empty_retrieval"
            
        if len(scores) != len(chunks):
            return 0.95, "score_count_mismatch"
            
        if len(scores) > 1:
            std = np.std(scores)
            if std == 0:
                return 0.85, "uniform_scores"
                
            top_score = scores[0]
            second_score = scores[1]
            if top_score > 0.99 and second_score < 0.5:
                return 0.85, "poisoned_top_chunk"
                
        return 0.0, "normal"
