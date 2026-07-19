package api

import (
	"encoding/json"
	"net/http"

	"github.com/aegis-security/aegis/apps/policy-engine/internal/hierarchy"
	"github.com/aegis-security/aegis/apps/policy-engine/internal/opa"
	"github.com/aegis-security/aegis/apps/policy-engine/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	store     store.PolicyStoreInterface
	evaluator *opa.Evaluator
	resolver  *hierarchy.Resolver
	logger    *zap.Logger
}

func NewHandler(s store.PolicyStoreInterface, e *opa.Evaluator, r *hierarchy.Resolver, l *zap.Logger) *Handler {
	return &Handler{
		store:     s,
		evaluator: e,
		resolver:  r,
		logger:    l,
	}
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/api/v1/health", h.HealthCheck)
	r.Get("/api/v1/policies", h.ListPolicies)
	r.Post("/api/v1/policies", h.CreatePolicy)
	r.Get("/api/v1/policies/{id}", h.GetPolicy)
	r.Put("/api/v1/policies/{id}", h.UpdatePolicy)
	r.Get("/api/v1/policies/{id}/versions", h.GetPolicyVersions)
	r.Post("/api/v1/check", h.CheckHandler)
}

func writeError(w http.ResponseWriter, code int, reqID, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":       "POLICY_ERROR",
			"category":   "policy",
			"message":    msg,
			"request_id": reqID,
		},
	})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "policy-engine"})
}

func (h *Handler) ListPolicies(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	policies, next, err := h.store.ListPolicies(r.Context(), cursor, 20)
	if err != nil {
		h.logger.Error("failed to list policies", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "", "failed to list policies")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":        policies,
		"next_cursor": next,
	})
}

func (h *Handler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	var input store.CreatePolicyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "", "invalid JSON body")
		return
	}

	validScopes := map[string]bool{"organization": true, "application": true, "model_endpoint": true, "environment": true}
	if !validScopes[input.ScopeType] {
		writeError(w, http.StatusBadRequest, "", "invalid scope_type")
		return
	}

	pol, err := h.store.Create(r.Context(), input)
	if err != nil {
		h.logger.Error("failed to create policy", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "", "failed to create policy")
		return
	}

	if err := h.evaluator.LoadRegoContent(r.Context(), pol.ID, input.RegoContent); err != nil {
		h.logger.Error("failed to load rego content", zap.Error(err))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pol)
}

func (h *Handler) GetPolicy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pol, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "", "policy not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pol)
}

func (h *Handler) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input struct {
		RegoContent string `json:"rego_content"`
		Actor       string `json:"actor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "", "invalid JSON")
		return
	}
	
	_, err := h.store.CreateVersion(r.Context(), id, input.RegoContent, "", input.Actor)
	if err != nil {
		h.logger.Error("failed to update policy", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "", "failed to update policy")
		return
	}

	if err := h.evaluator.LoadRegoContent(r.Context(), id, input.RegoContent); err != nil {
		h.logger.Error("failed to reload rego content", zap.Error(err))
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetPolicyVersions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	vers, err := h.store.GetVersions(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "", "failed to fetch versions")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": vers})
}
