package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Naming rule: https://prometheus.io/docs/practices/naming/
var (
	ServiceRequestsProcessedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_request_processed_total",
		Help: "The total number of requests",
	}, []string{"status", "endpoint"})
)

func CollectServiceRequestMetrics(endpointTag string, err error) {
	statusTag := OkStatusTag
	if err != nil {
		statusTag = ErrorStatusTag
	}
	ServiceRequestsProcessedTotal.WithLabelValues(endpointTag, statusTag).Inc()
}
