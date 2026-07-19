import json
from dataclasses import asdict, dataclass

import structlog
from confluent_kafka import Producer

logger = structlog.get_logger(__name__)

@dataclass
class DetectionEventMessage:
    id: str
    event_id: str
    session_id: str
    category: str
    confidence: float
    action_taken: str
    owasp_llm_id: str
    atlas_technique: str
    policy_ref: str
    model_a_verdict: str
    model_b_verdict: str
    disagreement: bool
    ts: str

class DetectionProducer:
    def __init__(self, brokers: str, topic: str = "aegis.events.detections"):
        self.topic = topic
        self._producer = Producer({
            'bootstrap.servers': brokers,
            'acks': 'all',
            'compression.type': 'lz4'
        })

    def produce(self, detection: DetectionEventMessage) -> None:
        try:
            json_bytes = json.dumps(asdict(detection)).encode('utf-8')
            self._producer.produce(
                self.topic,
                key=detection.session_id.encode('utf-8') if detection.session_id else None,
                value=json_bytes,
                on_delivery=self._delivery_callback
            )
            self._producer.poll(0)
        except Exception as e:
            logger.error("failed_to_produce_detection", error=str(e), session_id=detection.session_id)

    def _delivery_callback(self, err, msg) -> None:
        if err:
            logger.error("detection_delivery_failed", error=str(err))
        else:
            logger.debug("detection_delivered", topic=msg.topic(), partition=msg.partition())

    def flush(self) -> None:
        self._producer.flush(timeout=5.0)

    def close(self) -> None:
        self.flush()
