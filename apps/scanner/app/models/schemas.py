from pydantic import BaseModel, Field, field_validator
from enum import Enum
from typing import Optional

class VerdictAction(str, Enum):
    ALLOW = "allow"
    BLOCK = "block" 
    REDACT = "redact"
    TAG = "tag"

class OWASPCategory(str, Enum):
    LLM01 = "LLM01"  # Prompt Injection
    LLM02 = "LLM02"  # Insecure Output Handling
    LLM03 = "LLM03"  # Training Data Poisoning
    LLM04 = "LLM04"  # Model Denial of Service
    LLM05 = "LLM05"  # Supply Chain Vulnerabilities
    LLM06 = "LLM06"  # Sensitive Information Disclosure
    LLM07 = "LLM07"  # Insecure Plugin Design
    LLM08 = "LLM08"  # Excessive Agency
    LLM09 = "LLM09"  # Overreliance
    LLM10 = "LLM10"  # Model Theft

class MITREAtlasTechnique(str, Enum):
    T0051 = "AML.T0051"  # LLM Prompt Injection
    T0054 = "AML.T0054"  # LLM Data Leakage
    T0048 = "AML.T0048"  # Backdoor ML Model
    T0043 = "AML.T0043"  # Craft Adversarial Data
    T0040 = "AML.T0040"  # ML Model Inference API Access

class ScanRequest(BaseModel):
    session_id: str
    request_id: str
    payload: str = Field(..., max_length=32768)  # 32K char limit
    context: dict = Field(default_factory=dict)
    
    @field_validator('payload')
    @classmethod
    def payload_not_empty(cls, v: str) -> str:
        if not v.strip():
            raise ValueError('payload cannot be empty')
        return v

class ScanVerdict(BaseModel):
    session_id: str
    request_id: str
    action: VerdictAction
    confidence: float = Field(..., ge=0.0, le=1.0)
    owasp_category: Optional[str] = None
    atlas_technique: Optional[str] = None
    reason: str
    disagreement: bool = False
    latency_ms: float
    pre_filter_triggered: bool = False  # true if rule-based filter caught it

class EnsembleResult(BaseModel):
    action: VerdictAction
    confidence: float
    owasp_category: Optional[str] = None
    atlas_technique: Optional[str] = None
    model_a_verdict: str
    model_b_verdict: str
    disagreement: bool
