package grpcserver

import (
	"context"
	"crypto/tls"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type ServerConfig struct {
	Addr           string
	TLSConfig      *tls.Config
	MaxRecvMsgSize int
}

func NewGRPCServer(cfg ServerConfig, logger *zap.Logger) *grpc.Server {
	var opts []grpc.ServerOption

	if cfg.MaxRecvMsgSize == 0 {
		cfg.MaxRecvMsgSize = 4 * 1024 * 1024
	}
	opts = append(opts, grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize))

	if cfg.TLSConfig != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(cfg.TLSConfig)))
		logger.Info("Configuring gRPC server with mTLS")
	} else {
		logger.Warn("Configuring gRPC server in insecure mode (dev without SPIRE)")
	}

	// In a real app, unary interceptors for logging, recovery, OTel go here
	// opts = append(opts, grpc.ChainUnaryInterceptor(...))

	s := grpc.NewServer(opts...)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection for grpcurl
	reflection.Register(s)

	return s
}

func Start(ctx context.Context, s *grpc.Server, addr string, logger *zap.Logger) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		logger.Info("Starting gRPC server", zap.String("addr", addr))
		if err := s.Serve(lis); err != nil {
			logger.Error("gRPC server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("Gracefully stopping gRPC server")
	s.GracefulStop()
	return nil
}
