package eventbus

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.uber.org/zap"
)

type RawEventPayload struct {
    ID                  string            `json:"id"`
    SessionID           string            `json:"session_id"` // Kafka partition key
    RequestID           string            `json:"request_id"`
    OrgID               string            `json:"org_id"`
    AppID               string            `json:"app_id"`
    ModelID             string            `json:"model_id"`
    Environment         string            `json:"environment"`
    EventType           string            `json:"event_type"`  // request, response, tool_call
    Modality            string            `json:"modality"`    // text, image, audio
    PayloadHash         string            `json:"payload_hash"` // hex SHA-256 (never raw content)
    TokenCountEstimate  int               `json:"token_count_estimate"`
    PolicyVerdict       string            `json:"policy_verdict"`
    PolicyRef           string            `json:"policy_ref"`
    Metadata            map[string]string `json:"metadata"`
    Timestamp           time.Time         `json:"ts"`
}

type Producer struct {
    producer *kafka.Producer
    topic    string
    logger   *zap.Logger
}

func NewProducer(brokers, topic string, logger *zap.Logger) (*Producer, error) {
    p, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers":  brokers,
        "acks":               "all",
        "enable.idempotence": true,
        "compression.type":   "lz4",
        "message.timeout.ms": 5000,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create kafka producer: %w", err)
    }

    return &Producer{
        producer: p,
        topic:    topic,
        logger:   logger,
    }, nil
}

func (p *Producer) EmitRawEvent(ctx context.Context, event RawEventPayload) error {
    tracer := otel.Tracer("gateway-eventbus")
    ctx, span := tracer.Start(ctx, "eventbus.emit.raw_event")
    defer span.End()

    span.SetAttributes(
        attribute.String("event.type", event.EventType),
        attribute.String("event.session_id", event.SessionID),
    )

    data, err := json.Marshal(event)
    if err != nil {
        p.logger.Error("failed to marshal raw event", zap.Error(err))
        return fmt.Errorf("marshal error: %w", err)
    }

    deliveryChan := make(chan kafka.Event)
    err = p.producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
        Key:            []byte(event.SessionID),
        Value:          data,
    }, deliveryChan)
    if err != nil {
        p.logger.Error("failed to enqueue message", zap.Error(err))
        return fmt.Errorf("produce error: %w", err)
    }

    // Fire and forget - do not wait for deliveryChan in hot path
    go func() {
        e := <-deliveryChan
        m := e.(*kafka.Message)
        if m.TopicPartition.Error != nil {
            p.logger.Error("failed to deliver message", zap.Error(m.TopicPartition.Error))
        } else {
            p.logger.Debug("message delivered",
                zap.String("topic", *m.TopicPartition.Topic),
                zap.Int32("partition", m.TopicPartition.Partition),
                zap.String("session_id", event.SessionID),
            )
        }
    }()

    return nil
}

func (p *Producer) EmitAsync(ctx context.Context, event RawEventPayload) {
    if err := p.EmitRawEvent(ctx, event); err != nil {
        p.logger.Warn("async emit failed", zap.Error(err))
    }
}

func (p *Producer) Close() error {
    p.logger.Info("flushing kafka producer...")
    p.producer.Flush(5000)
    p.producer.Close()
    return nil
}
