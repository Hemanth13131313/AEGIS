-- ClickHouse Migration: 002_clickhouse_schema
-- High-volume append-only tables for traces, detections, and red-team results

CREATE DATABASE IF NOT EXISTS aegis;

-- Events (all request/response events from the gateway)
CREATE TABLE IF NOT EXISTS aegis.events (
    id UUID,
    session_id UUID,
    event_type LowCardinality(String), -- request, response, tool_call, error
    modality LowCardinality(String) DEFAULT 'text', -- text, image, audio (future)
    payload_redacted String DEFAULT '', -- redacted payload (never raw content at default)
    metadata Map(String, String),
    ts DateTime64(3, 'UTC')
) ENGINE = MergeTree()
ORDER BY (session_id, ts)
PARTITION BY toYYYYMM(ts);

-- Detections (one per detection event, linked to event)
CREATE TABLE IF NOT EXISTS aegis.detections (
    id UUID,
    event_id UUID,
    session_id UUID,
    category LowCardinality(String), -- prompt_injection, jailbreak, rag_poisoning, etc.
    confidence Float32,
    action_taken LowCardinality(String), -- allow, block, redact, tag
    owasp_llm_id LowCardinality(String), -- LLM01..LLM10
    atlas_technique String, -- AML.T0051 etc.
    policy_ref String,
    ts DateTime64(3, 'UTC')
) ENGINE = MergeTree()
ORDER BY (session_id, ts)
PARTITION BY toYYYYMM(ts);

-- Red-team run results
CREATE TABLE IF NOT EXISTS aegis.redteam_runs (
    id UUID,
    case_id UUID,
    technique String,
    owasp_category String,
    result LowCardinality(String), -- pass, fail, error
    action_received String,
    expected_action String,
    evidence String,
    run_at DateTime64(3, 'UTC')
) ENGINE = MergeTree()
ORDER BY (case_id, run_at)
PARTITION BY toYYYYMM(run_at);
