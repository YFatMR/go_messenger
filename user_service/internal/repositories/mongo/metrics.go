package mongo

import (
	"github.com/YFatMR/go_messenger/user_service/internal/metrics/prometheus"
	"time"
)

func collectDatabaseQueryMetrics(startTime time.Time, operationTag string, err *error) {
	functionDuration := time.Since(startTime).Seconds()
	statusTag := prometheus.ErrorStatusTag
	if *err == nil {
		prometheus.DatabaseSuccessQueryDurationSeconds.WithLabelValues(operationTag).Observe(functionDuration)
		statusTag = prometheus.OkStatusTag
	}
	prometheus.DatabaseQueryProcessedTotal.WithLabelValues(statusTag, operationTag).Inc()
}
