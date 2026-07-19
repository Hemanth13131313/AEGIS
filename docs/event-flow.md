# AEGIS Event Flow Architecture

## Synchronous Path (P95 ≤ 10ms budget)
```
Client → Envoy TLS → Gateway (auth + sanitize + policy check) → AI Backend
                                           ↓ (fire-and-forget, non-blocking)
                                       Kafka Producer
```

## Async Analysis Path
```
Kafka: aegis.events.raw
    ↓ Scanner Consumer (ensemble scan, ~500ms)
    ↓ Kafka: aegis.events.detections
    ↓ ClickHouse writer

Kafka: aegis.events.rag  
    ↓ RAG Monitor Consumer (anomaly detection, ~200ms)
    ↓ Kafka: aegis.events.detections

Kafka: aegis.control.policy-reload
    ↓ Policy Engine Consumer → OPA hot-reload
```

## Topic Reference
| Topic | Producer | Consumer | Retention |
|---|---|---|---|
| aegis.events.raw | Gateway | Scanner, RAG Monitor | 7d |
| aegis.events.detections | Scanner, RAG Monitor | ClickHouse Writer, UI | 30d |
| aegis.events.rag | Gateway | RAG Monitor | 7d |
| aegis.control.policy-reload | Policy Engine API | Policy Engine Consumer | 1d |
| aegis.redteam.jobs | Red Team Runner | Scanner | 1d |

## Delivery Guarantees
- Producer: `acks=all`, `enable.idempotence=true` — exactly-once delivery to broker
- Consumers: at-least-once with manual offset commit after processing
- All consumers MUST be idempotent (upsert by event.id)
