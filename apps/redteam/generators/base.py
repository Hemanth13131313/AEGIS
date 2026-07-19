from abc import ABC, abstractmethod
from dataclasses import dataclass, field

@dataclass
class TestCase:
    id: str
    name: str
    atlas_technique: str
    owasp_category: str
    prompt: str
    expected_action: str  # allow, block, redact, tag
    version: int = 1
    tags: list[str] = field(default_factory=list)
    meta: dict = field(default_factory=dict)  # extra metadata (e.g., variant_of)

class BaseGenerator(ABC):
    @abstractmethod
    def generate(self, count: int) -> list[TestCase]:
        ...
    
    def generate_id(self, owasp: str, seq: int) -> str:
        return f"RT-{owasp}-{seq:03d}"
