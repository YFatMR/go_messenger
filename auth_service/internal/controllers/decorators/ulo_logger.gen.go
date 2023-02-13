// Code generated by gowrap. DO NOT EDIT.
// template: ../../../../core/pkg/decorators/templates/ulo_logger.go
// gowrap: http://github.com/hexdigest/gowrap

package decorators

//go:generate gowrap gen -p github.com/YFatMR/go_messenger/auth_service/internal/controllers -i AccountController -t ../../../../core/pkg/decorators/templates/ulo_logger.go -o ulo_logger.gen.go -l ""

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/controllers"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.uber.org/zap"
)

// LoggingAccountControllerDecorator implements controllers.AccountController that is instrumented with custom zap logger
type LoggingAccountControllerDecorator struct {
	logger *loggers.OtelZapLoggerWithTraceID
	base   controllers.AccountController
}

// NewLoggingAccountControllerDecorator instruments an implementation of the controllers.AccountController with simple logging
func NewLoggingAccountControllerDecorator(base controllers.AccountController, logger *loggers.OtelZapLoggerWithTraceID) *LoggingAccountControllerDecorator {
	if base == nil {
		panic("LoggingAccountControllerDecorator got empty base")
	}
	if logger == nil {
		panic("LoggingAccountControllerDecorator got empty logger")
	}
	return &LoggingAccountControllerDecorator{
		base:   base,
		logger: logger,
	}
}

// CreateAccount implements controllers.AccountController
func (d *LoggingAccountControllerDecorator) CreateAccount(ctx context.Context, request *proto.Credential) (accountID *proto.AccountID, logstash ulo.LogStash, err error) {
	d.logger.DebugContextNoExport(ctx, "LoggingAccountControllerDecorator: calling CreateAccount")
	defer func() {
		d.logger.LogContextNoExportULO(ctx, logstash)
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.DebugContextNoExport(ctx, "LoggingAccountControllerDecorator: CreateAccount finished")
	}()
	return d.base.CreateAccount(ctx, request)
}

// GetToken implements controllers.AccountController
func (d *LoggingAccountControllerDecorator) GetToken(ctx context.Context, request *proto.Credential) (token *proto.Token, logstash ulo.LogStash, err error) {
	d.logger.DebugContextNoExport(ctx, "LoggingAccountControllerDecorator: calling GetToken")
	defer func() {
		d.logger.LogContextNoExportULO(ctx, logstash)
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.DebugContextNoExport(ctx, "LoggingAccountControllerDecorator: GetToken finished")
	}()
	return d.base.GetToken(ctx, request)
}

// GetTokenPayload implements controllers.AccountController
func (d *LoggingAccountControllerDecorator) GetTokenPayload(ctx context.Context, request *proto.Token) (tokenPayload *proto.TokenPayload, logstash ulo.LogStash, err error) {
	d.logger.DebugContextNoExport(ctx, "LoggingAccountControllerDecorator: calling GetTokenPayload")
	defer func() {
		d.logger.LogContextNoExportULO(ctx, logstash)
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.DebugContextNoExport(ctx, "LoggingAccountControllerDecorator: GetTokenPayload finished")
	}()
	return d.base.GetTokenPayload(ctx, request)
}

// Ping implements controllers.AccountController
func (d *LoggingAccountControllerDecorator) Ping(ctx context.Context, request *proto.Void) (pong *proto.Pong, logstash ulo.LogStash, err error) {
	d.logger.DebugContextNoExport(ctx, "LoggingAccountControllerDecorator: calling Ping")
	defer func() {
		d.logger.LogContextNoExportULO(ctx, logstash)
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.DebugContextNoExport(ctx, "LoggingAccountControllerDecorator: Ping finished")
	}()
	return d.base.Ping(ctx, request)
}