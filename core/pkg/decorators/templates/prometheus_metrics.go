import (
		"context"
		"time"

		 
)

type databaseOperatinTag string

const (
	okStatusTag    = "ok"
	errorStatusTag = "error"
)

// Prefixes:
// database_query: for database
// service_request: for seervices

{{ $errorType := "logerr.Error" }}
{{ $metricPrefix := .Vars.metricPrefix }}

// Naming rule: https://prometheus.io/docs/practices/naming/
var (
	databaseQueryDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "{{ $metricPrefix }}_duration_seconds",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 4},
			Help:    "Duration of one database query",
	}, []string{"status", "operation"})
	databaseQueryProcessedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "{{ $metricPrefix }}_processed_total",
			Help: "Count of query to database",
	}, []string{"status", "operation"})
	databaseQueryStartProcessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "{{ $metricPrefix }}_start_process_total",
		Help: "Count of started queries",
	})
)

{{ $decorator := (or .Vars.DecoratorName (printf "PrometheusMetrics%sDecorator" .Interface.Name)) }}

// {{$decorator}} implements {{.Interface.Type}} that is instrumented with custom zap logger
type {{$decorator}} struct {
	base {{.Interface.Type}}
}

// New{{$decorator}} instruments an implementation of the {{.Interface.Type}} with simple logging
func New{{$decorator}}(base {{.Interface.Type}}) *{{$decorator}} {
	if base == nil {
		panic("{{$decorator}} got empty base")
	}
	return &{{$decorator}}{
		base: base,
	}
}

{{range $method := .Interface.Methods}}
	// {{$method.Name}} implements {{$.Interface.Type}}
	func (d *{{$decorator}}) {{$method.Declaration}} {
		{{ if not ($method.AcceptsContext) }}
		panic("Expected context variable")
		{{ break }}
	{{ end }}
	{{ if not ($method.HasResults) }}
		panic("Expected return type from method. {{ $errorType }}  at least")
		{{ break }}
	{{ end }}

	{{ $errorsResultCount := 0}}
	{{ range $result := $method.Results }}
		{{ if eq $result.Type $errorType }}
			{{ $errorsResultCount = add $errorsResultCount 1 }}
		{{ end }}
	{{ end }}

	{{ if ne $errorsResultCount 1 }}
		panic("Expected exact one {{ $errorType }} type as last argument")
		{{ break }}
	{{ end }}

	{{ $errorResult := last $method.Results }}
	{{ if not (eq $errorResult.Type $errorType) }}
		panic("Expected exact one {{ $errorType }} type as last argument")
		{{ break }}
	{{ end }}
		startTime := time.Now()
		databaseQueryStartProcessTotal.Inc()
		defer func () {
			functionDuration := time.Since(startTime).Seconds()
			statusTag := okStatusTag 
			if {{ $errorResult.Name }} != nil && {{ $errorResult.Name }}.HasError() {
				statusTag = errorStatusTag
			}
			databaseQueryDurationSeconds.WithLabelValues(statusTag, "{{ snakecase $method.Name }}").Observe(functionDuration)
			databaseQueryProcessedTotal.WithLabelValues(statusTag, "{{ snakecase $method.Name }}").Inc()
		}()
		{{ $method.Pass "d.base." -}}
	}
{{end}}
