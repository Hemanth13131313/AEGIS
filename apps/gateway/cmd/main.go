package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aegis-security/aegis/apps/gateway/internal/auth"
	"github.com/aegis-security/aegis/apps/gateway/internal/identity"
	"github.com/aegis-security/aegis/apps/gateway/internal/policyclient"
	"github.com/aegis-security/aegis/apps/gateway/internal/proxy"
	"github.com/aegis-security/aegis/apps/gateway/internal/sanitize"
	"github.com/aegis-security/aegis/apps/gateway/internal/telemetry"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	listenAddr := getEnv("AEGIS_GATEWAY_LISTEN_ADDR", ":8080")
	upstreamURL := getEnv("AEGIS_GATEWAY_UPSTREAM_URL", "http://localhost:5000")
	policyAddr := getEnv("AEGIS_GATEWAY_POLICY_ENGINE_ADDR", "localhost:9090")
	failMode := getEnv("AEGIS_GATEWAY_FAIL_MODE", "closed")
	jwksURL := getEnv("AEGIS_GATEWAY_JWKS_URL", "")
	skipVerifyStr := getEnv("AEGIS_GATEWAY_AUTH_SKIP_VERIFY", "false")
	otelEndpoint := getEnv("AEGIS_OTEL_EXPORTER_ENDPOINT", "")
	spireSocket := getEnv("AEGIS_SPIRE_SOCKET", "unix:///tmp/spire-agent/public/api.sock")

	skipVerify, _ := strconv.ParseBool(skipVerifyStr)

	logger.Info("Starting AEGIS Gateway",
		zap.String("listen_addr", listenAddr),
		zap.String("upstream", upstreamURL),
		zap.String("policy_engine", policyAddr),
		zap.String("fail_mode", failMode),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init OTel
	shutdownOTel, err := telemetry.Init(ctx, "aegis-gateway", "0.1.0", otelEndpoint)
	if err != nil {
		logger.Fatal("Failed to initialize telemetry", zap.Error(err))
	}
	defer func() { _ = shutdownOTel(context.Background()) }()

	// Fetch SPIFFE Identity
	identCtx, identCancel := context.WithTimeout(ctx, 5*time.Second)
	ident, err := identity.Fetch(identCtx, spireSocket, logger)
	identCancel()
	if err != nil {
		logger.Warn("Failed to fetch SPIFFE identity, proceeding without SPIRE", zap.Error(err))
	}
	if ident != nil {
		defer ident.Close()
	}

	// Auth JWKS Cache
	var jwksCache *auth.JWKSCache
	if jwksURL != "" {
		jwksCache = auth.NewJWKSCache(jwksURL, logger)
		go jwksCache.AutoRefresh(ctx)
		fetchCtx, fetchCancel := context.WithTimeout(ctx, 5*time.Second)
		_ = jwksCache.Fetch(fetchCtx)
		fetchCancel()
	} else if !skipVerify {
		logger.Warn("AEGIS_GATEWAY_JWKS_URL is empty but auth verification is not skipped")
	}

	// Sanitizer
	_ = sanitize.NewSanitizerFromEnv(logger)

	// Policy Client with SPIFFE Identity
	pClient, err := policyclient.NewPolicyClient(policyAddr, failMode, logger, ident)
	if err != nil {
		logger.Fatal("Failed to create policy client", zap.Error(err))
	}
	defer pClient.Close()

	// Reverse Proxy
	proxyCfg := proxy.Config{
		UpstreamURL:      upstreamURL,
		FailMode:         failMode,
		PolicyEngineAddr: policyAddr,
	}
	rp, err := proxy.NewReverseProxy(proxyCfg, logger)
	if err != nil {
		logger.Fatal("Failed to create reverse proxy", zap.Error(err))
	}

	// Router
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"aegis-gateway"}`))
	})

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// Replace with actual metrics handler
		w.WriteHeader(http.StatusOK)
	})

	r.Group(func(r chi.Router) {
		if !skipVerify || jwksCache != nil {
			r.Use(auth.Middleware(logger, jwksCache, skipVerify))
		}
		// Metrics middleware would go here
		r.Handle("/*", rp.Handler())
	})

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}

	go func() {
		logger.Info("Listening for requests", zap.String("addr", listenAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}
	logger.Info("Shutdown complete")
}
