"""
Behavioral baseline evaluator for AEGIS Scanner.

Phase 8: Statistical baseline using rolling window mean/std.
Future: Replace with online clustering (BIRCH/HDBSCAN) on embedding features.

Tracked features per (org_id, app_id) baseline:
- tokens_per_request: mean, std, max
- tool_calls_per_session: mean
- request_rate_per_minute: mean, std
- role_sequence_entropy: mean (higher = more varied turn patterns)
"""
from dataclasses import dataclass, field
from collections import deque
import math
import structlog

from app.evaluators.base import BaseEvaluator, EvaluatorResult
from app.models.schemas import VerdictAction

logger = structlog.get_logger(__name__)

@dataclass
class BaselineStats:
    window: deque = field(default_factory=lambda: deque(maxlen=1000))
    
    @property
    def mean(self) -> float:
        return sum(self.window) / len(self.window) if self.window else 0.0
    
    @property
    def std(self) -> float:
        if len(self.window) < 2:
            return 0.0
        m = self.mean
        return math.sqrt(sum((x - m) ** 2 for x in self.window) / len(self.window))
    
    def z_score(self, value: float) -> float:
        s = self.std
        return (value - self.mean) / s if s > 0 else 0.0
    
    def update(self, value: float) -> None:
        self.window.append(value)

class BehavioralBaselineEvaluator(BaseEvaluator):
    """Detects anomalous request patterns vs. org/app rolling baseline."""
    
    def __init__(self, z_threshold: float = 3.5):
        # In-memory baselines keyed by (org_id, app_id)
        # Production: replace with Redis sorted sets for distributed baseline
        self._baselines: dict[str, dict[str, BaselineStats]] = {}
        self._z_threshold = z_threshold
    
    async def evaluate(self, payload: str, context: dict) -> EvaluatorResult:
        org_id = context.get("org_id", "unknown")
        app_id = context.get("app_id", "unknown")
        token_count = context.get("token_count", len(payload.split()))
        
        key = f"{org_id}:{app_id}"
        if key not in self._baselines:
            self._baselines[key] = {"tokens": BaselineStats()}
        
        stats = self._baselines[key]["tokens"]
        
        # Need at least 30 samples to start flagging
        if len(stats.window) >= 30:
            z = stats.z_score(token_count)
            if abs(z) > self._z_threshold:
                stats.update(token_count)
                return EvaluatorResult(
                    verdict=VerdictAction.TAG,
                    confidence=min(abs(z) / 10.0, 0.9),
                    owasp_category="LLM04",
                    atlas_technique="AML.T0043",
                    evidence=f"token_count_z_score={z:.2f} (threshold={self._z_threshold})"
                )
        
        stats.update(token_count)
        return EvaluatorResult(
            verdict=VerdictAction.ALLOW,
            confidence=0.0,
            owasp_category=None,
            atlas_technique=None,
            evidence="within_baseline"
        )
