package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// BedrockAdapter implements Adapter for AWS Bedrock (Converse API).
type BedrockAdapter struct {
	backendURL string
	model      string
}

func NewBedrockAdapter(backendURL, model string) *BedrockAdapter {
	return &BedrockAdapter{
		backendURL: backendURL,
		model:      model,
	}
}

func (a *BedrockAdapter) Name() string {
	return "bedrock"
}

func (a *BedrockAdapter) BackendURL() string {
	return a.backendURL
}

func (a *BedrockAdapter) NormalizeRequest(ctx context.Context, r *http.Request) (*CanonicalRequest, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("empty request body")
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var payload struct {
		System []struct {
			Text string `json:"text"`
		} `json:"system"`
		Messages []struct {
			Role    string `json:"role"`
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"messages"`
		ToolConfig struct {
			Tools []struct {
				ToolSpec struct {
					Name string `json:"name"`
				} `json:"toolSpec"`
			} `json:"tools"`
		} `json:"toolConfig"`
	}

	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse bedrock request: %w", err)
	}

	req := &CanonicalRequest{
		Model: a.model,
	}

	if len(payload.System) > 0 {
		sysText := ""
		for _, s := range payload.System {
			sysText += s.Text
		}
		req.Messages = append(req.Messages, CanonicalMessage{
			Role:    "system",
			Content: sysText,
		})
		req.RoleSequence = append(req.RoleSequence, "system")
	}

	for _, msg := range payload.Messages {
		text := ""
		for _, c := range msg.Content {
			text += c.Text
		}
		req.Messages = append(req.Messages, CanonicalMessage{
			Role:    msg.Role,
			Content: text,
		})
		req.RoleSequence = append(req.RoleSequence, msg.Role)
	}

	for _, tool := range payload.ToolConfig.Tools {
		if tool.ToolSpec.Name != "" {
			req.ToolNames = append(req.ToolNames, tool.ToolSpec.Name)
		}
	}

	return req, nil
}

func (a *BedrockAdapter) NormalizeResponse(ctx context.Context, body []byte) (*CanonicalResponse, error) {
	var payload struct {
		Output struct {
			Message struct {
				Content []struct {
					Text     string `json:"text"`
					ToolUse  struct {
						Name  string                 `json:"name"`
						Input map[string]interface{} `json:"input"`
					} `json:"toolUse"`
				} `json:"content"`
			} `json:"message"`
		} `json:"output"`
		StopReason string `json:"stopReason"`
		Usage      struct {
			TotalTokens int `json:"totalTokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse bedrock response: %w", err)
	}

	resp := &CanonicalResponse{
		TokensUsed:   payload.Usage.TotalTokens,
		FinishReason: payload.StopReason,
	}

	for _, c := range payload.Output.Message.Content {
		if c.Text != "" {
			resp.Content += c.Text
		}
		if c.ToolUse.Name != "" {
			resp.ToolCalls = append(resp.ToolCalls, CanonicalToolCall{
				Name:   c.ToolUse.Name,
				Params: c.ToolUse.Input,
			})
		}
	}

	return resp, nil
}
