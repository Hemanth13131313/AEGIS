# AEGIS Scanner

## Purpose
The dual-model ensemble scanner acts as the core detection engine for AEGIS, assessing incoming LLM payloads for prompt injection, jailbreaks, and data exfiltration.

## Architecture
The scanner pipeline consists of two stages:
1. **Pre-filter (Rule-based)**: Fast regex matching (<1ms) to catch known jailbreak/injection patterns instantly.
2. **Ensemble Scanner**: Two LLMs (Model A and Model B) running in parallel via `asyncio.gather`. They evaluate payloads using strictly isolated structures to prevent meta-injection. Disagreements result in a `tag` action.

## Environment Variables
| Variable | Description | Default |
|----------|-------------|---------|
| `AEGIS_SCANNER_MODEL_A_URL` | Base URL for Model A | `http://localhost:8000/v1` |
| `AEGIS_SCANNER_MODEL_A_NAME` | Model name for A | `mistral-7b-instruct` |
| `AEGIS_SCANNER_MODEL_B_URL` | Base URL for Model B | `http://localhost:8000/v1` |
| `AEGIS_SCANNER_MODEL_B_NAME` | Model name for B | `llama-3-8b-instruct` |
| `AEGIS_SCANNER_API_KEY` | API Key for external endpoints | `""` (empty for local) |
| `AEGIS_KAFKA_BROKERS` | Kafka broker list | `localhost:9092` |
| `AEGIS_KAFKA_EVENTS_TOPIC` | Kafka topic for events | `aegis.events` |

## Local vLLM Setup
1. Install vLLM: `pip install vllm`
2. Start model A: `python -m vllm.entrypoints.openai.api_server --model mistralai/Mistral-7B-Instruct-v0.2 --port 8000`
3. Start model B: `python -m vllm.entrypoints.openai.api_server --model meta-llama/Meta-Llama-3-8B-Instruct --port 8001`

## Tests
Run tests using: `uv run pytest`

## Metrics
Metrics are available at the `/metrics` endpoint in Prometheus format.

## Milestones
- **M4.1**: Pre-filter successfully blocks known injection payloads.
- **M4.2**: Dual-model ensemble handles disagreements correctly with `tag` actions.
- **M4.3**: RAG anomaly detection statistically flags poisoned retrievals.
