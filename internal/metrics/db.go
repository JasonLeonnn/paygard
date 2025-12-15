package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var DbQueryDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "db_query_duration_seconds",
		Help:    "Database query latency",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"query"},
)

func init() {
	prometheus.MustRegister(DbQueryDuration)
}
