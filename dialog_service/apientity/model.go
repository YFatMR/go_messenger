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
	CreateDialogMessage(ctx context.Context, request *entity.CreateDialogMessageRequest) (
		msg *entity.DialogMessage, err error,
	)
	CreateDialogMessageWithCode(ctx context.Context, request *entity.CreateDialogMessageWithCodeRequest) (
		*entity.DialogMessage, error,
	)
	GetDialogMessages(ctx context.Context, dialogID *entity.DialogID, messageID *entity.MessageID,
		limit uint64, offsetType entity.DialogMessagesOffserType,
	) (
		messages []*entity.DialogMessage, err error,
	)
	ReadMessage(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID,
		messageID *entity.MessageID,
	) error
	CreateInstruction(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID, instructionTitle string,
		instructionText string,
	) (
		instructionID *entity.InstructionID, err error,
	)
	GetInstructions(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID, limit uint64) (
		instructions []*entity.Instruction, err error,
	)
	GetInstructionsByID(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID,
		instructionID *entity.InstructionID, offsetType entity.InstructionOffserType, limit uint64,
	) (
		instructions []*entity.Instruction, err error,
	)
	GetLinks(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID, limit uint64) (
		links []*entity.Link, err error,
	)
	GetLinksByID(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID,
		linkID *entity.LinkID, offsetType entity.LinkOffserType, limit uint64,
	) (
		links []*entity.Link, err error,
	)
	GetDialogMembers(ctx context.Context, selfID *entity.UserID, dialogID *entity.DialogID) (
		_selfID *entity.UserID, _memberID *entity.UserID, err error,
	)
	GetUnreadDialogMessagesCount(ctx context.Context, selfID *entity.UserID, dialogID *entity.DialogID) (
		count uint64, err error,
	)
}
