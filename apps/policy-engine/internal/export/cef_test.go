package export

import (
	"strings"
	"testing"
	"time"
)

func TestFormatCEF_ValidFormat(t *testing.T) {
	e := CEFEvent{
		DeviceVendor:  "AEGIS",
		DeviceProduct: "Security Platform",
		DeviceVersion: "1.0",
		SignatureID:   "LLM01",
		Name:          "Prompt Injection",
		Severity:      9,
		Extension: map[string]string{
			"act": "block",
		},
	}
	result := FormatCEF(e)
	if !strings.HasPrefix(result, "CEF:0|AEGIS|Security Platform|1.0|LLM01|Prompt Injection|9|") {
		t.Errorf("Unexpected CEF format: %s", result)
	}
	if !strings.Contains(result, "act=block") {
		t.Errorf("Missing extension: %s", result)
	}
}

func TestFormatCEF_EscapesSpecialChars(t *testing.T) {
	e := CEFEvent{
		DeviceVendor:  "AEGIS|Vendor",
		DeviceProduct: "Platform\\Test",
		DeviceVersion: "1.0",
		SignatureID:   "LLM01",
		Name:          "Test",
		Severity:      1,
		Extension: map[string]string{
			"msg": "test=test\ntest\\test",
		},
	}
	result := FormatCEF(e)
	if !strings.Contains(result, "AEGIS\\|Vendor") {
		t.Errorf("Failed to escape pipe: %s", result)
	}
	if !strings.Contains(result, "Platform\\\\Test") {
		t.Errorf("Failed to escape backslash: %s", result)
	}
	if !strings.Contains(result, "msg=test\\=test\\ntest\\\\test") {
		t.Errorf("Failed to escape extension: %s", result)
	}
}

func TestDetectionToCEF_SeverityMapping(t *testing.T) {
	ts := time.Now()
	e := DetectionToCEF("sess1", "Injection", "LLM01", "AML.T0051", "block", "pol1", 0.9, ts)
	if e.Severity != 9 {
		t.Errorf("Expected severity 9, got %d", e.Severity)
	}
}
