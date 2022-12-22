package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Naming rule: https://prometheus.io/docs/practices/naming/

// Common tags
const (
	OkStatusTag    = "ok"
	ErrorStatusTag = "error"
)

// Metrics for high level endpoints (gRPC endpoints)
const (
	ServerSideErrorRequestTag = "server_side"
	ClientSideErrorRequestTag = "user_side"
)

var (
	RequestProcessingTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_processing_total",
		Help: "The total number of requests",
	}, []string{"endpoint"})
	RequestProcessingErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_processing",
		Help: "The total number of error requests",
	}, []string{"endpoint", "error_type"})
)

// Metrics for database
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
