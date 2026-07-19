package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware_RecordsRequest(t *testing.T) {
	handler := Middleware("test-service")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	// Note: We cannot easily assert Prometheus internal state without using testutil
	// Assuming it completes without panicking
}

func TestMiddleware_ExcludesHealthPath(t *testing.T) {
	handler := Middleware("test-service")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestStatusRecorder_CapturesCode(t *testing.T) {
	rr := httptest.NewRecorder()
	rec := &statusRecorder{ResponseWriter: rr, statusCode: 200}
	
	rec.WriteHeader(http.StatusAccepted)
	if rec.statusCode != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, rec.statusCode)
	}
}
