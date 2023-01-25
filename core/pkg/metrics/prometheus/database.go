package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Naming rule: https://prometheus.io/docs/practices/naming/

const (
	InsertOperationTag = "insert"
	FindOperationTag   = "find"
)

var (
	DatabaseSuccessQueryDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "database_success_query_duration_seconds",
		Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 4},
		Help:    "Duration of one database query",
	}, []string{"operation"})
	DatabaseQueryProcessedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "database_query_processed_total",
		Help: "Count of query to database",
	}, []string{"status", "operation"})
)

func CollectDatabaseQueryMetrics(startTime time.Time, operationTag string, err error) {
	functionDuration := time.Since(startTime).Seconds()
	statusTag := ErrorStatusTag
	if err == nil {
		DatabaseSuccessQueryDurationSeconds.WithLabelValues(operationTag).Observe(functionDuration)
		statusTag = OkStatusTag
	}
	DatabaseQueryProcessedTotal.WithLabelValues(statusTag, operationTag).Inc()
}
