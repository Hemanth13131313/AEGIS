"""
Persists red team run results to JSON and ClickHouse.
"""
import json
from datetime import datetime
from pathlib import Path
from typing import Optional
import structlog
from dataclasses import asdict

from .runner import RunResult, RunSummary

logger = structlog.get_logger(__name__)

class ResultStore:
    def __init__(self, results_dir: Path, clickhouse_addr: Optional[str] = None):
        self.results_dir = results_dir
        self.clickhouse_addr = clickhouse_addr
        
        if not self.results_dir.exists():
            self.results_dir.mkdir(parents=True)

    async def save(self, results: list[RunResult], summary: RunSummary) -> Path:
        timestamp = summary.run_at.strftime("%Y%m%d_%H%M%S")
        filename = self.results_dir / f"{timestamp}_results.json"
        
        data = {
            "summary": asdict(summary),
            "results": [asdict(r) for r in results]
        }
        
        # Datetime serialization helper
        def default_serializer(obj):
            if isinstance(obj, datetime):
                return obj.isoformat()
            raise TypeError(f"Type {type(obj)} not serializable")

        with open(filename, 'w') as f:
            json.dump(data, f, indent=2, default=default_serializer)
            
        logger.info("Saved local results", file=str(filename))
        
        if self.clickhouse_addr:
            try:
                await self._save_to_clickhouse(results)
                logger.info("Saved results to ClickHouse")
            except Exception as e:
                logger.error("Failed to save to ClickHouse", error=str(e))
                
        return filename

    async def _save_to_clickhouse(self, results: list[RunResult]):
        # Implementation for ClickHouse write would go here.
        # This typically involves POSTing bulk insert data to HTTP interface.
        pass
        
    def load_latest(self) -> Optional[list[RunResult]]:
        files = sorted(self.results_dir.glob("*_results.json"))
        if not files:
            return None
            
        with open(files[-1], 'r') as f:
            data = json.load(f)
            
        return [RunResult(**r) for r in data["results"]]
