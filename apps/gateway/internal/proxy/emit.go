package proxy

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "io"
    "net/http"
    "time"

    "github.com/aegis-security/aegis/apps/gateway/internal/eventbus"
    "github.com/google/uuid"
)

// emitRequestEvent constructs and emits a RawEvent from an incoming HTTP request.
// Payload content is NEVER included — only the hash.
func emitRequestEvent(
    ctx context.Context,
    producer *eventbus.Producer,
    r *http.Request,
    orgID, appID, modelID, environment, verdict, policyRef, requestID string,
) {
    if producer == nil {
        return // no bus configured (tests/dev)
    }

    sessionID := r.Header.Get("X-AEGIS-Session-ID")
    if sessionID == "" {
        sessionID = uuid.New().String()
    }

    // Read body for hashing only — do NOT store raw content
    payloadHash := ""
    if r.Body != nil {
        body, err := io.ReadAll(io.LimitReader(r.Body, 2<<20)) // max 2MB
        if err == nil {
            h := sha256.Sum256(body)
            payloadHash = hex.EncodeToString(h[:])
        }
        // NOTE: body must be re-set for the proxy to forward it in caller
    }

    event := eventbus.RawEventPayload{
        ID:                 uuid.New().String(),
        SessionID:          sessionID,
        RequestID:          requestID,
        OrgID:              orgID,
        AppID:              appID,
        ModelID:            modelID,
        Environment:        environment,
        EventType:          "request",
        Modality:           "text",
        PayloadHash:        payloadHash,
        TokenCountEstimate: 0, // populated by sanitizer in Phase 4
        PolicyVerdict:      verdict,
        PolicyRef:          policyRef,
        Metadata:           map[string]string{"method": r.Method, "path": r.URL.Path},
        Timestamp:          time.Now().UTC(),
    }

    producer.EmitAsync(ctx, event)
}
