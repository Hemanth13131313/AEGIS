package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PolicyChecksTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "aegis_policy_checks_total",
		Help: "Total policy checks performed",
	}, []string{"action", "scope_type"})

	OPAEvalDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "aegis_opa_eval_duration_seconds",
		Help:    "OPA evaluation duration in seconds",
		Buckets: []float64{0.001, 0.005, 0.010, 0.025, 0.050, 0.100},
	}, []string{"policy_id"})

	PolicyReloadsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "aegis_policy_reloads_total",
		Help: "Total policy hot-reload events",
	})

	ActivePoliciesGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "aegis_policy_active_gauge",
		Help: "Number of active policies currently loaded",
	})

	CacheHitsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "aegis_policy_cache_hits_total",
		Help: "Total policy cache hits and misses",
	}, []string{"result"}) // hit, miss
)
