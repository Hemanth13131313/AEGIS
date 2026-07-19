import { useState, useEffect } from 'react';

export interface TraceEvent {
  id: string;
  timestamp: string;
  stage: 'request' | 'pre_filter' | 'policy' | 'scanner_a' | 'scanner_b' | 'ensemble' | 'gateway' | 'upstream';
  label: string;
  details: string;
  latency_ms: number;
  verdict?: 'allow' | 'block' | 'tag' | 'pass' | 'fail';
}

const mockData: Record<string, { events: TraceEvent[], session_id: string, total_latency_ms: number }> = {
  'sess-001': {
    session_id: 'sess-001',
    total_latency_ms: 20,
    events: [
      { id: '1', timestamp: '10:23:01.001', stage: 'request', label: 'REQUEST received', details: 'POST /v1/chat/completions (tokens: ~45)', latency_ms: 0 },
      { id: '2', timestamp: '10:23:01.003', stage: 'pre_filter', label: 'PRE-FILTER', details: 'prompt_injection_evaluator → NO MATCH', latency_ms: 2, verdict: 'pass' },
      { id: '3', timestamp: '10:23:01.004', stage: 'policy', label: 'POLICY CHECK', details: 'default_policy → ALLOW', latency_ms: 1, verdict: 'allow' },
      { id: '4', timestamp: '10:23:01.012', stage: 'scanner_a', label: 'SCANNER A: MistralGuard', details: 'UNSAFE (confidence: 0.94)', latency_ms: 8, verdict: 'block' },
      { id: '5', timestamp: '10:23:01.019', stage: 'scanner_b', label: 'SCANNER B: LlamaGuard', details: 'UNSAFE (confidence: 0.91)', latency_ms: 15, verdict: 'block' },
      { id: '6', timestamp: '10:23:01.019', stage: 'ensemble', label: 'ENSEMBLE', details: 'BLOCK (both models agree, no disagreement)', latency_ms: 0, verdict: 'block' },
      { id: '7', timestamp: '10:23:01.021', stage: 'gateway', label: 'GATEWAY', details: 'REQUEST BLOCKED — POLICY_BLOCKED', latency_ms: 2, verdict: 'block' },
    ]
  },
  'sess-002': {
    session_id: 'sess-002',
    total_latency_ms: 32,
    events: [
      { id: '1', timestamp: '11:05:12.001', stage: 'request', label: 'REQUEST received', details: 'POST /v1/chat/completions', latency_ms: 0 },
      { id: '2', timestamp: '11:05:12.003', stage: 'pre_filter', label: 'PRE-FILTER', details: 'NO MATCH', latency_ms: 2, verdict: 'pass' },
      { id: '3', timestamp: '11:05:12.004', stage: 'policy', label: 'POLICY CHECK', details: 'default_policy → ALLOW', latency_ms: 1, verdict: 'allow' },
      { id: '4', timestamp: '11:05:12.029', stage: 'scanner_a', label: 'SCANNER A', details: 'SAFE', latency_ms: 25, verdict: 'allow' },
      { id: '5', timestamp: '11:05:12.030', stage: 'scanner_b', label: 'SCANNER B', details: 'SAFE', latency_ms: 26, verdict: 'allow' },
      { id: '6', timestamp: '11:05:12.031', stage: 'ensemble', label: 'ENSEMBLE', details: 'ALLOW', latency_ms: 1, verdict: 'allow' },
      { id: '7', timestamp: '11:05:12.033', stage: 'upstream', label: 'UPSTREAM', details: 'Forwarded to upstream LLM', latency_ms: 2, verdict: 'allow' },
    ]
  }
};

export function useSessionTrace(sessionId: string | null) {
  const [data, setData] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!sessionId) {
      setData(null);
      return;
    }

    setIsLoading(true);
    // Simulate network delay
    const timer = setTimeout(() => {
      setData(mockData[sessionId] || null);
      setIsLoading(false);
    }, 400);

    return () => clearTimeout(timer);
  }, [sessionId]);

  return { data, isLoading };
}
