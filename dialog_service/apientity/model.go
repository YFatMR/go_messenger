package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/dialog_service/entity"
)

type DialogModel interface {
	CreateDialog(ctx context.Context, userID1 *entity.UserID, userID2 *entity.UserID) (
		dialog *entity.Dialog, err error,
	)
	GetDialogs(ctx context.Context, userID *entity.UserID, offset uint64, limit uint64) (
		dialogs []*entity.Dialog, err error,
	)
	// CreateDialogMessage(ctx context.Context, message *entity.DialogMessage) (
	// 	err error,
	// )
	// GetDialogMessages(ctx context.Context, userID1 *entity.UserID, userID2 *entity.UserID,
	// 	offset int64, limit int64,
	// ) (
	// 	messages []*entity.DialogMessage, err error,
	// )
}
