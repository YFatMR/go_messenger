import (
    "context"

    "github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
    "go.opentelemetry.io/otel/trace"
    "github.com/YFatMR/go_messenger/user_service/internal/entities/account_id"
)

{{ $decorator := (or .Vars.DecoratorName (printf "OpentelemetryTracing%sDecorator" .Interface.Name)) }}

// {{$decorator}} implements {{.Interface.Type}} that is instrumented with custom zap logger
type {{$decorator}} struct {
  base {{.Interface.Type}}
  tracer trace.Tracer
  recordErrors bool
}

// New{{$decorator}} instruments an implementation of the {{.Interface.Type}} with simple logging
func New{{$decorator}}(base {{.Interface.Type}}, tracer trace.Tracer, recordErrors bool) *{{$decorator}} {
  return &{{$decorator}}{
    base: base,
    tracer: tracer,
    recordErrors: recordErrors,
  }
}

{{range $method := .Interface.Methods}}
  // {{$method.Name}} implements {{$.Interface.Type}}
  func (d *{{$decorator}}) {{$method.Declaration}} {
    {{ if not ($method.AcceptsContext) }}
      panic("Expected context variable")
    {{ end }}

    {{ if not ($method.HasResults) }}
      var span trace.Span
      ctx, span = d.tracer.Start(ctx, "/{{$method.Name}}")
      defer span.End()
      {{ $method.Pass "d.base." }}
    {{ end }}

    {{ range $result := $method.Results }}
        {{ if eq $result.Type "cerrors.Error" }}
          var span trace.Span
          ctx, span = d.tracer.Start(ctx, "/{{$method.Name}}")
          defer func() {
            if {{ $result.Name }} != nil && d.recordErrors {
                span.RecordError({{ $result.Name }}.GetInternalError())
            }
            span.End()
          }()
          {{ $method.Pass "d.base." }}
        {{ end }}
    {{ end }}
  }
{{end}}
