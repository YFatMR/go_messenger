// Code generated by gowrap. DO NOT EDIT.
// template: ../../../../core/pkg/decorators/templates/loggers.go
// gowrap: http://github.com/hexdigest/gowrap

package decorators

//go:generate gowrap gen -p github.com/YFatMR/go_messenger/user_service/internal/services -i UserService -t ../../../../core/pkg/decorators/templates/loggers.go -o loggers.gen.go -l ""

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
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
	return &LoggingUserServiceDecorator{
		base:   base,
		logger: logger,
	}
}

// Create implements services.UserService
func (d *LoggingUserServiceDecorator) Create(ctx context.Context, user *user.Entity, accountID *accountid.Entity) (userID *userid.Entity, cerr cerrors.Error) {

	d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: calling Create")
	defer func() {
		if cerr != nil {
			d.logger.ErrorContext(
				ctx, cerr.GetInternalErrorMessage(), zap.Error(cerr.GetInternalError()),
			)
		} else {
			d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: Create finished")
		}
		d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: Create finished")
	}()
	return d.base.Create(ctx, user, accountID)
}

// DeleteByID implements services.UserService
func (d *LoggingUserServiceDecorator) DeleteByID(ctx context.Context, userID *userid.Entity) (cerr cerrors.Error) {

	d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: calling DeleteByID")
	defer func() {
		if cerr != nil {
			d.logger.ErrorContext(
				ctx, cerr.GetInternalErrorMessage(), zap.Error(cerr.GetInternalError()),
			)
		} else {
			d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: DeleteByID finished")
		}
		d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: DeleteByID finished")
	}()
	return d.base.DeleteByID(ctx, userID)
}

// GetByID implements services.UserService
func (d *LoggingUserServiceDecorator) GetByID(ctx context.Context, userID *userid.Entity) (user *user.Entity, cerr cerrors.Error) {

	d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: calling GetByID")
	defer func() {
		if cerr != nil {
			d.logger.ErrorContext(
				ctx, cerr.GetInternalErrorMessage(), zap.Error(cerr.GetInternalError()),
			)
		} else {
			d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: GetByID finished")
		}
		d.logger.DebugContextNoExport(ctx, "LoggingUserServiceDecorator: GetByID finished")
	}()
	return d.base.GetByID(ctx, userID)
}
