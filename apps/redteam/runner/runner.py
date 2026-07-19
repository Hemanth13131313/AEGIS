"""
AEGIS Red Team Runner.
Submits test cases to the AEGIS /scan endpoint and reports results.
"""
import asyncio
import time
import httpx
import uuid
from dataclasses import dataclass
from datetime import datetime, timezone
from typing import Optional, Dict
import structlog

from generators.base import TestCase

logger = structlog.get_logger(__name__)

@dataclass
class RunConfig:
    target_url: str
    api_key: str = ""
    concurrency: int = 5
    timeout_seconds: float = 10.0
    fail_on_disagreement: bool = False
    dry_run: bool = False

@dataclass
class RunResult:
    test_case_id: str
    test_case_name: str
    owasp_category: str
    atlas_technique: str
    passed: bool
    action_received: str
    expected_action: str
    disagreement: bool
    confidence: float
    latency_ms: float
    run_at: datetime
    error: Optional[str] = None

@dataclass
class RunSummary:
    total: int
    passed: int
    failed: int
    errors: int
    disagreements: int
    pass_rate: float
    by_owasp: Dict[str, Dict[str, int]]
    by_atlas: Dict[str, Dict[str, int]]
    run_at: datetime
    duration_seconds: float

class RedTeamRunner:
    def __init__(self, config: RunConfig):
        self.config = config
    
    async def run_all(self, test_cases: list[TestCase]) -> list[RunResult]:
        semaphore = asyncio.Semaphore(self.config.concurrency)
        
        async with httpx.AsyncClient(timeout=self.config.timeout_seconds) as client:
            tasks = [self.run_one(case, semaphore, client) for case in test_cases]
            results = await asyncio.gather(*tasks)
            
        # Sort failures first, then passed
        return sorted(results, key=lambda r: (r.passed, r.test_case_id))

    async def run_one(self, case: TestCase, semaphore: asyncio.Semaphore, client: httpx.AsyncClient) -> RunResult:
        async with semaphore:
            start_time = time.monotonic()
            run_at = datetime.now(timezone.utc)
            
            if self.config.dry_run:
                # Mock result for dry run
                await asyncio.sleep(0.01)
                return RunResult(
                    test_case_id=case.id,
                    test_case_name=case.name,
                    owasp_category=case.owasp_category,
                    atlas_technique=case.atlas_technique,
                    passed=True,
                    action_received=case.expected_action,
                    expected_action=case.expected_action,
                    disagreement=False,
                    confidence=0.99,
                    latency_ms=10.0,
                    run_at=run_at
                )
            
            headers = {"Content-Type": "application/json"}
            if self.config.api_key:
                headers["Authorization"] = f"Bearer {self.config.api_key}"
                
            payload = {
                "session_id": str(uuid.uuid4()),
                "request_id": str(uuid.uuid4()),
                "payload": case.prompt,
                "context": {"redteam": True, "case_id": case.id}
            }
            
            try:
                response = await client.post(f"{self.config.target_url}/api/v1/scan", json=payload, headers=headers)
                response.raise_for_status()
                data = response.json()
                action_received = data.get("action", "allow")
                confidence = data.get("confidence", 0.0)
                
                passed = action_received == case.expected_action
                disagreement = False
                
                if case.expected_action == "block" and action_received == "tag":
                    disagreement = True
                    if self.config.fail_on_disagreement:
                        passed = False
                    else:
                        passed = True  # Partial detection might be acceptable unless strict
                
                latency_ms = (time.monotonic() - start_time) * 1000
                
                return RunResult(
                    test_case_id=case.id,
                    test_case_name=case.name,
                    owasp_category=case.owasp_category,
                    atlas_technique=case.atlas_technique,
                    passed=passed,
                    action_received=action_received,
                    expected_action=case.expected_action,
                    disagreement=disagreement,
                    confidence=confidence,
                    latency_ms=latency_ms,
                    run_at=run_at
                )
            except httpx.TimeoutException:
                return self._error_result(case, "timeout", start_time, run_at)
            except Exception as e:
                return self._error_result(case, str(e), start_time, run_at)

    def _error_result(self, case: TestCase, error: str, start_time: float, run_at: datetime) -> RunResult:
        latency_ms = (time.monotonic() - start_time) * 1000
        return RunResult(
            test_case_id=case.id,
            test_case_name=case.name,
            owasp_category=case.owasp_category,
            atlas_technique=case.atlas_technique,
            passed=False,
            action_received="error",
            expected_action=case.expected_action,
            disagreement=False,
            confidence=0.0,
            latency_ms=latency_ms,
            run_at=run_at,
            error=error
        )

    def summary(self, results: list[RunResult], start_time: datetime, end_time: datetime) -> RunSummary:
        total = len(results)
        passed = sum(1 for r in results if r.passed and not r.error)
        failed = sum(1 for r in results if not r.passed and not r.error)
        errors = sum(1 for r in results if r.error)
        disagreements = sum(1 for r in results if r.disagreement)
        
        pass_rate = passed / total if total > 0 else 0.0
        
        by_owasp: Dict[str, Dict[str, int]] = {}
        by_atlas: Dict[str, Dict[str, int]] = {}
        
        for r in results:
            if r.owasp_category not in by_owasp:
                by_owasp[r.owasp_category] = {"total": 0, "passed": 0, "failed": 0}
            by_owasp[r.owasp_category]["total"] += 1
            if r.passed:
                by_owasp[r.owasp_category]["passed"] += 1
            elif not r.error:
                by_owasp[r.owasp_category]["failed"] += 1
                
            if r.atlas_technique not in by_atlas:
                by_atlas[r.atlas_technique] = {"total": 0, "passed": 0, "failed": 0}
            by_atlas[r.atlas_technique]["total"] += 1
            if r.passed:
                by_atlas[r.atlas_technique]["passed"] += 1
            elif not r.error:
                by_atlas[r.atlas_technique]["failed"] += 1
                
        duration_seconds = (end_time - start_time).total_seconds()
        
        return RunSummary(
            total=total,
            passed=passed,
            failed=failed,
            errors=errors,
            disagreements=disagreements,
            pass_rate=pass_rate,
            by_owasp=by_owasp,
            by_atlas=by_atlas,
            run_at=start_time,
            duration_seconds=duration_seconds
        )
