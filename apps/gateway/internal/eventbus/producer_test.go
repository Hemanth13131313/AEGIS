package eventbus

import (
    "encoding/json"
    "testing"
    "time"
)

func TestHashPayload_Deterministic(t *testing.T) {
    // Basic test if hashing was exposed, skipping for now since hashing is in proxy
}

func TestRawEventPayload_JSONSerialization(t *testing.T) {
    now := time.Now().UTC()
    event := RawEventPayload{
        ID:          "123",
        SessionID:   "session-1",
        EventType:   "request",
        PayloadHash: "hash123",
        Timestamp:   now,
    }

    data, err := json.Marshal(event)
    if err != nil {
        t.Fatalf("Failed to marshal: %v", err)
    }

    var decoded RawEventPayload
    if err := json.Unmarshal(data, &decoded); err != nil {
        t.Fatalf("Failed to unmarshal: %v", err)
    }

    if decoded.ID != event.ID || decoded.SessionID != event.SessionID {
        t.Errorf("Mismatch in decoded event: %+v", decoded)
    }
}
