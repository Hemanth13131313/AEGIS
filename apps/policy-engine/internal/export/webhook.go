package export

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookPayload is the JSON body sent to webhook/SOAR targets.
type WebhookPayload struct {
	Source     string            `json:"source"`      // "aegis"
	Version    string            `json:"version"`     // "1.0"
	EventType  string            `json:"event_type"`  // "detection"
	SessionID  string            `json:"session_id"`
	Category   string            `json:"category"`
	OWASPId    string            `json:"owasp_id"`
	ATLASId    string            `json:"atlas_id"`
	Action     string            `json:"action"`
	Confidence float64           `json:"confidence"`
	PolicyRef  string            `json:"policy_ref"`
	Timestamp  time.Time         `json:"timestamp"`
	Metadata   map[string]string `json:"metadata"`
}

// WebhookExporter sends detections to a configured webhook endpoint.
type WebhookExporter struct {
	targetURL  string
	secret     string // HMAC-SHA256 signing secret
	httpClient *http.Client
}

func NewWebhookExporter(targetURL, secret string) *WebhookExporter {
	return &WebhookExporter{
		targetURL:  targetURL,
		secret:     secret,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *WebhookExporter) Export(ctx context.Context, payload WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	sig := w.sign(body)

	req, err := http.NewRequestWithContext(ctx, "POST", w.targetURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AEGIS-Signature", sig)

	// Retry once on 5xx
	var resp *http.Response
	for i := 0; i < 2; i++ {
		resp, err = w.httpClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to send webhook after retries: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status: %d", resp.StatusCode)
	}

	return nil
}

func (w *WebhookExporter) sign(data []byte) string {
	h := hmac.New(sha256.New, []byte(w.secret))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func (w *WebhookExporter) Verify(payload []byte, signature string) bool {
	expected := w.sign(payload)
	return hmac.Equal([]byte(expected), []byte(signature))
}
