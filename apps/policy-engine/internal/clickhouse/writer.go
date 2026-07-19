package clickhouse

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/ClickHouse/clickhouse-go/v2/stdlib" // ClickHouse stdlib driver
    "go.uber.org/zap"
)

type EventRecord struct {
    ID                 string
    SessionID          string
    RequestID          string
    OrgID              string
    AppID              string
    ModelID            string
    Environment        string
    EventType          string
    Modality           string
    PayloadHash        string
    TokenCountEstimate int32
    PolicyVerdict      string
    PolicyRef          string
    Timestamp          time.Time
}

type DetectionRecord struct {
    ID             string
    EventID        string
    SessionID      string
    Category       string
    Confidence     float32
    ActionTaken    string
    OWASPLLMID     string
    ATLASTechnique string
    PolicyRef      string
    ModelAVerdict  string
    ModelBVerdict  string
    Disagreement   bool
    Timestamp      time.Time
}

type Writer struct {
    db     *sql.DB
    logger *zap.Logger
}

func NewWriter(ctx context.Context, addr, database string, logger *zap.Logger) (*Writer, error) {
    dsn := fmt.Sprintf("clickhouse://%s?database=%s", addr, database)
    db, err := sql.Open("clickhouse", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open clickhouse connection: %w", err)
    }

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping clickhouse: %w", err)
    }

    return &Writer{
        db:     db,
        logger: logger,
    }, nil
}

func (w *Writer) WriteEvent(ctx context.Context, rec EventRecord) error {
    query := `
        INSERT INTO aegis.events (
            id, session_id, request_id, org_id, app_id, model_id, environment,
            event_type, modality, payload_hash, token_count_estimate,
            policy_verdict, policy_ref, ts
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    _, err := w.db.ExecContext(ctx, query,
        rec.ID, rec.SessionID, rec.RequestID, rec.OrgID, rec.AppID, rec.ModelID, rec.Environment,
        rec.EventType, rec.Modality, rec.PayloadHash, rec.TokenCountEstimate,
        rec.PolicyVerdict, rec.PolicyRef, rec.Timestamp,
    )
    if err != nil {
        return fmt.Errorf("failed to write event: %w", err)
    }
    return nil
}

func (w *Writer) WriteDetection(ctx context.Context, rec DetectionRecord) error {
    query := `
        INSERT INTO aegis.detections (
            id, event_id, session_id, category, confidence, action_taken,
            owasp_llm_id, atlas_technique, policy_ref, model_a_verdict,
            model_b_verdict, disagreement, ts
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    _, err := w.db.ExecContext(ctx, query,
        rec.ID, rec.EventID, rec.SessionID, rec.Category, rec.Confidence, rec.ActionTaken,
        rec.OWASPLLMID, rec.ATLASTechnique, rec.PolicyRef, rec.ModelAVerdict,
        rec.ModelBVerdict, rec.Disagreement, rec.Timestamp,
    )
    if err != nil {
        return fmt.Errorf("failed to write detection: %w", err)
    }
    return nil
}

func (w *Writer) WriteBatch(ctx context.Context, events []EventRecord) error {
    if len(events) == 0 {
        return nil
    }

    tx, err := w.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO aegis.events (
            id, session_id, request_id, org_id, app_id, model_id, environment,
            event_type, modality, payload_hash, token_count_estimate,
            policy_verdict, policy_ref, ts
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()

    for _, rec := range events {
        _, err = stmt.ExecContext(ctx,
            rec.ID, rec.SessionID, rec.RequestID, rec.OrgID, rec.AppID, rec.ModelID, rec.Environment,
            rec.EventType, rec.Modality, rec.PayloadHash, rec.TokenCountEstimate,
            rec.PolicyVerdict, rec.PolicyRef, rec.Timestamp,
        )
        if err != nil {
            return fmt.Errorf("failed to exec statement for event %s: %w", rec.ID, err)
        }
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit batch: %w", err)
    }
    return nil
}

func (w *Writer) Close() error {
    return w.db.Close()
}
