package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TransactionCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transactions_total",
			Help: "Total number of transactions created",
		},
		[]string{"category"},
	)

	AnomalyCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "anomalies_total",
			Help: "Total number of anomaly transactions",
		},
	)

	AnomalyBySeverityCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "anomalies_by_severity_total",
			Help: "Total number of anomaly transactions by severity",
		},
		[]string{"severity"},
	)
)

func init() {
	prometheus.MustRegister(TransactionCounter)
	prometheus.MustRegister(AnomalyCounter)
	prometheus.MustRegister(AnomalyBySeverityCounter)
}
