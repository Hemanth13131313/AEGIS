from prometheus_client import Counter, Histogram, Gauge

# Counters
scan_requests_total = Counter(
    'aegis_scanner_requests_total',
    'Total scan requests processed',
    ['action', 'pre_filter_triggered']
)

# Histograms  
scan_latency_seconds = Histogram(
    'aegis_scanner_latency_seconds',
    'Scan latency in seconds',
    ['stage'],  # pre_filter, model_a, model_b, ensemble, total
    buckets=[0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5]
)

# Gauges
ensemble_disagreement_ratio = Gauge(
    'aegis_scanner_disagreement_ratio',
    'Rolling ratio of disagreements between ensemble models'
)

def record_scan(action: str, pre_filter: bool, latency: float):
    scan_requests_total.labels(action=action, pre_filter_triggered=str(pre_filter).lower()).inc()
    scan_latency_seconds.labels(stage="total").observe(latency)

def record_disagreement(disagreement: bool):
    # This is a bit simplified for a Gauge, typically you'd calculate a rolling ratio elsewhere,
    # but for this requirement we can just set it to 1 or 0 for the most recent observation,
    # or maintain state. For simplicity:
    ensemble_disagreement_ratio.set(1.0 if disagreement else 0.0)
