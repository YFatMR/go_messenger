// Code generated by gowrap. DO NOT EDIT.
// template: ../../../../core/pkg/decorators/templates/prometheus_metrics.go
// gowrap: http://github.com/hexdigest/gowrap

package decorators

//go:generate gowrap gen -p github.com/YFatMR/go_messenger/auth_service/internal/services -i AccountService -t ../../../../core/pkg/decorators/templates/prometheus_metrics.go -o prometheus_metrics.gen.go -v metricPrefix=service_request -l ""

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/token"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"
	"github.com/YFatMR/go_messenger/auth_service/internal/services"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type databaseOperatinTag string

const (
	okStatusTag    = "ok"
	errorStatusTag = "error"
)

// Prefixes:
// database_query: for database
// service_request: for seervices

// Naming rule: https://prometheus.io/docs/practices/naming/
var (
	durationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "service_request_duration_seconds",
		Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 4},
		Help:    "Duration of one database query",
	}, []string{"status", "operation"})
	processedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_request_processed_total",
		Help: "Count of query to database",
	}, []string{"status", "operation"})
	startProcessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "service_request_start_process_total",
		Help: "Count of started queries",
	})
)

// PrometheusMetricsAccountServiceDecorator implements services.AccountService that is instrumented with custom zap logger
type PrometheusMetricsAccountServiceDecorator struct {
	base services.AccountService
}

// NewPrometheusMetricsAccountServiceDecorator instruments an implementation of the services.AccountService with simple logging
func NewPrometheusMetricsAccountServiceDecorator(base services.AccountService) *PrometheusMetricsAccountServiceDecorator {
	if base == nil {
		panic("PrometheusMetricsAccountServiceDecorator got empty base")
	}
	return &PrometheusMetricsAccountServiceDecorator{
		base: base,
	}
}

// CreateAccount implements services.AccountService
func (d *PrometheusMetricsAccountServiceDecorator) CreateAccount(ctx context.Context, credential *credential.Entity) (accountID *accountid.Entity, logStash ulo.LogStash, err error) {
	startTime := time.Now()
	startProcessTotal.Inc()
	defer func() {
		functionDuration := time.Since(startTime).Seconds()
		statusTag := okStatusTag

		if err != nil {
			statusTag = errorStatusTag
		}

		durationSeconds.WithLabelValues(statusTag, "create_account").Observe(functionDuration)
		processedTotal.WithLabelValues(statusTag, "create_account").Inc()
	}()
	return d.base.CreateAccount(ctx, credential)
}

// GetToken implements services.AccountService
func (d *PrometheusMetricsAccountServiceDecorator) GetToken(ctx context.Context, credential *credential.Entity) (token *token.Entity, logStash ulo.LogStash, err error) {
	startTime := time.Now()
	startProcessTotal.Inc()
	defer func() {
		functionDuration := time.Since(startTime).Seconds()
		statusTag := okStatusTag

		if err != nil {
			statusTag = errorStatusTag
		}

		durationSeconds.WithLabelValues(statusTag, "get_token").Observe(functionDuration)
		processedTotal.WithLabelValues(statusTag, "get_token").Inc()
	}()
	return d.base.GetToken(ctx, credential)
}

// GetTokenPayload implements services.AccountService
func (d *PrometheusMetricsAccountServiceDecorator) GetTokenPayload(ctx context.Context, token *token.Entity) (tokenPayload *tokenpayload.Entity, logStash ulo.LogStash, err error) {
	startTime := time.Now()
	startProcessTotal.Inc()
	defer func() {
		functionDuration := time.Since(startTime).Seconds()
		statusTag := okStatusTag

		if err != nil {
			statusTag = errorStatusTag
		}

		durationSeconds.WithLabelValues(statusTag, "get_token_payload").Observe(functionDuration)
		processedTotal.WithLabelValues(statusTag, "get_token_payload").Inc()
	}()
	return d.base.GetTokenPayload(ctx, token)
}
