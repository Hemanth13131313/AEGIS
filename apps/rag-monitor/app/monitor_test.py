import pytest
import asyncio
from app.monitor import RAGMonitor, RetrievalEvent, RAGVerdict

@pytest.fixture
def monitor():
    return RAGMonitor()

@pytest.mark.asyncio
async def test_normal_retrieval_returns_allow(monitor):
    event = RetrievalEvent("s1", [], ["a", "b"], [0.8, 0.75])
    verdict = await monitor.analyze_retrieval_event(event)
    assert verdict.action == "allow"
    assert not verdict.is_anomalous
    assert verdict.reason == "normal"

@pytest.mark.asyncio
async def test_uniform_scores_flagged(monitor):
    event = RetrievalEvent("s1", [], ["a", "b"], [0.5, 0.5])
    verdict = await monitor.analyze_retrieval_event(event)
    assert verdict.is_anomalous
    assert verdict.reason == "uniform_scores"

@pytest.mark.asyncio
async def test_poisoned_top_chunk_flagged(monitor):
    event = RetrievalEvent("s1", [], ["a", "b"], [0.999, 0.2])
    verdict = await monitor.analyze_retrieval_event(event)
    assert verdict.is_anomalous
    assert verdict.reason == "poisoned_top_chunk"

@pytest.mark.asyncio
async def test_empty_retrieval_low_anomaly(monitor):
    event = RetrievalEvent("s1", [], [], [])
    verdict = await monitor.analyze_retrieval_event(event)
    assert not verdict.is_anomalous
    assert verdict.anomaly_score == 0.3
    assert verdict.reason == "empty_retrieval"

@pytest.mark.asyncio
async def test_score_count_mismatch_flagged(monitor):
    event = RetrievalEvent("s1", [], ["a", "b"], [0.8])
    verdict = await monitor.analyze_retrieval_event(event)
    assert verdict.is_anomalous
    assert verdict.reason == "score_count_mismatch"
    assert verdict.action == "block"

@pytest.mark.asyncio
async def test_high_anomaly_score_blocks(monitor):
    # Same as count mismatch which returns 0.95 and blocks
    event = RetrievalEvent("s1", [], ["a"], [0.8, 0.9])
    verdict = await monitor.analyze_retrieval_event(event)
    assert verdict.action == "block"
