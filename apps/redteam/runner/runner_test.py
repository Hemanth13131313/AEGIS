import pytest
import pytest_asyncio
from datetime import datetime, timezone
import httpx

from runner.runner import RedTeamRunner, RunConfig, RunResult
from generators.base import TestCase

@pytest.fixture
def mock_case():
    return TestCase(
        id="RT-LLM01-001",
        name="Test",
        atlas_technique="AML.T0051",
        owasp_category="LLM01",
        prompt="Bad prompt",
        expected_action="block"
    )

class MockResponse:
    def __init__(self, data, status_code=200):
        self._data = data
        self.status_code = status_code
        
    def json(self):
        return self._data
        
    def raise_for_status(self):
        if self.status_code >= 400:
            raise httpx.HTTPStatusError("error", request=None, response=self)

@pytest.mark.asyncio
async def test_run_one_pass(mock_case, monkeypatch):
    async def mock_post(*args, **kwargs):
        return MockResponse({"action": "block", "confidence": 0.95})
        
    monkeypatch.setattr(httpx.AsyncClient, "post", mock_post)
    
    runner = RedTeamRunner(RunConfig(target_url="http://test"))
    
    import asyncio
    async with httpx.AsyncClient() as client:
        result = await runner.run_one(mock_case, asyncio.Semaphore(1), client)
    
    assert result.passed is True
    assert result.action_received == "block"
    assert result.disagreement is False

@pytest.mark.asyncio
async def test_run_one_fail(mock_case, monkeypatch):
    async def mock_post(*args, **kwargs):
        return MockResponse({"action": "allow", "confidence": 0.1})
        
    monkeypatch.setattr(httpx.AsyncClient, "post", mock_post)
    
    runner = RedTeamRunner(RunConfig(target_url="http://test"))
    
    import asyncio
    async with httpx.AsyncClient() as client:
        result = await runner.run_one(mock_case, asyncio.Semaphore(1), client)
    
    assert result.passed is False
    assert result.action_received == "allow"

@pytest.mark.asyncio
async def test_run_one_disagreement(mock_case, monkeypatch):
    async def mock_post(*args, **kwargs):
        return MockResponse({"action": "tag", "confidence": 0.5})
        
    monkeypatch.setattr(httpx.AsyncClient, "post", mock_post)
    
    runner = RedTeamRunner(RunConfig(target_url="http://test"))
    
    import asyncio
    async with httpx.AsyncClient() as client:
        result = await runner.run_one(mock_case, asyncio.Semaphore(1), client)
    
    assert result.passed is True  # Default allows disagreement to pass
    assert result.disagreement is True

@pytest.mark.asyncio
async def test_run_one_timeout(mock_case, monkeypatch):
    async def mock_post(*args, **kwargs):
        raise httpx.TimeoutException("Timeout")
        
    monkeypatch.setattr(httpx.AsyncClient, "post", mock_post)
    
    runner = RedTeamRunner(RunConfig(target_url="http://test"))
    
    import asyncio
    async with httpx.AsyncClient() as client:
        result = await runner.run_one(mock_case, asyncio.Semaphore(1), client)
    
    assert result.passed is False
    assert result.error == "timeout"

def test_summary_pass_rate():
    runner = RedTeamRunner(RunConfig(target_url="http://test"))
    results = [
        RunResult("1", "T", "LLM01", "A", True, "block", "block", False, 1.0, 10, datetime.now()),
        RunResult("2", "T", "LLM01", "A", True, "block", "block", False, 1.0, 10, datetime.now()),
        RunResult("3", "T", "LLM01", "A", False, "allow", "block", False, 1.0, 10, datetime.now())
    ]
    summary = runner.summary(results, datetime.now(), datetime.now())
    assert summary.total == 3
    assert summary.passed == 2
    assert summary.failed == 1
    assert abs(summary.pass_rate - 0.666) < 0.01

def test_summary_by_owasp():
    runner = RedTeamRunner(RunConfig(target_url="http://test"))
    results = [
        RunResult("1", "T", "LLM01", "A", True, "block", "block", False, 1.0, 10, datetime.now()),
        RunResult("2", "T", "LLM06", "A", False, "allow", "block", False, 1.0, 10, datetime.now())
    ]
    summary = runner.summary(results, datetime.now(), datetime.now())
    assert summary.by_owasp["LLM01"]["total"] == 1
    assert summary.by_owasp["LLM01"]["passed"] == 1
    assert summary.by_owasp["LLM06"]["total"] == 1
    assert summary.by_owasp["LLM06"]["failed"] == 1
