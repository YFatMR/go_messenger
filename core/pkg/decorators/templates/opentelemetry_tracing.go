import (
		"context"

		 
		"go.opentelemetry.io/otel/trace"
)

{{ $errorType := "logerr.Error" }}
{{ $decorator := (or .Vars.DecoratorName (printf "OpentelemetryTracing%sDecorator" .Interface.Name)) }}

// {{$decorator}} implements {{.Interface.Type}} that is instrumented with custom zap logger
type {{$decorator}} struct {
	base {{.Interface.Type}}
	tracer trace.Tracer
	recordErrors bool
}

// New{{$decorator}} instruments an implementation of the {{.Interface.Type}} with simple logging
func New{{$decorator}}(base {{.Interface.Type}}, tracer trace.Tracer, recordErrors bool) *{{$decorator}} {
	if base == nil {
		panic("{{$decorator}} got empty base")
	}
	if tracer == nil {
		panic("{{$decorator}} got empty tracer")
	}
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
		var span trace.Span
		ctx, span = d.tracer.Start(ctx, "/{{$method.Name}}")
		defer func() {
			defer span.End()
			if {{ $errorResult.Name }} == nil {
				return
			}
			if {{ $errorResult.Name }}.HasError() && d.recordErrors {
				span.RecordError({{ $errorResult.Name }}.GetAPIError())
			}
		}()
		{{ $method.Pass "d.base." -}}
	}
{{end}}
