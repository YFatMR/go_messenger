// Code generated by gowrap. DO NOT EDIT.
// template: ../../../../core/pkg/decorators/templates/opentelemetry_tracing.go
// gowrap: http://github.com/hexdigest/gowrap

package decorators

//go:generate gowrap gen -p github.com/YFatMR/go_messenger/user_service/internal/services -i UserService -t ../../../../core/pkg/decorators/templates/opentelemetry_tracing.go -o opentelemetry_tracing.gen.go -l ""

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/userid"
	"github.com/YFatMR/go_messenger/user_service/internal/services"
	"go.opentelemetry.io/otel/trace"
)

// OpentelemetryTracingUserServiceDecorator implements services.UserService that is instrumented with custom zap logger
type OpentelemetryTracingUserServiceDecorator struct {
	base         services.UserService
	tracer       trace.Tracer
	recordErrors bool
}

// NewOpentelemetryTracingUserServiceDecorator instruments an implementation of the services.UserService with simple logging
func NewOpentelemetryTracingUserServiceDecorator(base services.UserService, tracer trace.Tracer, recordErrors bool) *OpentelemetryTracingUserServiceDecorator {
	if base == nil {
		panic("OpentelemetryTracingUserServiceDecorator got empty base")
	}
	if tracer == nil {
		panic("OpentelemetryTracingUserServiceDecorator got empty tracer")
	}
	return &OpentelemetryTracingUserServiceDecorator{
		base:         base,
		tracer:       tracer,
		recordErrors: recordErrors,
	}
}

// Create implements services.UserService
func (d *OpentelemetryTracingUserServiceDecorator) Create(ctx context.Context, user *user.Entity, accountID *accountid.Entity) (userID *userid.Entity, logstash ulo.LogStash, err error) {
	var span trace.Span
	ctx, span = d.tracer.Start(ctx, "/Create")
	defer func() {
		if err != nil && d.recordErrors {
			span.RecordError(err)
		}
		span.End()
	}()
	return d.base.Create(ctx, user, accountID)
}

// DeleteByID implements services.UserService
func (d *OpentelemetryTracingUserServiceDecorator) DeleteByID(ctx context.Context, userID *userid.Entity) (logstash ulo.LogStash, err error) {
	var span trace.Span
	ctx, span = d.tracer.Start(ctx, "/DeleteByID")
	defer func() {
		if err != nil && d.recordErrors {
			span.RecordError(err)
		}
		span.End()
	}()
	return d.base.DeleteByID(ctx, userID)
}

// GetByID implements services.UserService
func (d *OpentelemetryTracingUserServiceDecorator) GetByID(ctx context.Context, userID *userid.Entity) (user *user.Entity, logstash ulo.LogStash, err error) {
	var span trace.Span
	ctx, span = d.tracer.Start(ctx, "/GetByID")
	defer func() {
		if err != nil && d.recordErrors {
			span.RecordError(err)
		}
		span.End()
	}()
	return d.base.GetByID(ctx, userID)
}