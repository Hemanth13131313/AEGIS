module github.com/aegis-security/aegis/apps/gateway

go 1.22

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.4.0
	github.com/go-chi/chi/v5 v5.1.0
	github.com/lestrrat-go/jwx/v2 v2.1.0
	github.com/prometheus/client_golang v1.19.1
	github.com/redis/go-redis/v9 v9.5.1
	go.opentelemetry.io/otel v1.27.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.27.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.2
)
