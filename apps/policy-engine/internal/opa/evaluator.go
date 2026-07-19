package opa

import (
	"context"
	"fmt"
	"sync"

	"github.com/open-policy-agent/opa/rego"
	"go.uber.org/zap"
)

type EvalInput struct {
	OrgID        string            `json:"org_id"`
	AppID        string            `json:"app_id"`
	ModelID      string            `json:"model_id"`
	Environment  string            `json:"environment"`
	Auth         AuthInput         `json:"auth"`
	Payload      PayloadInput      `json:"payload"`
	AllowedTools []string          `json:"allowed_tools"`
	Metadata     map[string]string `json:"metadata"`
}

type AuthInput struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id"`
}

type PayloadInput struct {
	TokenCount   int      `json:"token_count"`
	CharsetValid bool     `json:"charset_valid"`
	ToolName     string   `json:"tool_name"`
	TurnCount    int      `json:"turn_count"`
	RoleSequence []string `json:"role_sequence"`
}

type EvalResult struct {
	Allow     bool     `json:"allow"`
	Deny      []string `json:"deny"`
	PolicyRef string   `json:"policy_ref"`
}

type Evaluator struct {
	logger *zap.Logger
	mu     sync.RWMutex
	queries map[string]rego.PreparedEvalQuery
}

func NewEvaluator(logger *zap.Logger) *Evaluator {
	return &Evaluator{
		logger:  logger,
		queries: make(map[string]rego.PreparedEvalQuery),
	}
}

func (e *Evaluator) LoadBundle(ctx context.Context, bundlePath string) error {
	e.logger.Info("LoadBundle not fully implemented yet", zap.String("path", bundlePath))
	return nil
}

func (e *Evaluator) LoadRegoContent(ctx context.Context, policyID, regoContent string) error {
	query, err := rego.New(
		rego.Query("data.aegis.policy"),
		rego.Module(fmt.Sprintf("%s.rego", policyID), regoContent),
	).PrepareForEval(ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare rego content: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.queries[policyID] = query
	e.logger.Info("Loaded policy", zap.String("policy_id", policyID))
	return nil
}

func (e *Evaluator) Evaluate(ctx context.Context, policyID string, input EvalInput) (*EvalResult, error) {
	e.mu.RLock()
	query, ok := e.queries[policyID]
	e.mu.RUnlock()

	if !ok {
		e.logger.Warn("Policy not found", zap.String("policy_id", policyID))
		return &EvalResult{
			Allow:     false,
			Deny:      []string{"POLICY_NOT_FOUND"},
			PolicyRef: policyID,
		}, nil
	}

	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate policy %s: %w", policyID, err)
	}

	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return &EvalResult{Allow: false, Deny: []string{"NO_RESULT"}, PolicyRef: policyID}, nil
	}

	resMap, ok := rs[0].Expressions[0].Value.(map[string]interface{})
	if !ok {
		return &EvalResult{Allow: false, Deny: []string{"INVALID_RESULT"}, PolicyRef: policyID}, nil
	}

	var allow bool
	var deny []string

	for pkgName, pkgData := range resMap {
		if pkgMap, isMap := pkgData.(map[string]interface{}); isMap {
			if a, exists := pkgMap["allow"]; exists {
				if ab, ok := a.(bool); ok && ab {
					allow = true
				}
			}
			if d, exists := pkgMap["deny"]; exists {
				if dSlice, isSlice := d.([]interface{}); isSlice {
					for _, dItem := range dSlice {
						if dStr, ok := dItem.(string); ok {
							deny = append(deny, dStr)
						}
					}
				}
			}
			if !allow {
				e.logger.Debug("Policy package denied", zap.String("pkg", pkgName))
			}
		}
	}

	return &EvalResult{
		Allow:     allow,
		Deny:      deny,
		PolicyRef: policyID,
	}, nil
}
