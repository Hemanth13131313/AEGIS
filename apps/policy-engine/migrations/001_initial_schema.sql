-- Migration: 001_initial_schema
-- Creates the core AEGIS policy and tenancy tables

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Organizations (top of policy hierarchy)
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Applications (owned by org)
CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    environment TEXT NOT NULL DEFAULT 'production', -- production, staging, development
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_applications_org_id ON applications(org_id);

-- Model endpoints (associated with application)
CREATE TABLE model_endpoints (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    app_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    provider TEXT NOT NULL, -- openai, anthropic, vllm, tgi
    endpoint_url TEXT NOT NULL,
    model_digest TEXT, -- sha256 of model weights for supply-chain verification
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_model_endpoints_app_id ON model_endpoints(app_id);

-- Policies (scoped to org, app, model, or environment level)
CREATE TABLE policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    scope_id UUID NOT NULL, -- references org_id, app_id, or model_endpoint_id
    scope_type TEXT NOT NULL CHECK (scope_type IN ('organization', 'application', 'model_endpoint', 'environment')),
    rego_bundle_ref TEXT NOT NULL, -- path/ref to the Rego bundle artifact
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_policies_scope ON policies(scope_id, scope_type);
CREATE INDEX idx_policies_active ON policies(active);

-- Policy versions (immutable audit log of all policy changes)
CREATE TABLE policy_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    rego_content TEXT NOT NULL, -- full Rego content snapshot
    diff TEXT, -- unified diff from previous version
    actor TEXT NOT NULL, -- user or system that made the change
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_policy_versions_policy_id ON policy_versions(policy_id);
CREATE UNIQUE INDEX idx_policy_versions_unique ON policy_versions(policy_id, version);

-- Sessions (AI interaction sessions)
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    app_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    user_ref TEXT, -- external user identifier (hashed/pseudonymized)
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_sessions_app_id ON sessions(app_id);

-- Red-team test case definitions
CREATE TABLE redteam_cases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    technique TEXT NOT NULL, -- MITRE ATLAS technique ID
    owasp_category TEXT NOT NULL, -- OWASP LLM Top 10 ID
    name TEXT NOT NULL,
    prompt_template TEXT NOT NULL,
    expected_action TEXT NOT NULL DEFAULT 'block',
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
