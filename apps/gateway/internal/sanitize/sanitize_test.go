package sanitize

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestSanitize_ValidPayload_Passes(t *testing.T) {
	logger := zap.NewNop()
	s := &Sanitizer{
		config: Config{MaxPayloadBytes: 1024, EnforceUTF8: true},
		logger: logger,
	}

	payload := []byte("hello world")
	result, err := s.Sanitize(context.Background(), payload)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(result) != string(payload) {
		t.Errorf("Expected unchanged payload")
	}
}

func TestSanitize_InvalidUTF8_Rejected(t *testing.T) {
	logger := zap.NewNop()
	s := &Sanitizer{
		config: Config{MaxPayloadBytes: 1024, EnforceUTF8: true},
		logger: logger,
	}

	payload := []byte{0xFF, 0xFE}
	_, err := s.Sanitize(context.Background(), payload)
	if err == nil {
		t.Errorf("Expected error for invalid UTF-8")
	}

	sErr, ok := err.(SanitizeError)
	if !ok || sErr.Code != "INVALID_CHARSET" {
		t.Errorf("Expected INVALID_CHARSET error, got %v", err)
	}
}

func TestSanitize_OversizePayload_Rejected(t *testing.T) {
	logger := zap.NewNop()
	s := &Sanitizer{
		config: Config{MaxPayloadBytes: 5, EnforceUTF8: true},
		logger: logger,
	}

	payload := []byte("too long")
	_, err := s.Sanitize(context.Background(), payload)
	if err == nil {
		t.Errorf("Expected error for oversize payload")
	}

	sErr, ok := err.(SanitizeError)
	if !ok || sErr.Code != "PAYLOAD_TOO_LARGE" {
		t.Errorf("Expected PAYLOAD_TOO_LARGE error, got %v", err)
	}
}

func TestSanitize_CustomConfig(t *testing.T) {
	logger := zap.NewNop()
	s := &Sanitizer{
		config: Config{MaxPayloadBytes: 10, EnforceUTF8: true},
		logger: logger,
	}

	payload := []byte("12345678901") // 11 bytes
	_, err := s.Sanitize(context.Background(), payload)
	if err == nil {
		t.Errorf("Expected error")
	}
}
