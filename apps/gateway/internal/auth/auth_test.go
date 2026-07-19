package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestMiddleware_MissingAuthHeader(t *testing.T) {
	logger := zap.NewNop()
	cache := NewJWKSCache("http://localhost/jwks", logger)
	handler := Middleware(logger, cache, false)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "test-req-id")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}

	var resp APIError
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Error.RequestID != "test-req-id" {
		t.Errorf("Expected RequestID test-req-id, got %s", resp.Error.RequestID)
	}
}

func TestMiddleware_InvalidBearerFormat(t *testing.T) {
	logger := zap.NewNop()
	cache := NewJWKSCache("http://localhost/jwks", logger)
	handler := Middleware(logger, cache, false)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic some_token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}

func TestMiddleware_SkipVerify_PassThrough(t *testing.T) {
	logger := zap.NewNop()
	cache := NewJWKSCache("http://localhost/jwks", logger)
	passed := false
	handler := Middleware(logger, cache, true)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passed = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !passed {
		t.Error("Expected request to pass through when skipVerify is true")
	}
}

func TestWriteError_Shape(t *testing.T) {
	rec := httptest.NewRecorder()
	writeError(rec, http.StatusBadRequest, "BAD_REQ", "client", "Bad request", "req-123")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}

	var resp APIError
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Error.Code != "BAD_REQ" {
		t.Errorf("Expected Code BAD_REQ, got %s", resp.Error.Code)
	}
	if resp.Error.Category != "client" {
		t.Errorf("Expected Category client, got %s", resp.Error.Category)
	}
	if resp.Error.Message != "Bad request" {
		t.Errorf("Expected Message Bad request, got %s", resp.Error.Message)
	}
	if resp.Error.RequestID != "req-123" {
		t.Errorf("Expected RequestID req-123, got %s", resp.Error.RequestID)
	}
}
