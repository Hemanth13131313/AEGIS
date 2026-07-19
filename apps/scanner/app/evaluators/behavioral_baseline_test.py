import pytest
from app.evaluators.behavioral_baseline import BehavioralBaselineEvaluator
from app.models.schemas import VerdictAction

@pytest.mark.asyncio
async def test_insufficient_samples_always_allows():
    evaluator = BehavioralBaselineEvaluator()
    for _ in range(10):
        res = await evaluator.evaluate("some text", {"token_count": 100})
        assert res.verdict == VerdictAction.ALLOW

@pytest.mark.asyncio
async def test_normal_value_within_baseline_allows():
    evaluator = BehavioralBaselineEvaluator()
    # build baseline
    for _ in range(40):
        await evaluator.evaluate("text", {"token_count": 100})
    
    res = await evaluator.evaluate("text", {"token_count": 105})
    assert res.verdict == VerdictAction.ALLOW

@pytest.mark.asyncio
async def test_extreme_spike_tags():
    evaluator = BehavioralBaselineEvaluator()
    for _ in range(40):
        await evaluator.evaluate("text", {"token_count": 100})
    
    res = await evaluator.evaluate("text", {"token_count": 10000})
    assert res.verdict == VerdictAction.TAG
    assert res.owasp_category == "LLM04"

@pytest.mark.asyncio
async def test_baseline_updates_on_every_call():
    evaluator = BehavioralBaselineEvaluator()
    for _ in range(35):
        await evaluator.evaluate("text", {"token_count": 100})
        
    stats = evaluator._baselines["unknown:unknown"]["tokens"]
    assert len(stats.window) == 35
