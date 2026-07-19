package adapters

import (
	"context"
	"net/http"
)

// Adapter normalizes provider-specific request/response formats.
type Adapter interface {
	// Name returns the adapter identifier (openai, anthropic, gemini, bedrock)
	Name() string
	// NormalizeRequest converts a provider request to AEGIS canonical format for inspection
	NormalizeRequest(ctx context.Context, r *http.Request) (*CanonicalRequest, error)
	// NormalizeResponse converts a provider response to AEGIS canonical format
	NormalizeResponse(ctx context.Context, body []byte) (*CanonicalResponse, error)
	// BackendURL returns the actual upstream URL for this provider
	BackendURL() string
}

// CanonicalRequest is the AEGIS-internal normalized request for inspection.
type CanonicalRequest struct {
	Messages     []CanonicalMessage
	Model        string
	MaxTokens    int
	ToolNames    []string
	TokenCount   int
	RoleSequence []string
}

type CanonicalMessage struct {
	Role    string // system, user, assistant, tool
	Content string // plain text (no raw embedding)
}

// CanonicalResponse is the AEGIS-internal normalized response.
type CanonicalResponse struct {
	Content      string
	TokensUsed   int
	FinishReason string
	ToolCalls    []CanonicalToolCall
}

type CanonicalToolCall struct {
	Name   string
	Params map[string]interface{}
}
