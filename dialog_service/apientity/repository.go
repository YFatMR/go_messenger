package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/dialog_service/entity"
)

type DialogRepository interface {
	CreateDialog(ctx context.Context, userID1 *entity.UserID, userData1 *entity.UserData, userID2 *entity.UserID,
		userData2 *entity.UserData,
	) (
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
	GetDialogMessagesBefore(ctx context.Context, dialogID *entity.DialogID, messageID *entity.MessageID,
		limit uint64,
	) (
		messages []*entity.DialogMessage, err error,
	)
	GetDialogMessagesBeforeAndInclude(ctx context.Context, dialogID *entity.DialogID, messageID *entity.MessageID,
		limit uint64,
	) (
		messages []*entity.DialogMessage, err error,
	)
	GetDialogMessagesAfter(ctx context.Context, dialogID *entity.DialogID, messageID *entity.MessageID,
		limit uint64,
	) (
		messages []*entity.DialogMessage, err error,
	)
	GetDialogMessagesAfterAndInclude(ctx context.Context, dialogID *entity.DialogID, messageID *entity.MessageID,
		limit uint64,
	) (
		messages []*entity.DialogMessage, err error,
	)
	GetDialogMessageByID(ctx context.Context, dialogID *entity.DialogID,
		messageID *entity.MessageID,
	) (
		message *entity.DialogMessage, err error,
	)
	GetDialogMembers(ctx context.Context, dialogID *entity.DialogID) (
		userIDs []*entity.UserID, err error,
	)
	ReadAllMessagesBeforeAndIncl(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID,
		messageID *entity.MessageID,
	) error
}
