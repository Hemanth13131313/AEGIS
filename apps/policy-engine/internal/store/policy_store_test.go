package store

import (
	"context"
	"testing"
)

func TestCreate_Success(t *testing.T) {
	s := NewPolicyStore(nil, nil) // nil pool for mock
	pol, err := s.Create(context.Background(), CreatePolicyInput{
		ScopeID: "test-scope", ScopeType: "organization", RegoContent: "test", Actor: "test-actor",
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if pol.ScopeID != "test-scope" {
		t.Errorf("Expected test-scope, got %v", pol.ScopeID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	s := NewPolicyStore(nil, nil)
	_, err := s.GetByID(context.Background(), "unknown")
	if err == nil {
		t.Errorf("Expected error for uninitialized pool / not found")
	}
}

func TestListPolicies_Pagination(t *testing.T) {
	s := NewPolicyStore(nil, nil)
	pols, cursor, err := s.ListPolicies(context.Background(), "", 10)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(pols) != 0 || cursor != "" {
		t.Errorf("Expected empty response for mock, got %v %v", len(pols), cursor)
	}
}
