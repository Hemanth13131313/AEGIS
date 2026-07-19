import pytest
from unittest.mock import AsyncMock, MagicMock, patch
from app.consumer.kafka_consumer import ScannerKafkaConsumer, KafkaConsumerConfig, RawEventMessage

@pytest.fixture
def mock_scanner():
    scanner = AsyncMock()
    return scanner

@pytest.fixture
def mock_producer():
    producer = MagicMock()
    return producer

@pytest.mark.asyncio
async def test_process_event_allow_verdict(mock_scanner, mock_producer):
    mock_scanner.scan.return_value = {"verdict": "allow", "confidence": 0.9}
    consumer = ScannerKafkaConsumer(
        config=KafkaConsumerConfig(brokers="localhost:9092"),
        scanner=mock_scanner,
        producer=mock_producer
    )
    
    event = RawEventMessage(
        id="123", session_id="s1", request_id="r1", org_id="o1", app_id="a1",
        model_id="m1", environment="test", event_type="request", modality="text",
        payload_hash="hash", token_count_estimate=10, policy_verdict="allow",
        policy_ref="p1", metadata={}, ts="2023-01-01T00:00:00Z"
    )
    
    await consumer._process_event(event)
    mock_producer.produce.assert_not_called()

@pytest.mark.asyncio
async def test_process_event_block_verdict(mock_scanner, mock_producer):
    mock_scanner.scan.return_value = {"verdict": "block", "confidence": 0.95, "category": "jailbreak"}
    consumer = ScannerKafkaConsumer(
        config=KafkaConsumerConfig(brokers="localhost:9092"),
        scanner=mock_scanner,
        producer=mock_producer
    )
    
    event = RawEventMessage(
        id="123", session_id="s1", request_id="r1", org_id="o1", app_id="a1",
        model_id="m1", environment="test", event_type="request", modality="text",
        payload_hash="hash", token_count_estimate=10, policy_verdict="allow",
        policy_ref="p1", metadata={}, ts="2023-01-01T00:00:00Z"
    )
    
    await consumer._process_event(event)
    mock_producer.produce.assert_called_once()
    
    called_detection = mock_producer.produce.call_args[0][0]
    assert called_detection.action_taken == "block"
    assert called_detection.category == "jailbreak"

@pytest.mark.asyncio
async def test_stop_cancels_loop(mock_scanner, mock_producer):
    consumer = ScannerKafkaConsumer(
        config=KafkaConsumerConfig(brokers="localhost:9092"),
        scanner=mock_scanner,
        producer=mock_producer
    )
    
    with patch('app.consumer.kafka_consumer.Consumer') as MockConfluentConsumer:
        mock_instance = MockConfluentConsumer.return_value
        await consumer.start()
        assert consumer._running is True
        
        await consumer.stop()
        assert consumer._running is False
        mock_instance.close.assert_called_once()
