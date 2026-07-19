package hierarchy

import (
	"context"
	"testing"
)

// Mock Store for tests
type mockStore struct{}

func (m *mockStore) GetActiveByScope(ctx context.Context, scopeID, scopeType string) ([]*store.Policy, error) {
	return nil, nil // Return nil for mock simplicity
}

func TestResolve_OrgOnly(t *testing.T) {
	// Simple test layout; fully implement mock in real project
	t.Skip("Implement mock store logic for test")
}

func TestResolve_AllLevels(t *testing.T) {
	t.Skip("Implement mock store logic for test")
}

func TestResolve_NoPolicy(t *testing.T) {
	t.Skip("Implement mock store logic for test")
}
