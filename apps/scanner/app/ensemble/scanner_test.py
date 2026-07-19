import pytest
import asyncio
from unittest.mock import AsyncMock, patch
from app.ensemble.scanner import EnsembleScanner
from app.models.llm_client import LLMClientConfig, ClassificationResult
from app.models.schemas import VerdictAction

@pytest.fixture
def scanner():
    cfg_a = LLMClientConfig(base_url="http://mock", model="A")
    cfg_b = LLMClientConfig(base_url="http://mock", model="B")
    return EnsembleScanner(cfg_a, cfg_b)

@pytest.mark.asyncio
async def test_both_safe_returns_allow(scanner):
    res_safe = ClassificationResult("SAFE", 0.9, "none", "none", "ok", 10, "A")
    with patch.object(scanner._client_a, 'classify', new_callable=AsyncMock) as m_a, \
         patch.object(scanner._client_b, 'classify', new_callable=AsyncMock) as m_b:
        m_a.return_value = res_safe
        m_b.return_value = ClassificationResult("SAFE", 0.9, "none", "none", "ok", 10, "B")
        
        verdict = await scanner.scan("test", {})
        assert verdict.action == VerdictAction.ALLOW
        assert not verdict.disagreement

@pytest.mark.asyncio
async def test_both_unsafe_returns_block(scanner):
    res_unsafe_a = ClassificationResult("UNSAFE", 0.9, "LLM01", "AML.T0051", "bad", 10, "A")
    res_unsafe_b = ClassificationResult("UNSAFE", 0.8, "LLM01", "AML.T0051", "bad", 10, "B")
    with patch.object(scanner._client_a, 'classify', new_callable=AsyncMock) as m_a, \
         patch.object(scanner._client_b, 'classify', new_callable=AsyncMock) as m_b:
        m_a.return_value = res_unsafe_a
        m_b.return_value = res_unsafe_b
        
        verdict = await scanner.scan("test", {})
        assert verdict.action == VerdictAction.BLOCK
        assert not verdict.disagreement
        assert verdict.confidence == 0.9

@pytest.mark.asyncio
async def test_disagreement_a_unsafe_b_safe_returns_tag_with_disagreement_true(scanner):
    res_unsafe = ClassificationResult("UNSAFE", 0.9, "LLM01", "AML.T0051", "bad", 10, "A")
    res_safe = ClassificationResult("SAFE", 0.9, "none", "none", "ok", 10, "B")
    with patch.object(scanner._client_a, 'classify', new_callable=AsyncMock) as m_a, \
         patch.object(scanner._client_b, 'classify', new_callable=AsyncMock) as m_b:
        m_a.return_value = res_unsafe
        m_b.return_value = res_safe
        
        verdict = await scanner.scan("test", {})
        assert verdict.action == VerdictAction.TAG
        assert verdict.disagreement

@pytest.mark.asyncio
async def test_disagreement_a_safe_b_unsafe_returns_tag_with_disagreement_true(scanner):
    res_safe = ClassificationResult("SAFE", 0.9, "none", "none", "ok", 10, "A")
    res_unsafe = ClassificationResult("UNSAFE", 0.9, "LLM01", "AML.T0051", "bad", 10, "B")
    with patch.object(scanner._client_a, 'classify', new_callable=AsyncMock) as m_a, \
         patch.object(scanner._client_b, 'classify', new_callable=AsyncMock) as m_b:
        m_a.return_value = res_safe
        m_b.return_value = res_unsafe
        
        verdict = await scanner.scan("test", {})
        assert verdict.action == VerdictAction.TAG
        assert verdict.disagreement

@pytest.mark.asyncio
async def test_model_a_exception_does_not_crash_scan(scanner):
    res_safe = ClassificationResult("SAFE", 0.9, "none", "none", "ok", 10, "B")
    with patch.object(scanner._client_a, 'classify', new_callable=AsyncMock) as m_a, \
         patch.object(scanner._client_b, 'classify', new_callable=AsyncMock) as m_b:
        m_a.side_effect = Exception("boom")
        m_b.return_value = res_safe
        
        verdict = await scanner.scan("test", {})
        # A will failsafe to SAFE, B is SAFE -> ALLOW
        assert verdict.action == VerdictAction.ALLOW
        assert verdict.model_a_verdict == "SAFE"
        assert not verdict.disagreement

@pytest.mark.asyncio
async def test_scan_latency_both_called_concurrently(scanner):
    async def slow_classify_a(*args, **kwargs):
        await asyncio.sleep(0.1)
        return ClassificationResult("SAFE", 0.9, "none", "none", "ok", 100, "A")
    async def slow_classify_b(*args, **kwargs):
        await asyncio.sleep(0.1)
        return ClassificationResult("SAFE", 0.9, "none", "none", "ok", 100, "B")
        
    with patch.object(scanner._client_a, 'classify', new=slow_classify_a), \
         patch.object(scanner._client_b, 'classify', new=slow_classify_b):
        verdict = await scanner.scan("test", {})
        # If run sequentially, latency would be > 200ms. Since concurrent, < 150ms.
        assert verdict.latency_ms < 150
