"""
Kafka consumer for RAG Monitor.
Consumes from aegis.events.rag and publishes detections to aegis.events.detections.
"""
import asyncio
import json
import uuid
from dataclasses import dataclass, field
from datetime import UTC, datetime

import structlog
from confluent_kafka import Consumer, KafkaError

from app.consumer.detection_producer import DetectionEventMessage, DetectionProducer
from app.monitor.rag_monitor import RAGMonitor

logger = structlog.get_logger(__name__)

@dataclass
class KafkaConsumerConfig:
    brokers: str
    group_id: str = "aegis-rag-monitor"
    topics: list[str] = field(default_factory=lambda: ["aegis.events.rag"])
    max_poll_interval_ms: int = 300000
    auto_offset_reset: str = "earliest"

@dataclass
class RAGEventMessage:
    id: str
    session_id: str
    request_id: str
    chunk_embeddings: list
    chunks: list[str]
    scores: list[float]
    ts: str

class RAGKafkaConsumer:
    def __init__(self, config: KafkaConsumerConfig, monitor: RAGMonitor, producer: DetectionProducer):
        self.config = config
        self.monitor = monitor
        self.producer = producer
        
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
        logger.info("rag_consumer_started", topics=self.config.topics)

    async def stop(self) -> None:
        self._running = False
        if self._task:
            self._task.cancel()
            try:
                await self._task
            except asyncio.CancelledError:
                pass
        self._consumer.close()
        logger.info("rag_consumer_stopped")

    async def _consume_loop(self) -> None:
        while self._running:
            try:
                msg = await asyncio.to_thread(self._consumer.poll, 0.1)
                
                if msg is None:
                    continue
                if msg.error():
                    if msg.error().code() == KafkaError._PARTITION_EOF:
                        continue
                    logger.error("kafka_consumer_error", error=str(msg.error()))
                    continue

                value = msg.value()
                if value is None:
                    continue
                
                try:
                    data = json.loads(value.decode('utf-8'))
                    event = RAGEventMessage(**data)
                    await self._process_event(event)
                except Exception as e:
                    logger.error("event_processing_error", error=str(e), key=msg.key())
                finally:
                    self._consumer.commit(message=msg)
                    
            except asyncio.CancelledError:
                break
            except Exception as e:
                logger.error("consume_loop_error", error=str(e))
                await asyncio.sleep(1)

    async def _process_event(self, event: RAGEventMessage) -> None:
        # Pass data to monitor
        is_anomalous = await self.monitor.analyze_retrieval_event(
            event.session_id,
            event.chunks,
            event.chunk_embeddings,
            event.scores
        )
        
        if is_anomalous:
            detection = DetectionEventMessage(
                id=str(uuid.uuid4()),
                event_id=event.id,
                session_id=event.session_id,
                category="rag_poisoning",
                confidence=0.85, # dynamic from monitor in real impl
                action_taken="tag",
                owasp_llm_id="LLM03",
                atlas_technique="AML.T0000",
                policy_ref="policy-rag",
                model_a_verdict="anomaly",
                model_b_verdict="none",
                disagreement=False,
                ts=datetime.now(UTC).isoformat()
            )
            self.producer.produce(detection)
            logger.info("rag_anomaly_detected", session_id=event.session_id)
