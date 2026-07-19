"""
Rule-based prompt injection pre-filter.
Runs BEFORE LLM ensemble to catch obvious patterns instantly (<1ms).
Reduces LLM calls for clear-cut cases.
"""
from dataclasses import dataclass
from typing import Optional
import re
import structlog

from app.evaluators.base import BaseEvaluator, EvaluatorResult
from app.models.schemas import VerdictAction

logger = structlog.get_logger(__name__)

@dataclass
class PatternConfig:
    pattern: re.Pattern
    owasp: str
    atlas: str
    name: str

class PromptInjectionPatterns:
    PATTERNS = [
        PatternConfig(re.compile(r"ignore previous", re.I), "LLM01", "AML.T0051", "ignore_instructions"),
        PatternConfig(re.compile(r"forget your instructions", re.I), "LLM01", "AML.T0051", "forget_instructions"),
        PatternConfig(re.compile(r"disregard", re.I), "LLM01", "AML.T0051", "disregard_instructions"),
        PatternConfig(re.compile(r"you are now", re.I), "LLM01", "AML.T0051", "role_escape"),
        PatternConfig(re.compile(r"act as", re.I), "LLM01", "AML.T0051", "role_escape"),
        PatternConfig(re.compile(r"pretend you are", re.I), "LLM01", "AML.T0051", "role_escape"),
        PatternConfig(re.compile(r"your new persona", re.I), "LLM01", "AML.T0051", "role_escape"),
        PatternConfig(re.compile(r"\bDAN\b", re.I), "LLM01", "AML.T0051", "jailbreak_dan"),
        PatternConfig(re.compile(r"Do Anything Now", re.I), "LLM01", "AML.T0051", "jailbreak_dan"),
        PatternConfig(re.compile(r"developer mode", re.I), "LLM01", "AML.T0051", "jailbreak_dev"),
        PatternConfig(re.compile(r"god mode", re.I), "LLM01", "AML.T0051", "jailbreak_god"),
        PatternConfig(re.compile(r"\[\[INJECT\]\]", re.I), "LLM01", "AML.T0051", "indirect_inject"),
        PatternConfig(re.compile(r"<\|SYSTEM\|>", re.I), "LLM01", "AML.T0051", "system_prompt_leak"),
        PatternConfig(re.compile(r"print your system prompt", re.I), "LLM01", "AML.T0051", "exfiltration_prompt"),
        PatternConfig(re.compile(r"reveal your instructions", re.I), "LLM01", "AML.T0051", "exfiltration_prompt"),
        PatternConfig(re.compile(r"what were you told", re.I), "LLM01", "AML.T0051", "exfiltration_prompt")
    ]

class PromptInjectionEvaluator(BaseEvaluator):
    async def evaluate(self, payload: str, context: dict) -> EvaluatorResult:
        for p in PromptInjectionPatterns.PATTERNS:
            if p.pattern.search(payload):
                logger.info("prefilter.prompt_injection.match", name=p.name)
                return EvaluatorResult(
                    verdict=VerdictAction.BLOCK,
                    confidence=0.95,
                    evidence=p.name,
                    owasp_category=p.owasp,
                    atlas_technique=p.atlas
                )
        return EvaluatorResult(
            verdict=VerdictAction.ALLOW,
            confidence=0.0,
            evidence="no_pattern_match"
        )
