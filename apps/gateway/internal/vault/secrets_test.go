package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestGetSecret_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"data":{"secret_key":"secret_value"}}}`))
	}))
	defer mockServer.Close()

	client := &Client{
		addr:   mockServer.URL,
		token:  "test-token",
		logger: zap.NewNop(),
		http:   mockServer.Client(),
	}

	secret, err := client.GetSecret(context.Background(), "secret/data/my-secret")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val := secret["secret_key"]; val != "secret_value" {
		t.Errorf("Expected secret_value, got %s", val)
	}
}

func TestGetSecret_403(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer mockServer.Close()

	client := &Client{
		addr:   mockServer.URL,
		token:  "bad-token",
		logger: zap.NewNop(),
		http:   mockServer.Client(),
	}

	_, err := client.GetSecret(context.Background(), "secret/data/my-secret")
	if err == nil {
		t.Error("Expected error for 403 response")
	}
}
