package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenAIAdapter implements Adapter for OpenAI-compatible APIs.
type OpenAIAdapter struct {
	backendURL string
	model      string
}

func NewOpenAIAdapter(backendURL, model string) *OpenAIAdapter {
	return &OpenAIAdapter{
		backendURL: backendURL,
		model:      model,
	}
}

func (a *OpenAIAdapter) Name() string {
	return "openai"
}

func (a *OpenAIAdapter) BackendURL() string {
	return a.backendURL
}

func (a *OpenAIAdapter) NormalizeRequest(ctx context.Context, r *http.Request) (*CanonicalRequest, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("empty request body")
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var payload struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		MaxTokens int `json:"max_tokens"`
		Tools     []struct {
			Type     string `json:"type"`
			Function struct {
				Name string `json:"name"`
			} `json:"function"`
		} `json:"tools"`
	}

	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse openai request: %w", err)
	}

	req := &CanonicalRequest{
		Model:     payload.Model,
		MaxTokens: payload.MaxTokens,
	}
	if req.Model == "" {
		req.Model = a.model
	}

	for _, msg := range payload.Messages {
		req.Messages = append(req.Messages, CanonicalMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
		req.RoleSequence = append(req.RoleSequence, msg.Role)
	}

	for _, t := range payload.Tools {
		if t.Type == "function" {
			req.ToolNames = append(req.ToolNames, t.Function.Name)
		}
	}

	return req, nil
}

func (a *OpenAIAdapter) NormalizeResponse(ctx context.Context, body []byte) (*CanonicalResponse, error) {
	var payload struct {
		Choices []struct {
			Message struct {
				Content   string `json:"content"`
				ToolCalls []struct {
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse openai response: %w", err)
	}

	resp := &CanonicalResponse{
		TokensUsed: payload.Usage.TotalTokens,
	}

	if len(payload.Choices) > 0 {
		choice := payload.Choices[0]
		resp.Content = choice.Message.Content
		resp.FinishReason = choice.FinishReason

		for _, tc := range choice.Message.ToolCalls {
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err == nil {
				resp.ToolCalls = append(resp.ToolCalls, CanonicalToolCall{
					Name:   tc.Function.Name,
					Params: args,
				})
			}
		}
	}

	return resp, nil
}
