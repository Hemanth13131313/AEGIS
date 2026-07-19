"""
Kafka consumer for AEGIS Scanner.
Consumes from aegis.events.raw and publishes detections to aegis.events.detections.
"""
import asyncio
import json
from dataclasses import dataclass, field
from datetime import datetime, timezone

from confluent_kafka import Consumer, KafkaError
import structlog

from app.ensemble.scanner import EnsembleScanner
from app.consumer.detection_producer import DetectionProducer, DetectionEventMessage

logger = structlog.get_logger(__name__)

@dataclass
class KafkaConsumerConfig:
    brokers: str
    group_id: str = "aegis-scanner"
    topics: list[str] = field(default_factory=lambda: ["aegis.events.raw"])
    max_poll_interval_ms: int = 300000
    auto_offset_reset: str = "earliest"

@dataclass
class RawEventMessage:
    id: str
    session_id: str
    request_id: str
    org_id: str
    app_id: str
    model_id: str
    environment: str
    event_type: str
    modality: str
    payload_hash: str
    token_count_estimate: int
    policy_verdict: str
    policy_ref: str
    metadata: dict
    ts: str

class ScannerKafkaConsumer:
    def __init__(self, config: KafkaConsumerConfig, scanner: EnsembleScanner, producer: DetectionProducer, logger_instance=None):
        self.config = config
        self.scanner = scanner
        self.producer = producer
        self.logger = logger_instance or logger
        
        self._consumer = Consumer({
            'bootstrap.servers': self.config.brokers,
            'group.id': self.config.group_id,
            'auto.offset.reset': self.config.auto_offset_reset,
            'enable.auto.commit': False,
            'max.poll.interval.ms': self.config.max_poll_interval_ms
        })
        self._running = False
        self._task = None

    async def start(self) -> None:
        self._consumer.subscribe(self.config.topics)
        self._running = True
        self._task = asyncio.create_task(self._consume_loop())
        self.logger.info("scanner_consumer_started", topics=self.config.topics)

    async def stop(self) -> None:
        self._running = False
        if self._task:
            self._task.cancel()
            try:
                await self._task
            except asyncio.CancelledError:
                pass
        self._consumer.close()
        self.logger.info("scanner_consumer_stopped")

    async def _consume_loop(self) -> None:
        while self._running:
            try:
                # Use executor to avoid blocking the event loop
                msg = await asyncio.to_thread(self._consumer.poll, 0.1)
                
                if msg is None:
                    continue
                if msg.error():
                    if msg.error().code() == KafkaError._PARTITION_EOF:
                        continue
                    self.logger.error("kafka_consumer_error", error=str(msg.error()))
                    continue

                value = msg.value()
                if value is None:
                    continue
                
                try:
                    data = json.loads(value.decode('utf-8'))
                    event = RawEventMessage(**data)
                    await self._process_event(event)
                except Exception as e:
                    self.logger.error("event_processing_error", error=str(e), key=msg.key())
                finally:
                    self._consumer.commit(message=msg)
                    
            except asyncio.CancelledError:
                break
            except Exception as e:
                self.logger.error("consume_loop_error", error=str(e))
                await asyncio.sleep(1)

    async def _process_event(self, event: RawEventMessage) -> None:
        # payload is not forwarded via Kafka for privacy; scanning metadata/hash or fetched payload later
        verdict_result = await self.scanner.scan(payload="", context=event.metadata)
        
        if verdict_result["verdict"] in ["block", "tag", "redact"]:
            detection = DetectionEventMessage(
                id=event.id, # using event ID or new UUID
                event_id=event.id,
                session_id=event.session_id,
                category=verdict_result.get("category", "unknown"),
                confidence=verdict_result.get("confidence", 1.0),
                action_taken=verdict_result["verdict"],
                owasp_llm_id="LLM01", # mock for now
                atlas_technique="AML.T0000", # mock
                policy_ref=event.policy_ref,
                model_a_verdict=verdict_result.get("model_a_verdict", "unknown"),
                model_b_verdict=verdict_result.get("model_b_verdict", "unknown"),
                disagreement=False,
                ts=datetime.now(timezone.utc).isoformat()
            )
            self.producer.produce(detection)
            self.logger.info(
                "detection_produced",
                session_id=event.session_id,
                request_id=event.request_id,
                action=detection.action_taken,
                confidence=detection.confidence
            )
