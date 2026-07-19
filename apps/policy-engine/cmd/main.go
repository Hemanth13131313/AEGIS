package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aegis-security/aegis/apps/policy-engine/api"
	"github.com/aegis-security/aegis/apps/policy-engine/internal/cache"
	grpcserver "github.com/aegis-security/aegis/apps/policy-engine/internal/grpc"
	"github.com/aegis-security/aegis/apps/policy-engine/internal/hierarchy"
	"github.com/aegis-security/aegis/apps/policy-engine/internal/opa"
	"github.com/aegis-security/aegis/apps/policy-engine/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	listenAddr := os.Getenv("AEGIS_POLICY_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8081"
	}
	grpcAddr := os.Getenv("AEGIS_POLICY_GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = ":9090"
	}
	dbDsn := os.Getenv("AEGIS_POLICY_DB_DSN")
	redisAddr := os.Getenv("AEGIS_POLICY_REDIS_ADDR")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Postgres
	var pool *pgxpool.Pool
	if dbDsn != "" {
		p, err := pgxpool.New(ctx, dbDsn)
		if err != nil {
			logger.Fatal("failed to connect to database", zap.Error(err))
		}
		pool = p
		defer pool.Close()
	}

	// Initialize Redis Cache
	_, err := cache.NewPolicyCache(redisAddr, logger)
	if err != nil && redisAddr != "" {
		logger.Warn("failed to connect to redis cache", zap.Error(err))
	}

	// Initialize components
	evaluator := opa.NewEvaluator(logger)
	policyStore := store.NewPolicyStore(pool, logger)
	resolver := hierarchy.NewResolver(policyStore, logger)
	apiHandler := api.NewHandler(policyStore, evaluator, resolver, logger)

	// REST Server
	r := chi.NewRouter()
	apiHandler.Routes(r)

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}

	go func() {
		logger.Info("Starting REST server", zap.String("addr", listenAddr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("REST server error", zap.Error(err))
		}
	}()

	// gRPC Server (Stubbed for Phase 2)
	cfg := grpcserver.ServerConfig{Addr: grpcAddr}
	grpcSrv := grpcserver.NewGRPCServer(cfg, logger)

	// Shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("REST server shutdown error", zap.Error(err))
	}
	grpcSrv.GracefulStop()
	logger.Info("Shutdown complete")
}
