package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// Config holds configuration for the reverse proxy.
type Config struct {
	UpstreamURL      string
	FailMode         string
	PolicyEngineAddr string
}

// Sanitizer defines the interface for payload sanitization.
type Sanitizer interface {
	Sanitize(ctx context.Context, payload []byte) ([]byte, error)
}

// ReverseProxy handles proxying requests to the upstream service.
type ReverseProxy struct {
	config    Config
	logger    *zap.Logger
	proxy     *httputil.ReverseProxy
}

// APIError represents the standard AEGIS error shape.
type APIError struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains details of the error.
type ErrorDetail struct {
	Code      string `json:"code"`
	Category  string `json:"category"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

// NewReverseProxy creates a new ReverseProxy.
func NewReverseProxy(cfg Config, logger *zap.Logger) (*ReverseProxy, error) {
	parsedURL, err := url.Parse(cfg.UpstreamURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid upstream URL: %w", err)
	}

	rp := &ReverseProxy{
		config: cfg,
		logger: logger,
	}

	proxy := &httputil.ReverseProxy{
		Director:       rp.Director(parsedURL),
		ModifyResponse: rp.ModifyResponse(),
		ErrorHandler:   rp.ErrorHandler(),
	}

	rp.proxy = proxy
	return rp, nil
}

// Director modifies the request before it is sent to the upstream.
func (rp *ReverseProxy) Director(target *url.URL) func(*http.Request) {
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path + req.URL.Path // Assume simple path append for now
		req.Host = target.Host

		// Remove hop-by-hop headers
		req.Header.Del("Connection")
		req.Header.Del("Keep-Alive")
		req.Header.Del("Proxy-Authenticate")
		req.Header.Del("Proxy-Authorization")
		req.Header.Del("Te")
		req.Header.Del("Trailer")
		req.Header.Del("Transfer-Encoding")
		req.Header.Del("Upgrade")

		reqID := req.Header.Get("X-Request-Id")
		req.Header.Set("X-AEGIS-Request-ID", reqID)

		if clientIP := req.RemoteAddr; clientIP != "" {
			req.Header.Set("X-Forwarded-For", clientIP)
		}
	}
}

// ModifyResponse modifies the response from the upstream.
func (rp *ReverseProxy) ModifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-AEGIS-Policy-Version", "phase1-stub")
		rp.logger.Info("Upstream response", zap.Int("status", resp.StatusCode))
		return nil
	}
}

// ErrorHandler handles errors during proxying.
func (rp *ReverseProxy) ErrorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		reqID := r.Header.Get("X-Request-Id")
		status := http.StatusBadGateway

		if err == context.Canceled {
			rp.logger.Debug("Client disconnected", zap.String("request_id", reqID))
			status = 499 // Client Closed Request
		} else {
			rp.logger.Error("Upstream error", zap.Error(err), zap.String("request_id", reqID))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		errResp := APIError{
			Error: ErrorDetail{
				Code:      "UPSTREAM_ERROR",
				Category:  "upstream",
				Message:   "Failed to communicate with upstream service",
				RequestID: reqID,
			},
		}
		_ = json.NewEncoder(w).Encode(errResp)
	}
}

// Handler returns the http.Handler for the proxy.
func (rp *ReverseProxy) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		ctx := r.Context()

		tracer := otel.Tracer("aegis.gateway")
		ctx, span := tracer.Start(ctx, "proxy.request")
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("request_id", reqID),
		)

		rp.logger.Info("Proxying request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("upstream", rp.config.UpstreamURL),
			zap.String("request_id", reqID),
		)

		start := time.Now()
		rp.proxy.ServeHTTP(w, r.WithContext(ctx))
		duration := time.Since(start)

		rp.logger.Info("Request completed", zap.Duration("duration", duration), zap.String("request_id", reqID))
	})
}
