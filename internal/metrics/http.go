package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var HttpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP requests latency",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"path", "method", "status"},
)

func init() {
	prometheus.MustRegister(HttpRequestDuration)
}
