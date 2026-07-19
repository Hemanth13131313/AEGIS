package identity

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"go.uber.org/zap"
)

type WorkloadIdentity struct {
	source      *workloadapi.X509Source
	trustDomain spiffeid.TrustDomain
	logger      *zap.Logger
	socketPath  string
}

func Fetch(ctx context.Context, socketPath string, logger *zap.Logger) (*WorkloadIdentity, error) {
	if socketPath == "" {
		socketPath = "unix:///tmp/spire-agent/public/api.sock"
	}

	logger.Info("Fetching SPIFFE identity", zap.String("socket", socketPath))

	td, err := spiffeid.TrustDomainFromString("aegis.cluster")
	if err != nil {
		return nil, fmt.Errorf("failed to parse trust domain: %w", err)
	}

	source, err := workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)))
	if err != nil {
		logger.Warn("SPIRE not running or unavailable. Falling back to no-op mode (dev only).", zap.Error(err))
		return nil, nil // Return nil for dev fallback
	}

	svid, err := source.GetX509SVID()
	if err == nil {
		logger.Info("Successfully fetched SPIFFE SVID", zap.String("spiffe_id", svid.ID.String()))
	} else {
		logger.Warn("Could not get initial SVID", zap.Error(err))
	}

	return &WorkloadIdentity{
		source:      source,
		trustDomain: td,
		logger:      logger,
		socketPath:  socketPath,
	}, nil
}

func (w *WorkloadIdentity) TLSConfig() *tls.Config {
	if w == nil || w.source == nil {
		return nil
	}
	return tlsconfig.MTLSClientConfig(w.source, w.source, tlsconfig.AuthorizeAny())
}

func (w *WorkloadIdentity) ServerTLSConfig() *tls.Config {
	if w == nil || w.source == nil {
		return nil
	}
	return tlsconfig.MTLSServerConfig(w.source, w.source, tlsconfig.AuthorizeAny())
}

func (w *WorkloadIdentity) SVID() string {
	if w == nil || w.source == nil {
		return ""
	}
	svid, err := w.source.GetX509SVID()
	if err != nil {
		return ""
	}
	return svid.ID.String()
}

func (w *WorkloadIdentity) Close() error {
	if w == nil || w.source == nil {
		return nil
	}
	return w.source.Close()
}
