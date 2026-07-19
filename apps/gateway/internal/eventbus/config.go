package eventbus

import "os"

// Config holds Kafka producer configuration from environment.
type Config struct {
    Brokers string // AEGIS_KAFKA_BROKERS (comma-separated)
    Topic   string // AEGIS_KAFKA_EVENTS_TOPIC (default: aegis.events.raw)
    Enabled bool   // false if AEGIS_KAFKA_BROKERS is empty
}

// ConfigFromEnv reads Kafka config from environment variables.
func ConfigFromEnv() Config {
    brokers := os.Getenv("AEGIS_KAFKA_BROKERS")
    topic := os.Getenv("AEGIS_KAFKA_EVENTS_TOPIC")
    if topic == "" {
        topic = "aegis.events.raw"
    }
    return Config{
        Brokers: brokers,
        Topic:   topic,
        Enabled: brokers != "",
    }
}
