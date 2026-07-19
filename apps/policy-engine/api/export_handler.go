package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aegis-security/aegis/apps/policy-engine/internal/export"
)

type ConfigSyslogRequest struct {
	Protocol string `json:"protocol"` // udp|tcp
	Addr     string `json:"addr"`
}

type ConfigWebhookRequest struct {
	URL    string `json:"url"`
	Secret string `json:"secret"`
}

type ExportHandler struct {
	// In a real app, these would be managed in a registry or config struct
	syslogWriter    *export.SyslogWriter
	webhookExporter *export.WebhookExporter
}

func NewExportHandler() *ExportHandler {
	return &ExportHandler{}
}

func (h *ExportHandler) ConfigureSIEM(w http.ResponseWriter, r *http.Request) {
	var req ConfigSyslogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writer, err := export.NewSyslogWriter(req.Protocol, req.Addr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.syslogWriter = writer
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "target": req.Addr})
}

func (h *ExportHandler) ConfigureWebhook(w http.ResponseWriter, r *http.Request) {
	var req ConfigWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.webhookExporter = export.NewWebhookExporter(req.URL, req.Secret)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *ExportHandler) TestExport(w http.ResponseWriter, r *http.Request) {
	ts := time.Now()
	
	if h.syslogWriter != nil {
		event := export.DetectionToCEF("test-sess", "Test Detection", "LLM01", "AML.T0051", "log", "test-policy", 0.9, ts)
		if err := h.syslogWriter.Write(event); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if h.webhookExporter != nil {
		payload := export.WebhookPayload{
			Source:     "aegis",
			Version:    "1.0",
			EventType:  "detection",
			SessionID:  "test-sess",
			Category:   "Test Detection",
			OWASPId:    "LLM01",
			ATLASId:    "AML.T0051",
			Action:     "log",
			Confidence: 0.9,
			PolicyRef:  "test-policy",
			Timestamp:  ts,
		}
		if err := h.webhookExporter.Export(r.Context(), payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
