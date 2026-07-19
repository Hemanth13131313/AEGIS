package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aegis-security/aegis/apps/gateway/internal/policyclient"
	"go.uber.org/zap"
)

// PolicyBlockError is returned when a request is blocked by policy.
type PolicyBlockError struct {
	PolicyRef string
	Reason    string
	RequestID string
}

func (e *PolicyBlockError) Error() string {
	return fmt.Sprintf("policy block: %s (ref: %s)", e.Reason, e.PolicyRef)
}

// writePolicyBlock writes a policy block response in AEGIS error shape.
func writePolicyBlock(w http.ResponseWriter, e *PolicyBlockError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	resp := map[string]interface{}{
		"error": map[string]interface{}{
			"code":       "POLICY_BLOCKED",
			"category":   "policy",
			"message":    "Request blocked by policy: " + e.Reason,
			"policy_ref": e.PolicyRef,
			"request_id": e.RequestID,
		},
	}
	_ = json.NewEncoder(w).Encode(resp)
}
