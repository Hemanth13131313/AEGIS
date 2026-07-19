package hierarchy

import (
	"context"
	"fmt"

	"github.com/aegis-security/aegis/apps/policy-engine/internal/store"
	"go.uber.org/zap"
)

type ScopeLevel string

const (
	Organization  ScopeLevel = "organization"
	Application   ScopeLevel = "application"
	ModelEndpoint ScopeLevel = "model_endpoint"
	Environment   ScopeLevel = "environment"
)

type ResolvedPolicy struct {
	PolicyID    string
	ScopeLevel  ScopeLevel
	RegoContent string
	PolicyRef   string
}

type Resolver struct {
	store  store.PolicyStoreInterface
	logger *zap.Logger
}

func NewResolver(s store.PolicyStoreInterface, logger *zap.Logger) *Resolver {
	return &Resolver{
		store:  s,
		logger: logger,
	}
}

// Resolve fetches active policies ordered from broadest (org) to most specific (env).
func (r *Resolver) Resolve(ctx context.Context, orgID, appID, modelID, env string) ([]*ResolvedPolicy, error) {
	var resolved []*ResolvedPolicy

	scopes := []struct {
		id    string
		level ScopeLevel
	}{
		{orgID, Organization},
		{appID, Application},
		{modelID, ModelEndpoint},
		{env, Environment},
	}

	for _, s := range scopes {
		if s.id == "" {
			continue
		}
		pols, err := r.store.GetActiveByScope(ctx, s.id, string(s.level))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch scope %s %s: %w", s.level, s.id, err)
		}
		for _, p := range pols {
			resolved = append(resolved, &ResolvedPolicy{
				PolicyID:   p.ID,
				ScopeLevel: s.level,
				PolicyRef:  p.RegoBundleRef,
			})
		}
	}

	return resolved, nil
}

func MostSpecific(policies []*ResolvedPolicy) *ResolvedPolicy {
	if len(policies) == 0 {
		return nil
	}
	// Since Resolve orders them org->app->model->env, the last is the most specific.
	return policies[len(policies)-1]
}
