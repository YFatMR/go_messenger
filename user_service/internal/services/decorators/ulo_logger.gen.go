// Code generated by gowrap. DO NOT EDIT.
// template: ../../../../core/pkg/decorators/templates/ulo_logger.go
// gowrap: http://github.com/hexdigest/gowrap

package decorators

//go:generate gowrap gen -p github.com/YFatMR/go_messenger/user_service/internal/services -i UserService -t ../../../../core/pkg/decorators/templates/ulo_logger.go -o ulo_logger.gen.go -l ""

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/userid"
	"github.com/YFatMR/go_messenger/user_service/internal/services"
	"go.uber.org/zap"
)

// LoggingUserServiceDecorator implements services.UserService that is instrumented with custom zap logger
type LoggingUserServiceDecorator struct {
	logger *loggers.OtelZapLoggerWithTraceID
	base   services.UserService
}

// NewLoggingUserServiceDecorator instruments an implementation of the services.UserService with simple logging
func NewLoggingUserServiceDecorator(base services.UserService, logger *loggers.OtelZapLoggerWithTraceID) *LoggingUserServiceDecorator {
	if base == nil {
		panic("LoggingUserServiceDecorator got empty base")
	}
	if logger == nil {
		panic("LoggingUserServiceDecorator got empty logger")
	}
	return &LoggingUserServiceDecorator{
		base:   base,
		logger: logger,
	}
}

// Create implements services.UserService
func (d *LoggingUserServiceDecorator) Create(ctx context.Context, user *user.Entity, accountID *accountid.Entity) (userID *userid.Entity, logstash ulo.LogStash, err error) {
	d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: calling Create")
	defer func() {
		d.logger.LogContextNoExportULO(ctx, logstash)
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: Create finished")
	}()
	return d.base.Create(ctx, user, accountID)
}

// DeleteByID implements services.UserService
func (d *LoggingUserServiceDecorator) DeleteByID(ctx context.Context, userID *userid.Entity) (logstash ulo.LogStash, err error) {
	d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: calling DeleteByID")
	defer func() {
		d.logger.LogContextNoExportULO(ctx, logstash)
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: DeleteByID finished")
	}()
	return d.base.DeleteByID(ctx, userID)
}

// GetByID implements services.UserService
func (d *LoggingUserServiceDecorator) GetByID(ctx context.Context, userID *userid.Entity) (user *user.Entity, logstash ulo.LogStash, err error) {
	d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: calling GetByID")
	defer func() {
		d.logger.LogContextNoExportULO(ctx, logstash)
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: GetByID finished")
	}()
	return d.base.GetByID(ctx, userID)
}
