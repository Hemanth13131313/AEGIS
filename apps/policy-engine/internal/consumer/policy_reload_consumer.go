package consumer

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
    "go.uber.org/zap"
)

type PolicyReloadMessage struct {
    PolicyID    string    `json:"policy_id"`
    ScopeID     string    `json:"scope_id"`
    ScopeType   string    `json:"scope_type"`
    TriggeredBy string    `json:"triggered_by"`
    Timestamp   time.Time `json:"ts"`
}

type PolicyLoader interface {
    ReloadPolicy(ctx context.Context, policyID string) error
}

type PolicyReloadConsumer struct {
    consumer *kafka.Consumer
    logger   *zap.Logger
    loader   PolicyLoader
}

func NewPolicyReloadConsumer(brokers, groupID string, loader PolicyLoader, logger *zap.Logger) (*PolicyReloadConsumer, error) {
    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": brokers,
        "group.id":          groupID,
        "auto.offset.reset": "earliest",
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create consumer: %w", err)
    }

    return &PolicyReloadConsumer{
        consumer: c,
        logger:   logger,
        loader:   loader,
    }, nil
}

func (c *PolicyReloadConsumer) Start(ctx context.Context) error {
    if err := c.consumer.Subscribe("aegis.control.policy-reload", nil); err != nil {
        return fmt.Errorf("failed to subscribe: %w", err)
    }

    c.logger.Info("started policy reload consumer")

    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
                if err != nil {
                    if err.(kafka.Error).Code() != kafka.ErrTimedOut {
                        c.logger.Warn("kafka read error", zap.Error(err))
                    }
                    continue
                }

                var cmd PolicyReloadMessage
                if err := json.Unmarshal(msg.Value, &cmd); err != nil {
                    c.logger.Error("failed to unmarshal policy reload message", zap.Error(err))
                    continue
                }

                c.logger.Info("received policy reload command", zap.String("policy_id", cmd.PolicyID))
                if err := c.loader.ReloadPolicy(ctx, cmd.PolicyID); err != nil {
                    c.logger.Error("failed to reload policy", zap.String("policy_id", cmd.PolicyID), zap.Error(err))
                } else {
                    c.logger.Info("policy reloaded successfully", zap.String("policy_id", cmd.PolicyID))
                }
            }
        }
    }()

    return nil
}

func (c *PolicyReloadConsumer) Close() error {
    return c.consumer.Close()
}
