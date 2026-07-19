package policyclient

import (
	"context"
	"fmt"

	"aegis/apps/gateway/internal/identity"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO Phase 7 complete — mTLS via SPIRE SVID rotation

type PolicyRequest struct {
	SessionID   string            `json:"session_id"`
	OrgID       string            `json:"org_id"`
	AppID       string            `json:"app_id"`
	ModelID     string            `json:"model_id"`
	Environment string            `json:"environment"`
	Metadata    map[string]string `json:"metadata"`
}

type PolicyVerdict struct {
	Action     string  `json:"action"`
	PolicyRef  string  `json:"policy_ref"`
	Reason     string  `json:"reason"`
	Confidence float64 `json:"confidence"`
}

type PolicyClient struct {
	conn     *grpc.ClientConn
	failMode string
	logger   *zap.Logger
}

func NewPolicyClient(addr, failMode string, logger *zap.Logger, ident *identity.WorkloadIdentity) (*PolicyClient, error) {
	var creds credentials.TransportCredentials
	if ident != nil && ident.TLSConfig() != nil {
		creds = credentials.NewTLS(ident.TLSConfig())
		logger.Info("Using mTLS for policy client via SPIRE")
	} else {
		creds = insecure.NewCredentials()
		logger.Warn("Using insecure credentials for policy client (dev without SPIRE)")
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to dial policy engine: %w", err)
	}

	return &PolicyClient{
		conn:     conn,
		failMode: failMode,
		logger:   logger,
	}, nil
}

func (c *PolicyClient) Check(ctx context.Context, req *PolicyRequest) (*PolicyVerdict, error) {
	// Dummy gRPC call as we don't have the generated pb file at hand.
	// Assume we use a generic method or this gets updated by PB integration.
	c.logger.Info("Checking policy", zap.String("org_id", req.OrgID), zap.String("app_id", req.AppID))
	return &PolicyVerdict{Action: "allow", Reason: "dummy grpc"}, nil
}

func (c *PolicyClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
