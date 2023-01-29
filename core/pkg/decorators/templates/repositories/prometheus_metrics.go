import (
    "context"
    "time"

    "github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
    "github.com/YFatMR/go_messenger/user_service/internal/entities/account_id"
)

// Naming rule: https://prometheus.io/docs/practices/naming/

type databaseOperatinTag string

const (
    okStatusTag    = "ok"
    errorStatusTag = "error"
)

var (
    databaseQueryDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "database_success_query_duration_seconds",
        Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 4},
        Help:    "Duration of one database query",
    }, []string{"status", "operation"})
    databaseQueryProcessedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "database_query_processed_total",
        Help: "Count of query to database",
    }, []string{"status", "operation"})
    databaseQueryStartProcessTotal = promauto.NewCounter(prometheus.CounterOpts{
      Name: "database_query_start_process_total",
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
  return &{{$decorator}}{
    base: base,
  }
}

{{range $method := .Interface.Methods}}
  // {{$method.Name}} implements {{$.Interface.Type}}
  func (d *{{$decorator}}) {{$method.Declaration}} {
    {{ if not ($method.AcceptsContext) }}
      panic("Expected context variable")
    {{ end }}

    {{ if not ($method.HasResults) }}
      panic("Expected return type from method. cerrors.Error at least")
    {{ end }}

    {{ range $result := $method.Results }}
      {{ if eq $result.Type "cerrors.Error" }}
        startTime := time.Now()
        databaseQueryStartProcessTotal.Inc()
        defer func () {
          functionDuration := time.Since(startTime).Seconds()
          statusTag := errorStatusTag
          if {{ $result.Name }} == nil {
            statusTag = okStatusTag
          }
          databaseQueryDurationSeconds.WithLabelValues(statusTag, "{{ down $method.Name }}").Observe(functionDuration)
          databaseQueryProcessedTotal.WithLabelValues(statusTag, "{{ down $method.Name }}").Inc()
        }()
        {{ $method.Pass "d.base." }}
      {{ end }}
    {{ end }}
  }
{{end}}
