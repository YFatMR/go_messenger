// Code generated by gowrap. DO NOT EDIT.
// template: ../../core/pkg/decorators/templates/czap_logger.template.go
// gowrap: http://github.com/hexdigest/gowrap

package decorator

//go:generate gowrap gen -p github.com/YFatMR/go_messenger/dialog_service/apientity -i DialogRepository -t ../../core/pkg/decorators/templates/czap_logger.template.go -o dialog_repository_czap_logger.gen.go -l ""

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/dialog_service/entity"
	"go.uber.org/zap"
)

// LoggingDialogRepositoryDecorator implements apientity.DialogRepository that is instrumented with custom zap logger
type LoggingDialogRepositoryDecorator struct {
	logger *czap.Logger
	base   apientity.DialogRepository
}

// NewLoggingDialogRepositoryDecorator instruments an implementation of the apientity.DialogRepository with simple logging
func NewLoggingDialogRepositoryDecorator(base apientity.DialogRepository, logger *czap.Logger) *LoggingDialogRepositoryDecorator {
	if base == nil {
		panic("LoggingDialogRepositoryDecorator got empty base")
	}
	if logger == nil {
		panic("LoggingDialogRepositoryDecorator got empty logger")
	}
	return &LoggingDialogRepositoryDecorator{
		base:   base,
		logger: logger,
	}
}

// CreateDialog implements apientity.DialogRepository
func (d *LoggingDialogRepositoryDecorator) CreateDialog(ctx context.Context, userID1 *entity.UserID, userData1 *entity.UserData, userID2 *entity.UserID, userData2 *entity.UserData) (dialog *entity.Dialog, err error) {

	d.logger.InfoContext(ctx, "LoggingDialogRepositoryDecorator: calling CreateDialog")
	defer func() {
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.InfoContext(ctx, "LoggingDialogRepositoryDecorator: CreateDialog finished")
	}()
	return d.base.CreateDialog(ctx, userID1, userData1, userID2, userData2)
}

// CreateDialogMessage implements apientity.DialogRepository
func (d *LoggingDialogRepositoryDecorator) CreateDialogMessage(ctx context.Context, dialogID *entity.DialogID, message *entity.DialogMessage) (msg *entity.DialogMessage, err error) {

	d.logger.InfoContext(ctx, "LoggingDialogRepositoryDecorator: calling CreateDialogMessage")
	defer func() {
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.InfoContext(ctx, "LoggingDialogRepositoryDecorator: CreateDialogMessage finished")
	}()
	return d.base.CreateDialogMessage(ctx, dialogID, message)
}

// GetDialogMessages implements apientity.DialogRepository
func (d *LoggingDialogRepositoryDecorator) GetDialogMessages(ctx context.Context, dialogID *entity.DialogID, offset uint64, limit uint64) (messages []*entity.DialogMessage, err error) {

	d.logger.InfoContext(ctx, "LoggingDialogRepositoryDecorator: calling GetDialogMessages")
	defer func() {
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.InfoContext(ctx, "LoggingDialogRepositoryDecorator: GetDialogMessages finished")
	}()
	return d.base.GetDialogMessages(ctx, dialogID, offset, limit)
}

// GetDialogs implements apientity.DialogRepository
func (d *LoggingDialogRepositoryDecorator) GetDialogs(ctx context.Context, userID *entity.UserID, offset uint64, limit uint64) (dialogs []*entity.Dialog, err error) {

	d.logger.InfoContext(ctx, "LoggingDialogRepositoryDecorator: calling GetDialogs")
	defer func() {
		if err != nil {
			d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
		}
		d.logger.InfoContext(ctx, "LoggingDialogRepositoryDecorator: GetDialogs finished")
	}()
	return d.base.GetDialogs(ctx, userID, offset, limit)
}
