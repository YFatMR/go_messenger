package decorators

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/metrics/prometheus"
)

// Naming rule: https://prometheus.io/docs/practices/naming/

const (
	InsertOperationTag = "insert"
	FindOperationTag   = "find"
	DeleteOperationTag = "find"
)

func CollectMetricForDatabaseCallback(ctx context.Context, operationTag prometheus.DatabaseOperatinTag,
	callback func(ctx context.Context) error,
) (err error) {
	startTime := time.Now()
	defer prometheus.CollectDatabaseQueryMetrics(startTime, operationTag, err)
	return callback(ctx)
}

func CollectMetricForDatabaseCallbackWithReturnType[T any](ctx context.Context,
	operationTag prometheus.DatabaseOperatinTag, callback func(ctx context.Context) (T, error),
) (_ T, err error) {
	startTime := time.Now()
	defer prometheus.CollectDatabaseQueryMetrics(startTime, operationTag, err)
	return callback(ctx)
}

func CollectMetricForDatabaseCallbackWithTwoReturnType[T any, G any](ctx context.Context,
	operationTag prometheus.DatabaseOperatinTag, callback func(ctx context.Context) (T, G, error),
) (_ T, _ G, err error) {
	startTime := time.Now()
	defer prometheus.CollectDatabaseQueryMetrics(startTime, operationTag, err)
	return callback(ctx)
}
