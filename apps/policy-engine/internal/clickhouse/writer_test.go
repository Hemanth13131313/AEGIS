package clickhouse

import (
    "context"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "go.uber.org/zap"
)

func TestWriteEvent_MapsFieldsCorrectly(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    logger := zap.NewNop()
    writer := &Writer{
        db:     db,
        logger: logger,
    }

    rec := EventRecord{
        ID:                 "123",
        SessionID:          "s1",
        RequestID:          "r1",
        OrgID:              "o1",
        AppID:              "a1",
        ModelID:            "m1",
        Environment:        "dev",
        EventType:          "request",
        Modality:           "text",
        PayloadHash:        "hash",
        TokenCountEstimate: 10,
        PolicyVerdict:      "allow",
        PolicyRef:          "p1",
        Timestamp:          time.Now(),
    }

    mock.ExpectExec("INSERT INTO aegis.events").
        WithArgs(rec.ID, rec.SessionID, rec.RequestID, rec.OrgID, rec.AppID, rec.ModelID, rec.Environment,
            rec.EventType, rec.Modality, rec.PayloadHash, rec.TokenCountEstimate,
            rec.PolicyVerdict, rec.PolicyRef, rec.Timestamp).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = writer.WriteEvent(context.Background(), rec)
    if err != nil {
        t.Errorf("error was not expected while inserting: %s", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}

func TestWriteBatch_HandlesEmpty(t *testing.T) {
    writer := &Writer{
        db:     nil, // shouldn't be accessed
        logger: zap.NewNop(),
    }

    err := writer.WriteBatch(context.Background(), []EventRecord{})
    if err != nil {
        t.Errorf("Expected nil error for empty batch, got %v", err)
    }
}
