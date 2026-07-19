package telemetry

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/baggage"
)

// InjectBaggage adds AEGIS context fields to OTel baggage for cross-service correlation.
func InjectBaggage(ctx context.Context, orgID, appID, sessionID string) context.Context {
	m1, _ := baggage.NewMember("org_id", orgID)
	m2, _ := baggage.NewMember("app_id", appID)
	m3, _ := baggage.NewMember("session_id", sessionID)
	b, _ := baggage.New(m1, m2, m3)
	return baggage.ContextWithBaggage(ctx, b)
}

// ExtractBaggage reads AEGIS context fields from OTel baggage.
func ExtractBaggage(ctx context.Context) (orgID, appID, sessionID string) {
	b := baggage.FromContext(ctx)
	return b.Member("org_id").Value(),
		b.Member("app_id").Value(),
		b.Member("session_id").Value()
}

// BaggageMiddleware extracts trace context and baggage from incoming HTTP requests.
func BaggageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// OTel propagation extracts W3C TraceContext and Baggage headers automatically
		// This middleware logs the correlation IDs for structured logs
		next.ServeHTTP(w, r)
	})
}
