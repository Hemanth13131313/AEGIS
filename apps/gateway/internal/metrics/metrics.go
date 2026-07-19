package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests processed by the gateway",
	}, []string{"method", "path", "status", "service"})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: []float64{0.001, 0.005, 0.010, 0.025, 0.050, 0.100, 0.250, 0.500, 1.0},
	}, []string{"method", "path", "service"})

	PolicyCheckDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "aegis_gateway_policy_check_duration_seconds",
		Help:    "Time spent waiting for policy engine response",
		Buckets: []float64{0.001, 0.002, 0.005, 0.010, 0.025},
	}, []string{"action"})

	AuthFailuresTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "aegis_gateway_auth_failures_total",
		Help: "Total authentication failures",
	}, []string{"reason"})

	EventBusPublishErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "aegis_eventbus_publish_errors_total",
		Help: "Total Kafka publish errors (fire-and-forget, non-blocking)",
	})

	SanitizeRejectionsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "aegis_sanitize_rejections_total",
		Help: "Total requests rejected by input sanitizer",
	}, []string{"reason"})
)

// statusRecorder embeds http.ResponseWriter to track the status code.
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// Middleware returns an http.Handler middleware that records request metrics.
func Middleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" || r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(recorder, r)
			
			duration := time.Since(start).Seconds()
			
			RequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(recorder.statusCode), serviceName).Inc()
			RequestDuration.WithLabelValues(r.Method, r.URL.Path, serviceName).Observe(duration)
		})
	}
}

// Handler returns the promhttp handler for the /metrics endpoint.
func Handler() http.Handler {
	return promhttp.Handler()
}
