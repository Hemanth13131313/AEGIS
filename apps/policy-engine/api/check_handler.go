package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aegis-security/aegis/apps/policy-engine/internal/opa"
	"go.uber.org/zap"
)

type CheckRequest struct {
	SessionID   string            `json:"session_id"`
	OrgID       string            `json:"org_id"`
	AppID       string            `json:"app_id"`
	ModelID     string            `json:"model_id"`
	Environment string            `json:"environment"`
	Auth        opa.AuthInput     `json:"auth"`
	Payload     opa.PayloadInput  `json:"payload"`
	AllowedTools []string         `json:"allowed_tools"`
	Metadata    map[string]string `json:"metadata"`
}

type CheckResponse struct {
	Action     string  `json:"action"`
	PolicyRef  string  `json:"policy_ref"`
	Reason     string  `json:"reason"`
	Confidence float64 `json:"confidence"`
}

func (h *Handler) CheckHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "", "invalid JSON")
		return
	}

	resolved, err := h.resolver.Resolve(r.Context(), req.OrgID, req.AppID, req.ModelID, req.Environment)
	if err != nil {
		h.logger.Error("Failed to resolve hierarchy", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "", "resolution error")
		return
	}

	input := opa.EvalInput{
		OrgID:        req.OrgID,
		AppID:        req.AppID,
		ModelID:      req.ModelID,
		Environment:  req.Environment,
		Auth:         req.Auth,
		Payload:      req.Payload,
		AllowedTools: req.AllowedTools,
		Metadata:     req.Metadata,
	}

	var allDenyReasons []string
	var finalPolicyRef string

	for _, pol := range resolved {
		res, err := h.evaluator.Evaluate(r.Context(), pol.PolicyID, input)
		if err != nil {
			h.logger.Error("Evaluation error", zap.Error(err))
			continue
		}
		if !res.Allow {
			allDenyReasons = append(allDenyReasons, res.Deny...)
			if finalPolicyRef == "" {
				finalPolicyRef = res.PolicyRef
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if len(allDenyReasons) > 0 {
		h.logger.Info("Policy block", zap.String("org_id", req.OrgID), zap.String("app_id", req.AppID))
		json.NewEncoder(w).Encode(CheckResponse{
			Action:     "block",
			PolicyRef:  finalPolicyRef,
			Reason:     strings.Join(allDenyReasons, ", "),
			Confidence: 1.0,
		})
		return
	}

	h.logger.Info("Policy allow", zap.String("org_id", req.OrgID), zap.String("app_id", req.AppID))
	json.NewEncoder(w).Encode(CheckResponse{
		Action:     "allow",
		PolicyRef:  "resolved-hierarchy",
		Reason:     "",
		Confidence: 1.0,
	})
}
