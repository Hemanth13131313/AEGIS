"""
AEGIS RAG Monitor — FastAPI entry point.
"""
import os
import uuid
from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager

import structlog
from fastapi import FastAPI, Request

from app.consumer.detection_producer import DetectionProducer
from app.consumer.kafka_consumer import KafkaConsumerConfig, RAGKafkaConsumer
from app.monitor import RAGMonitor

logger = structlog.get_logger(__name__)

monitor_instance = None
consumer_instance = None
producer_instance = None

@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    global monitor_instance, consumer_instance, producer_instance
    
    monitor_instance = RAGMonitor(anomaly_threshold=0.85)
    
    brokers = os.getenv("AEGIS_KAFKA_BROKERS")
    if brokers:
        producer_instance = DetectionProducer(brokers=brokers)
        config = KafkaConsumerConfig(brokers=brokers)
        consumer_instance = RAGKafkaConsumer(
            config=config,
            monitor=monitor_instance,
            producer=producer_instance
        )
        await consumer_instance.start()
        logger.info("rag_kafka_consumer_started")
    else:
        logger.warning("kafka_brokers_not_set_rag_consumer_disabled")

    yield

    if consumer_instance:
        await consumer_instance.stop()
    if producer_instance:
        producer_instance.close()

app = FastAPI(title="AEGIS RAG Monitor", lifespan=lifespan)

@app.middleware("http")
async def request_id_middleware(request: Request, call_next):
    req_id = request.headers.get("X-Request-ID", str(uuid.uuid4()))
    with structlog.contextvars.bound_contextvars(request_id=req_id):
        response = await call_next(request)
        response.headers["X-Request-ID"] = req_id
        return response



@app.get("/health")
async def health_check():
    return {
        "status": "ok", 
        "service": "rag-monitor",
        "kafka_enabled": consumer_instance is not None
    }
