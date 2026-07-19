package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GeminiAdapter implements Adapter for Google Gemini API.
type GeminiAdapter struct {
	backendURL string
	model      string
}

func NewGeminiAdapter(backendURL, model string) *GeminiAdapter {
	return &GeminiAdapter{
		backendURL: backendURL,
		model:      model,
	}
}

func (a *GeminiAdapter) Name() string {
	return "gemini"
}

func (a *GeminiAdapter) BackendURL() string {
	return a.backendURL
}

func (a *GeminiAdapter) NormalizeRequest(ctx context.Context, r *http.Request) (*CanonicalRequest, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("empty request body")
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var payload struct {
		Contents []struct {
			Role  string `json:"role"`
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"contents"`
		SystemInstruction struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"system_instruction"`
		Tools []struct {
			FunctionDeclarations []struct {
				Name string `json:"name"`
			} `json:"functionDeclarations"`
		} `json:"tools"`
	}

	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse gemini request: %w", err)
	}

	req := &CanonicalRequest{
		Model: a.model,
	}

	if len(payload.SystemInstruction.Parts) > 0 {
		sysText := ""
		for _, p := range payload.SystemInstruction.Parts {
			sysText += p.Text
		}
		req.Messages = append(req.Messages, CanonicalMessage{
			Role:    "system",
			Content: sysText,
		})
		req.RoleSequence = append(req.RoleSequence, "system")
	}

	for _, content := range payload.Contents {
		text := ""
		for _, part := range content.Parts {
			text += part.Text
		}
		req.Messages = append(req.Messages, CanonicalMessage{
			Role:    content.Role,
			Content: text,
		})
		req.RoleSequence = append(req.RoleSequence, content.Role)
	}

	for _, tool := range payload.Tools {
		for _, f := range tool.FunctionDeclarations {
			req.ToolNames = append(req.ToolNames, f.Name)
		}
	}

	return req, nil
}

func (a *GeminiAdapter) NormalizeResponse(ctx context.Context, body []byte) (*CanonicalResponse, error) {
	var payload struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text         string                 `json:"text"`
					FunctionCall map[string]interface{} `json:"functionCall"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
		UsageMetadata struct {
			TotalTokenCount int `json:"totalTokenCount"`
		} `json:"usageMetadata"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w", err)
	}

	resp := &CanonicalResponse{
		TokensUsed: payload.UsageMetadata.TotalTokenCount,
	}

	if len(payload.Candidates) > 0 {
		cand := payload.Candidates[0]
		resp.FinishReason = cand.FinishReason
		for _, p := range cand.Content.Parts {
			if p.Text != "" {
				resp.Content += p.Text
			}
			if p.FunctionCall != nil {
				if name, ok := p.FunctionCall["name"].(string); ok {
					var params map[string]interface{}
					if args, ok := p.FunctionCall["args"].(map[string]interface{}); ok {
						params = args
					}
					resp.ToolCalls = append(resp.ToolCalls, CanonicalToolCall{
						Name:   name,
						Params: params,
					})
				}
			}
		}
	}

	return resp, nil
}
