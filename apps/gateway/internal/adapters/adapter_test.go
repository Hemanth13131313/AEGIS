package adapters

import (
	"bytes"
	"context"
	"net/http"
	"testing"
)

func TestOpenAIAdapter_NormalizeRequest(t *testing.T) {
	body := []byte(`{"model": "gpt-4", "messages": [{"role": "user", "content": "hello"}], "tools": [{"type": "function", "function": {"name": "search"}}]}`)
	req, _ := http.NewRequestWithContext(context.Background(), "POST", "/", bytes.NewReader(body))
	
	adapter := NewOpenAIAdapter("http://localhost", "gpt-4")
	canonical, err := adapter.NormalizeRequest(context.Background(), req)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(canonical.Messages) != 1 || canonical.Messages[0].Content != "hello" {
		t.Errorf("expected 1 message 'hello', got %+v", canonical.Messages)
	}
	if len(canonical.ToolNames) != 1 || canonical.ToolNames[0] != "search" {
		t.Errorf("expected tool 'search', got %v", canonical.ToolNames)
	}
}

func TestAnthropicAdapter_NormalizeRequest(t *testing.T) {
	body := []byte(`{"model": "claude-3", "system": "You are a helpful assistant.", "messages": [{"role": "user", "content": "hi"}]}`)
	req, _ := http.NewRequestWithContext(context.Background(), "POST", "/", bytes.NewReader(body))
	
	adapter := NewAnthropicAdapter("http://localhost", "claude-3")
	canonical, err := adapter.NormalizeRequest(context.Background(), req)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(canonical.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(canonical.Messages))
	}
	if canonical.Messages[0].Role != "system" || canonical.Messages[0].Content != "You are a helpful assistant." {
		t.Errorf("expected system message, got %+v", canonical.Messages[0])
	}
}

func TestGeminiAdapter_NormalizeRequest(t *testing.T) {
	body := []byte(`{"contents": [{"role": "user", "parts": [{"text": "hi"}]}]}`)
	req, _ := http.NewRequestWithContext(context.Background(), "POST", "/", bytes.NewReader(body))
	
	adapter := NewGeminiAdapter("http://localhost", "gemini-1.5")
	canonical, err := adapter.NormalizeRequest(context.Background(), req)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(canonical.Messages) != 1 || canonical.Messages[0].Role != "user" || canonical.Messages[0].Content != "hi" {
		t.Errorf("expected user message 'hi', got %+v", canonical.Messages)
	}
}

func TestRegistry_Get_UnknownProvider(t *testing.T) {
	_, err := Get("unknown", "http://localhost", "model")
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}
