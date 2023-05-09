package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/dialog_service/entity"
)

type DialogModel interface {
	CreateDialog(ctx context.Context, userID1 *entity.UserID, userID2 *entity.UserID) (
		dialog *entity.Dialog, err error,
	)
	GetDialog(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID) (
		dialog *entity.Dialog, err error,
	)
	GetDialogs(ctx context.Context, userID *entity.UserID, offset uint64, limit uint64) (
		dialogs []*entity.Dialog, err error,
	)
	CreateDialogMessage(ctx context.Context, dialogID *entity.DialogID, message *entity.DialogMessage) (
		msg *entity.DialogMessage, err error,
	)
	GetDialogMessages(ctx context.Context, dialogID *entity.DialogID, messageID *entity.MessageID,
		limit uint64, offsetType entity.OffserType,
	) (
		messages []*entity.DialogMessage, err error,
	)
	ReadAllMessagesBeforeAndIncl(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID,
		messageID *entity.MessageID,
	) error
}
