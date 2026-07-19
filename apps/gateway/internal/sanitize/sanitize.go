package sanitize

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"

	"go.uber.org/zap"
)

// Config holds configuration for the sanitizer.
type Config struct {
	MaxPayloadBytes  int
	MaxTokenEstimate int
	EnforceUTF8      bool
}

// Sanitizer processes and validates payloads.
type Sanitizer struct {
	config Config
	logger *zap.Logger
}

// SanitizeError represents an error during sanitization.
type SanitizeError struct {
	Code   string
	Detail string
}

func (e SanitizeError) Error() string {
	return fmt.Sprintf("Sanitize error %s: %s", e.Code, e.Detail)
}

// NewSanitizerFromEnv creates a Sanitizer based on environment variables.
func NewSanitizerFromEnv(logger *zap.Logger) *Sanitizer {
	cfg := Config{
		MaxPayloadBytes:  1024 * 1024, // 1MB default
		MaxTokenEstimate: 8192,
		EnforceUTF8:      true,
	}

	if maxBytesStr := os.Getenv("AEGIS_SANITIZE_MAX_PAYLOAD_BYTES"); maxBytesStr != "" {
		if val, err := strconv.Atoi(maxBytesStr); err == nil {
			cfg.MaxPayloadBytes = val
		}
	}

	return &Sanitizer{
		config: cfg,
		logger: logger,
	}
}

// Sanitize checks and optionally modifies the payload.
func (s *Sanitizer) Sanitize(ctx context.Context, payload []byte) ([]byte, error) {
	if len(payload) > s.config.MaxPayloadBytes {
		s.logger.Warn("Payload too large", zap.Int("size", len(payload)))
		return nil, SanitizeError{Code: "PAYLOAD_TOO_LARGE", Detail: "Payload exceeds maximum allowed size"}
	}

	if s.config.EnforceUTF8 && !utf8.Valid(payload) {
		s.logger.Warn("Invalid UTF-8 payload")
		return nil, SanitizeError{Code: "INVALID_CHARSET", Detail: "Payload contains invalid UTF-8 sequences"}
	}

	// Token estimate heuristic
	tokenEstimate := len(payload) / 4
	s.logger.Debug("Token estimate", zap.Int("estimate", tokenEstimate))

	// Return unchanged for Phase 1
	return payload, nil
}
