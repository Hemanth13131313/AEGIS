package export

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// CEFEvent represents a AEGIS detection in CEF format.
// CEF format: CEF:Version|Device Vendor|Device Product|Device Version|Signature ID|Name|Severity|Extension
type CEFEvent struct {
	DeviceVendor  string
	DeviceProduct string
	DeviceVersion string
	SignatureID   string // OWASP LLM category (LLM01 etc.)
	Name          string
	Severity      int // 0-10 (0=Unknown, 1-3=Low, 4-6=Medium, 7-8=High, 9-10=Very-High)
	Extension     map[string]string
}

// FormatCEF returns the CEF 0 formatted string for a detection.
func FormatCEF(e CEFEvent) string {
	var extBuilder strings.Builder
	for k, v := range e.Extension {
		// Escape = and \n
		val := strings.ReplaceAll(v, "\\", "\\\\")
		val = strings.ReplaceAll(val, "=", "\\=")
		val = strings.ReplaceAll(val, "\n", "\\n")
		extBuilder.WriteString(fmt.Sprintf("%s=%s ", k, val))
	}
	
	extStr := strings.TrimSpace(extBuilder.String())
	
	// Escape pipes and backslashes in headers
	vendor := escapeHeader(e.DeviceVendor)
	product := escapeHeader(e.DeviceProduct)
	version := escapeHeader(e.DeviceVersion)
	sigID := escapeHeader(e.SignatureID)
	name := escapeHeader(e.Name)

	return fmt.Sprintf("CEF:0|%s|%s|%s|%s|%s|%d|%s",
		vendor, product, version, sigID, name, e.Severity, extStr)
}

func escapeHeader(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "|", "\\|")
	return s
}

// DetectionToCEF converts an AEGIS DetectionRecord to a CEFEvent.
func DetectionToCEF(sessionID, category, owaspID, atlasID, action, policyRef string, confidence float64, ts time.Time) CEFEvent {
	severity := int(confidence * 10)
	if severity > 10 {
		severity = 10
	} else if severity < 0 {
		severity = 0
	}

	return CEFEvent{
		DeviceVendor:  "AEGIS",
		DeviceProduct: "Security Platform",
		DeviceVersion: "1.0",
		SignatureID:   owaspID,
		Name:          category,
		Severity:      severity,
		Extension: map[string]string{
			"act": action,
			"msg": policyRef,
			"sessionid": sessionID,
			"cat": atlasID,
			"rt": fmt.Sprintf("%d", ts.UnixMilli()),
		},
	}
}

// SyslogWriter writes CEF events to a UDP/TCP syslog target.
type SyslogWriter struct {
	network string
	addr    string
	conn    net.Conn
}

func NewSyslogWriter(network, addr string) (*SyslogWriter, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to syslog target %s:%s: %w", network, addr, err)
	}
	return &SyslogWriter{
		network: network,
		addr:    addr,
		conn:    conn,
	}, nil
}

func (s *SyslogWriter) Write(event CEFEvent) error {
	msg := FormatCEF(event) + "\n"
	_, err := s.conn.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write CEF event: %w", err)
	}
	return nil
}

func (s *SyslogWriter) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
