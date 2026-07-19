package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

// APIError represents the standard AEGIS error shape.
type APIError struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains details of the error.
type ErrorDetail struct {
	Code      string `json:"code"`
	Category  string `json:"category"` // "policy", "upstream", "internal"
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

// JWKSCache holds a cached JSON Web Key Set.
type JWKSCache struct {
	mu        sync.RWMutex
	keys      jwk.Set
	jwksURL   string
	lastFetch time.Time
	ttl       time.Duration
	logger    *zap.Logger
}

// NewJWKSCache creates a new JWKSCache.
func NewJWKSCache(jwksURL string, logger *zap.Logger) *JWKSCache {
	return &JWKSCache{
		jwksURL: jwksURL,
		ttl:     5 * time.Minute,
		logger:  logger,
	}
}

// Fetch fetches the JWKS from the URL and updates the cache.
func (c *JWKSCache) Fetch(ctx context.Context) error {
	set, err := jwk.Fetch(ctx, c.jwksURL)
	if err != nil {
		c.logger.Error("Failed to fetch JWKS", zap.Error(err), zap.String("url", c.jwksURL))
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.keys = set
	c.lastFetch = time.Now()
	c.logger.Info("Successfully fetched JWKS", zap.String("url", c.jwksURL))
	return nil
}

// Get returns the cached JWK set.
func (c *JWKSCache) Get() jwk.Set {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.keys
}

// AutoRefresh starts a goroutine to refresh the JWKS periodically.
func (c *JWKSCache) AutoRefresh(ctx context.Context) {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Stopping JWKS AutoRefresh")
			return
		case <-ticker.C:
			_ = c.Fetch(ctx)
		}
	}
}

type contextKey struct {
	name string
}

var claimsKey = &contextKey{"claims"}

// WithClaims injects claims into the context.
func WithClaims(ctx context.Context, claims jwt.Token) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// ClaimsFrom extracts claims from the context.
func ClaimsFrom(ctx context.Context) (jwt.Token, bool) {
	claims, ok := ctx.Value(claimsKey).(jwt.Token)
	return claims, ok
}

// extractRequestID retrieves the request ID from the header or generates a new one.
func extractRequestID(r *http.Request) string {
	id := r.Header.Get("X-Request-Id")
	if id == "" {
		id = uuid.New().String()
	}
	return id
}

// writeError writes a standardized API error response.
func writeError(w http.ResponseWriter, status int, code, category, message, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errResp := APIError{
		Error: ErrorDetail{
			Code:      code,
			Category:  category,
			Message:   message,
			RequestID: requestID,
		},
	}
	_ = json.NewEncoder(w).Encode(errResp)
}

// Middleware creates a new authentication middleware.
func Middleware(logger *zap.Logger, cache *JWKSCache, skipVerify bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := extractRequestID(r)

			if skipVerify {
				logger.Warn("Auth verification skipped (AEGIS_GATEWAY_AUTH_SKIP_VERIFY is true)", zap.String("request_id", requestID))
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Missing Authorization header", zap.String("request_id", requestID))
				writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "policy", "Missing Authorization header", requestID)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				logger.Warn("Invalid Authorization header format", zap.String("request_id", requestID))
				writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "policy", "Invalid Authorization header format", requestID)
				return
			}

			tokenString := parts[1]
			keys := cache.Get()

			if keys == nil {
				// Attempt to fetch if missing
				ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
				defer cancel()
				if err := cache.Fetch(ctx); err != nil {
					logger.Error("Failed to fetch JWKS during request", zap.Error(err), zap.String("request_id", requestID))
					writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "internal", "Failed to validate token", requestID)
					return
				}
				keys = cache.Get()
			}

			token, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(keys))
			if err != nil {
				logger.Warn("Invalid token", zap.Error(err), zap.String("request_id", requestID))
				writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "policy", "Invalid or expired token", requestID)
				return
			}

			ctx := WithClaims(r.Context(), token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
