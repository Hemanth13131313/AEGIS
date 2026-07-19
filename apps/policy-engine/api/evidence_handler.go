package api

import (
	"encoding/json"
	"net/http"
	"context"
	"go.uber.org/zap"
	"fmt"
)

// Placeholder for bundler
type Bundler interface {
	Generate(ctx context.Context) (interface{}, string, error)
}

type EvidenceHandler struct {
	logger *zap.Logger
	bundler Bundler
}

func NewEvidenceHandler(logger *zap.Logger, bundler Bundler) *EvidenceHandler {
	return &EvidenceHandler{
		logger: logger,
		bundler: bundler,
	}
}

func (h *EvidenceHandler) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Actor       string `json:"actor"`
		Environment string `json:"environment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	manifest, bundlePath, err := h.bundler.Generate(r.Context())
	if err != nil {
		h.logger.Error("failed to generate bundle", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"manifest": manifest,
		"bundle_path": bundlePath,
		"download_url": fmt.Sprintf("/api/v1/evidence/download/%s", "latest"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *EvidenceHandler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// ID logic
	w.Header().Set("Content-Disposition", "attachment; filename=\"aegis-evidence-latest.zip\"")
	w.Header().Set("Content-Type", "application/zip")
	// write zip content
	w.Write([]byte{})
}

func (h *EvidenceHandler) LatestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "latest not implemented"})
}
