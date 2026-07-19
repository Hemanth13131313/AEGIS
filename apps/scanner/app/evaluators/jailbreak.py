"""
Jailbreak evaluator (rule-based pre-filter).
"""
import re
import structlog
from dataclasses import dataclass
from app.evaluators.base import BaseEvaluator, EvaluatorResult
from app.models.schemas import VerdictAction

logger = structlog.get_logger(__name__)

@dataclass
class PatternConfig:
    pattern: re.Pattern
    owasp: str
    atlas: str
    name: str

class JailbreakPatterns:
    PATTERNS = [
        PatternConfig(re.compile(r"in your last response you said", re.I), "LLM01", "AML.T0051", "multi_turn_confusion"),
        PatternConfig(re.compile(r"you previously told me", re.I), "LLM01", "AML.T0051", "multi_turn_confusion"),
        PatternConfig(re.compile(r"OpenAI says", re.I), "LLM01", "AML.T0051", "authority_spoofing"),
        PatternConfig(re.compile(r"Anthropic permits", re.I), "LLM01", "AML.T0051", "authority_spoofing"),
        PatternConfig(re.compile(r"your developer allows", re.I), "LLM01", "AML.T0051", "authority_spoofing"),
        PatternConfig(re.compile(r"base64 decode", re.I), "LLM01", "AML.T0051", "encoding_tricks"),
        PatternConfig(re.compile(r"```SYSTEM", re.I), "LLM01", "AML.T0051", "prompt_delimiters"),
        PatternConfig(re.compile(r"#### INSTRUCTIONS", re.I), "LLM01", "AML.T0051", "prompt_delimiters"),
        PatternConfig(re.compile(r"---OVERRIDE---", re.I), "LLM01", "AML.T0051", "prompt_delimiters"),
    ]

class JailbreakEvaluator(BaseEvaluator):
    async def evaluate(self, payload: str, context: dict) -> EvaluatorResult:
        for p in JailbreakPatterns.PATTERNS:
            if p.pattern.search(payload):
                logger.info("prefilter.jailbreak.match", name=p.name)
                return EvaluatorResult(
                    verdict=VerdictAction.BLOCK,
                    confidence=0.90,
                    evidence=p.name,
                    owasp_category=p.owasp,
                    atlas_technique=p.atlas
                )
        return EvaluatorResult(
            verdict=VerdictAction.ALLOW,
            confidence=0.0,
            evidence="no_pattern_match"
        )
