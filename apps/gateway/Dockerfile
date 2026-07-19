# syntax=docker/dockerfile:1.7
# Build stage
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /build

# Download dependencies first (better layer caching)
COPY go.mod go.sum* ./
RUN go mod download

# Copy source
COPY . .

# Build with hardening flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w -extldflags '-static'" \
    -o /app/aegis-gateway \
    ./cmd/main.go

# Verify binary
RUN file /app/aegis-gateway | grep -q 'statically linked' || true

# Final stage: distroless
FROM gcr.io/distroless/static-debian12:nonroot

# Labels (OCI standard)
LABEL org.opencontainers.image.title="AEGIS Gateway"
LABEL org.opencontainers.image.description="AEGIS AI Security Gateway Sidecar"
LABEL org.opencontainers.image.vendor="AEGIS Security"
LABEL org.opencontainers.image.licenses="Apache-2.0"

# Copy binary and CA certs only
COPY --from=builder /app/aegis-gateway /aegis-gateway
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Non-root user (distroless nonroot = uid 65532)
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/aegis-gateway"]
