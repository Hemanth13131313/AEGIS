package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Policy struct {
	ID            string    `json:"id"`
	ScopeID       string    `json:"scope_id"`
	ScopeType     string    `json:"scope_type"`
	RegoBundleRef string    `json:"rego_bundle_ref"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PolicyVersion struct {
	ID          string    `json:"id"`
	PolicyID    string    `json:"policy_id"`
	Version     int       `json:"version"`
	RegoContent string    `json:"rego_content"`
	Diff        string    `json:"diff"`
	Actor       string    `json:"actor"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreatePolicyInput struct {
	ScopeID     string `json:"scope_id"`
	ScopeType   string `json:"scope_type"`
	RegoContent string `json:"rego_content"`
	Actor       string `json:"actor"`
}

type PolicyStoreInterface interface {
	Create(ctx context.Context, input CreatePolicyInput) (*Policy, error)
	GetByID(ctx context.Context, id string) (*Policy, error)
	GetActiveByScope(ctx context.Context, scopeID, scopeType string) ([]*Policy, error)
	ListPolicies(ctx context.Context, cursor string, limit int) ([]*Policy, string, error)
	CreateVersion(ctx context.Context, policyID, regoContent, diff, actor string) (*PolicyVersion, error)
	GetVersions(ctx context.Context, policyID string) ([]*PolicyVersion, error)
}

type PolicyStore struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPolicyStore(pool *pgxpool.Pool, logger *zap.Logger) *PolicyStore {
	return &PolicyStore{
		pool:   pool,
		logger: logger,
	}
}

func (s *PolicyStore) Create(ctx context.Context, input CreatePolicyInput) (*Policy, error) {
	if s.pool == nil {
		return &Policy{ID: uuid.NewString(), ScopeID: input.ScopeID, ScopeType: input.ScopeType, RegoBundleRef: "mock"}, nil
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	polID := uuid.NewString()
	pol := &Policy{
		ID:            polID,
		ScopeID:       input.ScopeID,
		ScopeType:     input.ScopeType,
		RegoBundleRef: "",
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO policies (id, scope_id, scope_type, rego_bundle_ref, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		pol.ID, pol.ScopeID, pol.ScopeType, pol.RegoBundleRef, pol.Active, pol.CreatedAt, pol.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert policy: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO policy_versions (id, policy_id, version, rego_content, diff, actor, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		uuid.NewString(), pol.ID, 1, input.RegoContent, "", input.Actor, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to insert policy version: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return pol, nil
}

func (s *PolicyStore) GetByID(ctx context.Context, id string) (*Policy, error) {
	if s.pool == nil {
		return nil, fmt.Errorf("pool not initialized")
	}
	pol := &Policy{}
	err := s.pool.QueryRow(ctx, `
		SELECT id, scope_id, scope_type, rego_bundle_ref, active, created_at, updated_at
		FROM policies WHERE id = $1`, id).Scan(
		&pol.ID, &pol.ScopeID, &pol.ScopeType, &pol.RegoBundleRef, &pol.Active, &pol.CreatedAt, &pol.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}
	return pol, nil
}

func (s *PolicyStore) GetActiveByScope(ctx context.Context, scopeID, scopeType string) ([]*Policy, error) {
	if s.pool == nil {
		return nil, nil
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, scope_id, scope_type, rego_bundle_ref, active, created_at, updated_at
		FROM policies WHERE scope_id = $1 AND scope_type = $2 AND active = true`, scopeID, scopeType)
	if err != nil {
		return nil, fmt.Errorf("failed to get active policies by scope: %w", err)
	}
	defer rows.Close()

	var policies []*Policy
	for rows.Next() {
		var pol Policy
		if err := rows.Scan(&pol.ID, &pol.ScopeID, &pol.ScopeType, &pol.RegoBundleRef, &pol.Active, &pol.CreatedAt, &pol.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan policy: %w", err)
		}
		policies = append(policies, &pol)
	}
	return policies, nil
}

func (s *PolicyStore) ListPolicies(ctx context.Context, cursor string, limit int) ([]*Policy, string, error) {
	if s.pool == nil {
		return nil, "", nil
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	
	query := `SELECT id, scope_id, scope_type, rego_bundle_ref, active, created_at, updated_at FROM policies`
	var args []interface{}
	if cursor != "" {
		query += ` WHERE id > $1`
		args = append(args, cursor)
	}
	query += ` ORDER BY id ASC LIMIT $2`
	if cursor == "" {
		args = append(args, limit)
		query = `SELECT id, scope_id, scope_type, rego_bundle_ref, active, created_at, updated_at FROM policies ORDER BY id ASC LIMIT $1`
	} else {
		args = append(args, limit)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list policies: %w", err)
	}
	defer rows.Close()

	var policies []*Policy
	var nextCursor string
	for rows.Next() {
		var pol Policy
		if err := rows.Scan(&pol.ID, &pol.ScopeID, &pol.ScopeType, &pol.RegoBundleRef, &pol.Active, &pol.CreatedAt, &pol.UpdatedAt); err != nil {
			return nil, "", fmt.Errorf("failed to scan policy: %w", err)
		}
		policies = append(policies, &pol)
		nextCursor = pol.ID
	}
	return policies, nextCursor, nil
}

func (s *PolicyStore) CreateVersion(ctx context.Context, policyID, regoContent, diff, actor string) (*PolicyVersion, error) {
	if s.pool == nil {
		return &PolicyVersion{ID: uuid.NewString(), PolicyID: policyID}, nil
	}
	var maxVersion int
	err := s.pool.QueryRow(ctx, `SELECT COALESCE(MAX(version), 0) FROM policy_versions WHERE policy_id = $1`, policyID).Scan(&maxVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get max version: %w", err)
	}

	ver := &PolicyVersion{
		ID:          uuid.NewString(),
		PolicyID:    policyID,
		Version:     maxVersion + 1,
		RegoContent: regoContent,
		Diff:        diff,
		Actor:       actor,
		CreatedAt:   time.Now(),
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO policy_versions (id, policy_id, version, rego_content, diff, actor, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		ver.ID, ver.PolicyID, ver.Version, ver.RegoContent, ver.Diff, ver.Actor, ver.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert policy version: %w", err)
	}
	return ver, nil
}

func (s *PolicyStore) GetVersions(ctx context.Context, policyID string) ([]*PolicyVersion, error) {
	if s.pool == nil {
		return nil, nil
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, policy_id, version, rego_content, diff, actor, created_at
		FROM policy_versions WHERE policy_id = $1 ORDER BY version DESC`, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}
	defer rows.Close()

	var versions []*PolicyVersion
	for rows.Next() {
		var ver PolicyVersion
		if err := rows.Scan(&ver.ID, &ver.PolicyID, &ver.Version, &ver.RegoContent, &ver.Diff, &ver.Actor, &ver.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan version: %w", err)
		}
		versions = append(versions, &ver)
	}
	return versions, nil
}
