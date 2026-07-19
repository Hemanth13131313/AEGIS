from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import Optional
from app.models.schemas import VerdictAction

@dataclass
class EvaluatorResult:
    verdict: VerdictAction
    confidence: float
    evidence: str
    owasp_category: Optional[str] = None
    atlas_technique: Optional[str] = None

class BaseEvaluator(ABC):
    @abstractmethod
    async def evaluate(self, payload: str, context: dict) -> EvaluatorResult:
        pass
