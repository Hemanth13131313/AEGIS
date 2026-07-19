package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

var ErrCacheMiss = errors.New("cache miss")

type PolicyCache struct {
	addr   string
	logger *zap.Logger
}

func NewPolicyCache(addr string, logger *zap.Logger) (*PolicyCache, error) {
	return &PolicyCache{
		addr:   addr,
		logger: logger,
	}, nil
}

func (c *PolicyCache) SetPolicy(ctx context.Context, key string, regoContent string, ttl time.Duration) error {
	// Stub for redis SET
	return nil
}

func (c *PolicyCache) GetPolicy(ctx context.Context, key string) (string, error) {
	// Stub for redis GET
	return "", ErrCacheMiss
}

func (c *PolicyCache) InvalidatePolicy(ctx context.Context, key string) error {
	return nil
}

func (c *PolicyCache) Ping(ctx context.Context) error {
	if c.addr == "" {
		return fmt.Errorf("redis not configured")
	}
	return nil
}
