package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestNewReverseProxy_InvalidUpstreamURL(t *testing.T) {
	logger := zap.NewNop()
	cfg := Config{UpstreamURL: "not a url"}
	_, err := NewReverseProxy(cfg, logger)
	if err == nil {
		t.Error("Expected error for invalid upstream URL")
	}
}

func TestHandler_ProxiesToUpstream(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("upstream response"))
	}))
	defer upstream.Close()

	logger := zap.NewNop()
	cfg := Config{UpstreamURL: upstream.URL}
	proxy, err := NewReverseProxy(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	proxy.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if rec.Header().Get("X-AEGIS-Policy-Version") != "phase1-stub" {
		t.Errorf("Expected X-AEGIS-Policy-Version header")
	}
}

func TestErrorHandler_UpstreamError(t *testing.T) {
	logger := zap.NewNop()
	cfg := Config{UpstreamURL: "http://localhost:1"} // Guaranteed to fail
	proxy, err := NewReverseProxy(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "err-req-123")
	rec := httptest.NewRecorder()

	proxy.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("Expected status 502, got %d", rec.Code)
	}

	var resp APIError
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Error.Category != "upstream" {
		t.Errorf("Expected category upstream, got %s", resp.Error.Category)
	}
}
