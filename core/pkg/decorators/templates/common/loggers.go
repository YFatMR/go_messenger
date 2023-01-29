import (
  "context"

  "github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
  "github.com/YFatMR/go_messenger/core/pkg/loggers"
  "go.uber.org/zap"
  "github.com/YFatMR/go_messenger/user_service/internal/entities/account_id"
)


{{ $decorator := (or .Vars.DecoratorName (printf "Logging%sDecorator" .Interface.Name)) }}

// {{$decorator}} implements {{.Interface.Type}} that is instrumented with custom zap logger
type {{$decorator}} struct {
  logger *loggers.OtelZapLoggerWithTraceID
  base {{.Interface.Type}}
}

// New{{$decorator}} instruments an implementation of the {{.Interface.Type}} with simple logging
func New{{$decorator}}(base {{.Interface.Type}}, logger *loggers.OtelZapLoggerWithTraceID) *{{$decorator}} {
  return &{{$decorator}}{
    base: base,
    logger: logger,
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
          d.logger.DebugContextNoExport(ctx, "{{$decorator}}: calling {{$method.Name}}")
          defer func() {
            {{- range $result := $method.Results -}}
                {{- if eq $result.Type "cerrors.Error" -}}
                  if {{ $result.Name }} != nil {
                    d.logger.ErrorContext(
                      ctx, {{ $result.Name }}.GetInternalErrorMessage(), zap.Error({{ $result.Name }}.GetInternalError()),
                    )
                  } else {
                    d.logger.DebugContextNoExport(ctx, "{{$decorator}}: {{$method.Name}} finished")
                  }
                {{- end}}
            {{ end }}
            d.logger.DebugContextNoExport(ctx, "{{$decorator}}: {{$method.Name}} finished")
          }()
          {{ $method.Pass "d.base." }}
        {{ end }}
    {{ end }}
  }
{{end}}
