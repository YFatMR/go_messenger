import (
		"context"

		"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
		"go.opentelemetry.io/otel/trace"
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
		{{ break }}
	{{ end }}

	{{ if not ($method.HasResults) }}
		var span trace.Span
		ctx, span = d.tracer.Start(ctx, "/{{$method.Name}}")
		defer span.End()
		{{ $method.Pass "d.base." -}}
		{{ break }}
	{{ end }}

	{{ $errorsResultCount := 0}}
	{{ range $result := $method.Results }}
		{{ if eq $result.Type "cerrors.Error" }}
			{{ $errorsResultCount = add $errorsResultCount 1 }}
		{{ end }}
	{{ end }}

	{{ if ne $errorsResultCount 1 }}
		panic("Expected exact one cerrors.Error type as last argument")
		{{ break }}
	{{ end }}

	{{ $errorResult := last $method.Results }}
	{{ if not (eq $errorResult.Type "cerrors.Error") }}
		panic("Expected exact one cerrors.Error type as last argument")
		{{ break }}
	{{ end }}
		var span trace.Span
		ctx, span = d.tracer.Start(ctx, "/{{$method.Name}}")
		defer func() {
			if {{ $errorResult.Name }} != nil && d.recordErrors {
					span.RecordError({{ $errorResult.Name }}.GetInternalError())
			}
			span.End()
		}()
		{{ $method.Pass "d.base." -}}
	}
{{end}}
