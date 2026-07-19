package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AnthropicAdapter implements Adapter for Anthropic Messages API.
type AnthropicAdapter struct {
	backendURL string
	model      string
}

func NewAnthropicAdapter(backendURL, model string) *AnthropicAdapter {
	return &AnthropicAdapter{
		backendURL: backendURL,
		model:      model,
	}
}

func (a *AnthropicAdapter) Name() string {
	return "anthropic"
}

func (a *AnthropicAdapter) BackendURL() string {
	return a.backendURL
}

func (a *AnthropicAdapter) NormalizeRequest(ctx context.Context, r *http.Request) (*CanonicalRequest, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("empty request body")
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var payload struct {
		Model     string `json:"model"`
		System    string `json:"system"`
		MaxTokens int    `json:"max_tokens"`
		Messages  []struct {
			Role    string      `json:"role"`
			Content interface{} `json:"content"` // can be string or array of objects
		} `json:"messages"`
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}

	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse anthropic request: %w", err)
	}

	req := &CanonicalRequest{
		Model:     payload.Model,
		MaxTokens: payload.MaxTokens,
	}

	if payload.System != "" {
		req.Messages = append(req.Messages, CanonicalMessage{
			Role:    "system",
			Content: payload.System,
		})
		req.RoleSequence = append(req.RoleSequence, "system")
	}

	for _, msg := range payload.Messages {
		contentStr := ""
		switch v := msg.Content.(type) {
		case string:
			contentStr = v
		case []interface{}:
			// simple string extraction for text type
			for _, part := range v {
				if m, ok := part.(map[string]interface{}); ok {
					if t, ok := m["text"].(string); ok {
						contentStr += t
					}
				}
			}
		}

		req.Messages = append(req.Messages, CanonicalMessage{
			Role:    msg.Role,
			Content: contentStr,
		})
		req.RoleSequence = append(req.RoleSequence, msg.Role)
	}

	for _, t := range payload.Tools {
		req.ToolNames = append(req.ToolNames, t.Name)
	}

	return req, nil
}

func (a *AnthropicAdapter) NormalizeResponse(ctx context.Context, body []byte) (*CanonicalResponse, error) {
	var payload struct {
		Content []struct {
			Type  string                 `json:"type"`
			Text  string                 `json:"text"`
			Name  string                 `json:"name"`
			Input map[string]interface{} `json:"input"`
		} `json:"content"`
		StopReason string `json:"stop_reason"`
		Usage      struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse anthropic response: %w", err)
	}

	resp := &CanonicalResponse{
		TokensUsed:   payload.Usage.InputTokens + payload.Usage.OutputTokens,
		FinishReason: payload.StopReason,
	}

	for _, c := range payload.Content {
		if c.Type == "text" {
			resp.Content += c.Text
		} else if c.Type == "tool_use" {
			resp.ToolCalls = append(resp.ToolCalls, CanonicalToolCall{
				Name:   c.Name,
				Params: c.Input,
			})
		}
	}

	return resp, nil
}
