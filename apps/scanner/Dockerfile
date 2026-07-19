# syntax=docker/dockerfile:1.7
# Build stage
FROM python:3.12-slim AS builder

RUN pip install uv
WORKDIR /build

COPY pyproject.toml ./
RUN uv pip install --system --no-cache-dir -r pyproject.toml 2>/dev/null || \
    uv pip install --system --no-cache-dir fastapi uvicorn pydantic structlog

COPY app/ ./app/

# Final stage: distroless Python
FROM gcr.io/distroless/python3-debian12:nonroot

LABEL org.opencontainers.image.title="AEGIS Scanner"
LABEL org.opencontainers.image.description="AEGIS AI Dual-Model Ensemble Scanner"
LABEL org.opencontainers.image.vendor="AEGIS Security"
LABEL org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /usr/local/lib/python3.12 /usr/local/lib/python3.12
COPY --from=builder /usr/local/bin/uvicorn /usr/local/bin/uvicorn
COPY --from=builder /build/app /app

USER nonroot:nonroot
WORKDIR /app

EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
