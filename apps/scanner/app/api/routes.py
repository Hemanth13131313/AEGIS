from fastapi import APIRouter, HTTPException, Depends, Request
import asyncio
import time
from app.models.schemas import ScanRequest, ScanVerdict, VerdictAction
from app.evaluators.prompt_injection import PromptInjectionEvaluator
from app.evaluators.jailbreak import JailbreakEvaluator
from app.evaluators.data_exfiltration import DataExfiltrationEvaluator
from app.ensemble.scanner import EnsembleScanner
import structlog

logger = structlog.get_logger(__name__)
router = APIRouter()

@router.post("/scan", response_model=ScanVerdict)
async def scan(request: Request, payload: ScanRequest):
    start_time = time.monotonic()
    
    scanner: EnsembleScanner = request.app.state.scanner
    pre_filters = request.app.state.pre_filters
    
    # Run pre-filters in parallel
    pre_results = await asyncio.gather(
        *[pf.evaluate(payload.payload, payload.context) for pf in pre_filters],
        return_exceptions=True
    )
    
    # Check if any pre-filter blocked
    for res in pre_results:
        if isinstance(res, Exception):
            logger.error("pre_filter.error", error=str(res))
            continue
        if res.verdict == VerdictAction.BLOCK:
            latency_ms = (time.monotonic() - start_time) * 1000
            logger.info("scan.pre_filter.blocked", request_id=payload.request_id, evidence=res.evidence)
            return ScanVerdict(
                session_id=payload.session_id,
                request_id=payload.request_id,
                action=VerdictAction.BLOCK,
                confidence=res.confidence,
                owasp_category=res.owasp_category,
                atlas_technique=res.atlas_technique,
                reason=f"pre_filter_match: {res.evidence}",
                disagreement=False,
                latency_ms=latency_ms,
                pre_filter_triggered=True
            )
            
    # Run ensemble if pre-filters pass
    ensemble_result = await scanner.scan(payload.payload, payload.context)
    latency_ms = (time.monotonic() - start_time) * 1000
    
    logger.info(
        "scan.ensemble.completed",
        request_id=payload.request_id,
        action=ensemble_result.action.value,
        confidence=ensemble_result.confidence,
        latency_ms=latency_ms
    )
    
    return ScanVerdict(
        session_id=payload.session_id,
        request_id=payload.request_id,
        action=ensemble_result.action,
        confidence=ensemble_result.confidence,
        owasp_category=ensemble_result.owasp_category,
        atlas_technique=ensemble_result.atlas_technique,
        reason=ensemble_result.reason,
        disagreement=ensemble_result.disagreement,
        latency_ms=latency_ms,
        pre_filter_triggered=False
    )
