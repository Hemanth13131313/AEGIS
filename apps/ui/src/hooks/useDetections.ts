import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { formatDistanceToNow } from 'date-fns';

export interface Detection {
  id: string;
  session_id: string;
  category: string;
  owasp_category: string;
  atlas_technique: string;
  action_taken: 'allow' | 'block' | 'redact' | 'tag';
  confidence: number;
  timestamp: string;
  severity: 'critical' | 'high' | 'medium' | 'low' | 'safe';
  policy_ref: string;
}

export interface DetectionFilters {
  severity?: string;
  owasp_category?: string;
  search?: string;
  cursor?: string;
  limit?: number;
  live?: boolean;
}

const MOCK_DETECTIONS: Detection[] = [
  {
    id: 'det_1', session_id: 'sess_1ab2c3d4', category: 'prompt_injection',
    owasp_category: 'LLM01', atlas_technique: 'AML.T0051', action_taken: 'block',
    confidence: 0.95, timestamp: new Date(Date.now() - 1000 * 60 * 2).toISOString(),
    severity: 'critical', policy_ref: 'policy-sec-1'
  },
  {
    id: 'det_2', session_id: 'sess_9f8e7d6c', category: 'rag_poisoning',
    owasp_category: 'LLM03', atlas_technique: 'AML.T0000', action_taken: 'block',
    confidence: 0.85, timestamp: new Date(Date.now() - 1000 * 60 * 15).toISOString(),
    severity: 'high', policy_ref: 'policy-rag-1'
  },
  {
    id: 'det_3', session_id: 'sess_4x5y6z7w', category: 'jailbreak',
    owasp_category: 'LLM01', atlas_technique: 'AML.T0054', action_taken: 'block',
    confidence: 0.92, timestamp: new Date(Date.now() - 1000 * 60 * 45).toISOString(),
    severity: 'critical', policy_ref: 'policy-sec-2'
  },
  {
    id: 'det_4', session_id: 'sess_11223344', category: 'data_exfiltration',
    owasp_category: 'LLM06', atlas_technique: 'AML.T0056', action_taken: 'redact',
    confidence: 0.65, timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
    severity: 'medium', policy_ref: 'policy-privacy-1'
  },
  {
    id: 'det_5', session_id: 'sess_aabbccdd', category: 'safe_request',
    owasp_category: 'none', atlas_technique: 'none', action_taken: 'allow',
    confidence: 0.1, timestamp: new Date(Date.now() - 1000 * 60 * 60 * 5).toISOString(),
    severity: 'safe', policy_ref: 'policy-base'
  }
];

export function useDetections(filters: DetectionFilters = {}) {
  return useQuery({
    queryKey: ['detections', filters],
    queryFn: async () => {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 300));
      return {
        detections: MOCK_DETECTIONS,
        next_cursor: 'cursor_xyz'
      };
    },
    refetchInterval: filters.live ? 5000 : false,
  });
}

export function useDetectionSummary() {
  return useQuery({
    queryKey: ['detection_summary'],
    queryFn: async () => {
      await new Promise(resolve => setTimeout(resolve, 200));
      return {
        total_today: 47,
        active_policies: 12,
        endpoints: 3,
        redteam_coverage: 60
      };
    }
  });
}
