from contextlib import asynccontextmanager
from fastapi import FastAPI
from app.api.routes import router as scan_router
from app.api.metrics_route import router as metrics_router
from app.ensemble.scanner import EnsembleScanner
from app.evaluators.prompt_injection import PromptInjectionEvaluator
from app.evaluators.jailbreak import JailbreakEvaluator
from app.evaluators.data_exfiltration import DataExfiltrationEvaluator

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Setup
    app.state.scanner = EnsembleScanner.from_env()
    app.state.pre_filters = [
        PromptInjectionEvaluator(),
        JailbreakEvaluator(),
        DataExfiltrationEvaluator()
    ]
    # Keep Kafka consumer init from Phase 3 here if applicable
    yield
    # Teardown

app = FastAPI(title="AEGIS Scanner", lifespan=lifespan)

app.include_router(scan_router, prefix="/api/v1")
app.include_router(metrics_router)

@app.get("/health")
async def health():
    return {
        "status": "ok",
        "scanner_ready": hasattr(app.state, "scanner") and app.state.scanner is not None,
        "pre_filters": len(app.state.pre_filters) if hasattr(app.state, "pre_filters") else 0
    }
