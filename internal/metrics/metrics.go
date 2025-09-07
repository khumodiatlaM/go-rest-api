package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "go_rest_api",
	Name:      "http_requests_total",
	Help:      "Total number of HTTP requests",
}, []string{"path", "method", "status"})

var RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "go_rest_api",
	Name:      "http_request_duration_seconds",
	Help:      "Duration of HTTP requests",
	Buckets:   prometheus.DefBuckets,
}, []string{"path", "method"})
