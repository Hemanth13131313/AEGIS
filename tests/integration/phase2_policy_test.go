//go:build integration

package integration_test

import (
	"testing"
)

func TestPhase2_PolicyBlock_ToolNotAllowed(t *testing.T) {
	// Start gateway + policy engine, configure tool_allowlist policy,
	// send request with blocked tool, assert 403 with POLICY_BLOCKED error shape.
	t.Skip("Enable in CI integration stage")
}

func TestPhase2_PolicyCacheFailover(t *testing.T) {
	// Kill policy engine, verify gateway uses last-known-good cached policy.
	t.Skip("Enable in CI integration stage")
}

func TestPhase2_PolicyChange_PropagatesWithinSLA(t *testing.T) {
	// Create policy, update it, verify new policy evaluated within 5 seconds.
	t.Skip("Enable in CI integration stage")
}
