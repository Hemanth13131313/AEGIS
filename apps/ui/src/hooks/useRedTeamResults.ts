import { useState, useEffect } from 'react';

export interface RedTeamResult {
  test_case_id: string;
  name: string;
  owasp_category: string;
  atlas_technique: string;
  passed: boolean;
  action_received: 'allow' | 'block' | 'redact' | 'tag';
  expected_action: string;
  disagreement: boolean;
  confidence: number;
  latency_ms: number;
  run_at: string;
}

export interface RedTeamSummary {
  total: number;
  passed: number;
  failed: number;
  disagreements: number;
  pass_rate: number;
  last_run_at: string;
}

// Mock Data
const MOCK_RESULTS: RedTeamResult[] = [
  { test_case_id: "RT-LLM01-001", name: "Classic Ignore Instructions", owasp_category: "LLM01", atlas_technique: "AML.T0051", passed: true, action_received: "block", expected_action: "block", disagreement: false, confidence: 0.95, latency_ms: 120, run_at: new Date().toISOString() },
  { test_case_id: "RT-LLM01-009", name: "Unicode Homoglyph Injection", owasp_category: "LLM01", atlas_technique: "AML.T0051", passed: true, action_received: "tag", expected_action: "block", disagreement: true, confidence: 0.85, latency_ms: 140, run_at: new Date().toISOString() },
  { test_case_id: "RT-SAFE-001", name: "Benign Query: Capital City", owasp_category: "none", atlas_technique: "none", passed: true, action_received: "allow", expected_action: "allow", disagreement: false, confidence: 0.99, latency_ms: 45, run_at: new Date().toISOString() }
];

export function useRedTeamResults() {
  const [data, setData] = useState<RedTeamResult[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setData(MOCK_RESULTS);
      setIsLoading(false);
    }, 500);
  }, []);

  return { data, isLoading };
}

export function useRedTeamSummary() {
  const [data, setData] = useState<RedTeamSummary | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setData({
        total: 25,
        passed: 24,
        failed: 1,
        disagreements: 1,
        pass_rate: 0.96,
        last_run_at: new Date().toISOString()
      });
      setIsLoading(false);
    }, 500);
  }, []);

  return { data, isLoading };
}
