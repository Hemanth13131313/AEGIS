package opa

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

const defaultPolicy = `
package aegis.policy.default

import rego.v1

default allow := false

allow if {
    input.auth.valid == true
    not exceeds_token_limit
    not invalid_charset
    not invalid_role_sequence
}

deny contains "UNAUTHENTICATED" if {
    input.auth.valid != true
}
deny contains "TOKEN_LIMIT_EXCEEDED" if {
    exceeds_token_limit
}
deny contains "INVALID_CHARSET" if {
    invalid_charset
}
deny contains "INVALID_ROLE_SEQUENCE" if {
    invalid_role_sequence
}
deny contains "TOOL_NOT_ALLOWED" if {
    input.payload.tool_name != ""
    not tool_allowed
}

max_tokens := 8192

exceeds_token_limit if {
    input.payload.token_count > max_tokens
}
invalid_charset if {
    input.payload.charset_valid == false
}
valid_first_roles := {"system", "user"}
invalid_role_sequence if {
    count(input.payload.role_sequence) > 0
    not (input.payload.role_sequence[0] in valid_first_roles)
}
tool_allowed if {
    input.payload.tool_name in input.allowed_tools
}
`

func setupEvaluator(t *testing.T) *Evaluator {
	logger := zap.NewNop()
	e := NewEvaluator(logger)
	err := e.LoadRegoContent(context.Background(), "default_policy", defaultPolicy)
	if err != nil {
		t.Fatalf("failed to load policy: %v", err)
	}
	return e
}

func TestEvaluate_DefaultPolicy_AllowsValidInput(t *testing.T) {
	e := setupEvaluator(t)
	input := EvalInput{
		Auth: AuthInput{Valid: true},
		Payload: PayloadInput{
			TokenCount:   100,
			CharsetValid: true,
			RoleSequence: []string{"user"},
		},
	}
	res, err := e.Evaluate(context.Background(), "default_policy", input)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if !res.Allow {
		t.Errorf("Expected allow=true, got false. Deny: %v", res.Deny)
	}
}

func TestEvaluate_DefaultPolicy_DeniesExcessTokens(t *testing.T) {
	e := setupEvaluator(t)
	input := EvalInput{
		Auth: AuthInput{Valid: true},
		Payload: PayloadInput{
			TokenCount:   9000,
			CharsetValid: true,
			RoleSequence: []string{"user"},
		},
	}
	res, err := e.Evaluate(context.Background(), "default_policy", input)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if res.Allow {
		t.Errorf("Expected allow=false, got true")
	}
	found := false
	for _, d := range res.Deny {
		if d == "TOKEN_LIMIT_EXCEEDED" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected TOKEN_LIMIT_EXCEEDED in deny, got %v", res.Deny)
	}
}

func TestEvaluate_DefaultPolicy_DeniesDisallowedTool(t *testing.T) {
	e := setupEvaluator(t)
	input := EvalInput{
		Auth: AuthInput{Valid: true},
		Payload: PayloadInput{
			TokenCount:   100,
			CharsetValid: true,
			RoleSequence: []string{"user"},
			ToolName:     "shell_exec",
		},
		AllowedTools: []string{"web_search"},
	}
	res, err := e.Evaluate(context.Background(), "default_policy", input)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if res.Allow {
		t.Errorf("Expected allow=false for disallowed tool")
	}
	found := false
	for _, d := range res.Deny {
		if d == "TOOL_NOT_ALLOWED" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected TOOL_NOT_ALLOWED in deny, got %v", res.Deny)
	}
}

func TestEvaluate_UnknownPolicy(t *testing.T) {
	e := setupEvaluator(t)
	input := EvalInput{}
	res, err := e.Evaluate(context.Background(), "unknown", input)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if res.Allow {
		t.Errorf("Expected allow=false")
	}
	if len(res.Deny) == 0 || res.Deny[0] != "POLICY_NOT_FOUND" {
		t.Errorf("Expected POLICY_NOT_FOUND, got %v", res.Deny)
	}
}
