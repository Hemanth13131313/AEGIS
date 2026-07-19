"""
Dual-model ensemble scanner for AEGIS.

Design:
- Two models (A and B) are called INDEPENDENTLY and IN PARALLEL (asyncio.gather)
- No shared context between models — meta-injection prevention
- Disagreement (A!=B on block/allow) → action=tag, disagreement=True
- Both block → action=block
- Both allow → action=allow
- P95 target: ≤500ms (models called concurrently)
"""
import asyncio
import os
import time
from dataclasses import dataclass
from typing import Optional
import structlog

from app.models.llm_client import LLMClassifierClient, LLMClientConfig, ClassificationResult
from app.models.schemas import VerdictAction

logger = structlog.get_logger(__name__)

@dataclass
class EnsembleVerdict:
    action: VerdictAction
    confidence: float
    owasp_category: Optional[str]
    atlas_technique: Optional[str]
    reason: str
    model_a_verdict: str
    model_b_verdict: str
    disagreement: bool
    latency_ms: float

class EnsembleScanner:
    def __init__(self, model_a_config: LLMClientConfig, model_b_config: LLMClientConfig):
        self._client_a = LLMClassifierClient(model_a_config)
        self._client_b = LLMClassifierClient(model_b_config)
        
    async def scan(self, payload: str, context: dict) -> EnsembleVerdict:
        start_time = time.monotonic()
        results = await asyncio.gather(
            self._client_a.classify(payload),
            self._client_b.classify(payload),
            return_exceptions=True
        )
        
        res_a = results[0]
        res_b = results[1]
        
        # Handle exceptions gracefully by producing a failsafe result
        if isinstance(res_a, Exception):
            logger.error("scanner.model_a.exception", error=str(res_a))
            res_a = ClassificationResult("SAFE", 0.5, "none", "none", "exception", 0, self._client_a.config.model)
        if isinstance(res_b, Exception):
            logger.error("scanner.model_b.exception", error=str(res_b))
            res_b = ClassificationResult("SAFE", 0.5, "none", "none", "exception", 0, self._client_b.config.model)
            
        disagreement = res_a.verdict != res_b.verdict
        
        if disagreement:
            action = VerdictAction.TAG
            confidence = max(res_a.confidence, res_b.confidence)
            reason = "ensemble_disagreement"
            unsafe_res = res_a if res_a.verdict == "UNSAFE" else res_b
            owasp_category = unsafe_res.category
            atlas_technique = unsafe_res.atlas_technique
        elif res_a.verdict == "UNSAFE":
            action = VerdictAction.BLOCK
            confidence = max(res_a.confidence, res_b.confidence)
            reason = "ensemble_block"
            # Pick highest confidence category
            unsafe_res = res_a if res_a.confidence >= res_b.confidence else res_b
            owasp_category = unsafe_res.category
            atlas_technique = unsafe_res.atlas_technique
        else:
            action = VerdictAction.ALLOW
            confidence = min(res_a.confidence, res_b.confidence)
            reason = "ensemble_allow"
            owasp_category = None
            atlas_technique = None

        latency_ms = (time.monotonic() - start_time) * 1000
        
        logger.info(
            "scanner.ensemble.result",
            action=action.value,
            disagreement=disagreement,
            latency_ms=latency_ms,
            model_a_verdict=res_a.verdict,
            model_b_verdict=res_b.verdict
        )
        
        return EnsembleVerdict(
            action=action,
            confidence=confidence,
            owasp_category=owasp_category,
            atlas_technique=atlas_technique,
            reason=reason,
            model_a_verdict=res_a.verdict,
            model_b_verdict=res_b.verdict,
            disagreement=disagreement,
            latency_ms=latency_ms
        )

    @classmethod
    def from_env(cls) -> 'EnsembleScanner':
        api_key = os.getenv("AEGIS_SCANNER_API_KEY", "")
        
        url_a = os.getenv("AEGIS_SCANNER_MODEL_A_URL", "http://localhost:8000/v1")
        name_a = os.getenv("AEGIS_SCANNER_MODEL_A_NAME", "mistral-7b-instruct")
        
        url_b = os.getenv("AEGIS_SCANNER_MODEL_B_URL", "http://localhost:8000/v1")
        name_b = os.getenv("AEGIS_SCANNER_MODEL_B_NAME", "llama-3-8b-instruct")
        
        config_a = LLMClientConfig(base_url=url_a, model=name_a, api_key=api_key)
        config_b = LLMClientConfig(base_url=url_b, model=name_b, api_key=api_key)
        
        return cls(config_a, config_b)
