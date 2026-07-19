package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Client is a Vault client supporting both dev (token) and production (K8s) auth.
type Client struct {
	addr   string
	token  string
	logger *zap.Logger
	http   *http.Client
}

// NewClientFromEnv creates a Vault client.
func NewClientFromEnv(logger *zap.Logger) (*Client, error) {
	addr := os.Getenv("AEGIS_VAULT_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("AEGIS_VAULT_ADDR not set")
	}

	c := &Client{
		addr:   addr,
		logger: logger,
		http:   &http.Client{Timeout: 10 * time.Second},
	}

	// Dev mode: direct token
	if devToken := os.Getenv("VAULT_DEV_ROOT_TOKEN_ID"); devToken != "" {
		logger.Warn("Using Vault dev token — NOT for production",
			zap.String("note", "Set VAULT_ROLE and remove VAULT_DEV_ROOT_TOKEN_ID in production"))
		c.token = devToken
		return c, nil
	}

	// Production: K8s ServiceAccount auth
	token, err := c.k8sAuth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("vault: k8s auth failed: %w", err)
	}
	c.token = token
	return c, nil
}

func (c *Client) k8sAuth(ctx context.Context) (string, error) {
	role := os.Getenv("AEGIS_VAULT_ROLE")
	if role == "" {
		return "", fmt.Errorf("AEGIS_VAULT_ROLE not set")
	}

	jwtBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return "", fmt.Errorf("read k8s service account token: %w", err)
	}

	body := fmt.Sprintf(`{"jwt":%q,"role":%q}`, strings.TrimSpace(string(jwtBytes)), role)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.addr+"/v1/auth/kubernetes/login",
		strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("vault auth request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode auth response: %w", err)
	}
	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("vault: empty client token in auth response")
	}
	c.logger.Info("Vault K8s auth successful", zap.String("role", role))
	return result.Auth.ClientToken, nil
}

func (c *Client) GetSecret(ctx context.Context, path string) (map[string]string, error) {
	url := fmt.Sprintf("%s/v1/%s", c.addr, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("vault: create request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("vault: access denied to %s (check Vault policy for role %s)",
			path, os.Getenv("AEGIS_VAULT_ROLE"))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault: unexpected status %d at %s", resp.StatusCode, path)
	}

	var result struct {
		Data struct {
			Data map[string]string `json:"data"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("vault: decode response: %w", err)
	}
	c.logger.Info("Vault secret fetched", zap.String("path", path))
	return result.Data.Data, nil
}
