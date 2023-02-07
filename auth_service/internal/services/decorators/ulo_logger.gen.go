// Code generated by gowrap. DO NOT EDIT.
// template: ../../../../core/pkg/decorators/templates/ulo_logger.go
// gowrap: http://github.com/hexdigest/gowrap

package decorators

//go:generate gowrap gen -p github.com/YFatMR/go_messenger/auth_service/internal/services -i AccountService -t ../../../../core/pkg/decorators/templates/ulo_logger.go -o ulo_logger.gen.go -l ""

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/token"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"
	"github.com/YFatMR/go_messenger/auth_service/internal/services"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"go.uber.org/zap"
)

// LoggingAccountServiceDecorator implements services.AccountService that is instrumented with custom zap logger
type LoggingAccountServiceDecorator struct {
	logger *loggers.OtelZapLoggerWithTraceID
	base   services.AccountService
}

// NewLoggingAccountServiceDecorator instruments an implementation of the services.AccountService with simple logging
func NewLoggingAccountServiceDecorator(base services.AccountService, logger *loggers.OtelZapLoggerWithTraceID) *LoggingAccountServiceDecorator {
	if base == nil {
		panic("LoggingAccountServiceDecorator got empty base")
	}
	if logger == nil {
		panic("LoggingAccountServiceDecorator got empty logger")
	}
	return &LoggingAccountServiceDecorator{
		base:   base,
		logger: logger,
	}
}

// CreateAccount implements services.AccountService
func (d *LoggingAccountServiceDecorator) CreateAccount(ctx context.Context, credential *credential.Entity) (accountID *accountid.Entity, logStash ulo.LogStash, err error) {
	d.logger.DebugContextNoExport(ctx, "LoggingAccountServiceDecorator: calling CreateAccount")
	defer func() {
		d.logger.LogContextNoExportULO(ctx, logStash)
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.DebugContextNoExport(ctx, "LoggingAccountServiceDecorator: CreateAccount finished")
	}()
	return d.base.CreateAccount(ctx, credential)
}

// GetToken implements services.AccountService
func (d *LoggingAccountServiceDecorator) GetToken(ctx context.Context, credential *credential.Entity) (token *token.Entity, logStash ulo.LogStash, err error) {
	d.logger.DebugContextNoExport(ctx, "LoggingAccountServiceDecorator: calling GetToken")
	defer func() {
		d.logger.LogContextNoExportULO(ctx, logStash)
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.DebugContextNoExport(ctx, "LoggingAccountServiceDecorator: GetToken finished")
	}()
	return d.base.GetToken(ctx, credential)
}

// GetTokenPayload implements services.AccountService
func (d *LoggingAccountServiceDecorator) GetTokenPayload(ctx context.Context, token *token.Entity) (tokenPayload *tokenpayload.Entity, logStash ulo.LogStash, err error) {
	d.logger.DebugContextNoExport(ctx, "LoggingAccountServiceDecorator: calling GetTokenPayload")
	defer func() {
		d.logger.LogContextNoExportULO(ctx, logStash)
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.DebugContextNoExport(ctx, "LoggingAccountServiceDecorator: GetTokenPayload finished")
	}()
	return d.base.GetTokenPayload(ctx, token)
}
