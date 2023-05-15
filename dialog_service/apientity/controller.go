package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type DialogController interface {
	CreateDialogWith(ctx context.Context, request *proto.UserID) (
		response *proto.Dialog, err error,
	)
	GetDialogByID(ctx context.Context, dialogID *proto.DialogID) (
		response *proto.Dialog, err error,
	)
	GetDialogs(ctx context.Context, request *proto.GetDialogsRequest) (
		response *proto.GetDialogsResponse, err error,
	)
	CreateDialogMessageWithCode(ctx context.Context, request *proto.CreateDialogMessageWithCodeRequest) (
		*proto.CreateDialogMessageResponse, error,
	)
	CreateDialogMessage(ctx context.Context, request *proto.CreateDialogMessageRequest) (
		response *proto.CreateDialogMessageResponse, err error,
	)
	GetDialogMessages(ctx context.Context, request *proto.GetDialogMessagesRequest) (
		response *proto.GetDialogMessagesResponse, err error,
	)
	ReadAllMessagesBeforeAndInclude(ctx context.Context, request *proto.ReadAllMessagesBeforeRequest) (
		void *proto.Void, err error,
	)
	CreateInstruction(ctx context.Context, request *proto.CreateInstructionRequest) (
		response *proto.InstructionID, err error,
	)
	GetInstructions(ctx context.Context, request *proto.GetInstructionsRequest) (
		response *proto.GetInstructionsResponse, err error,
	)
	GetInstructionsByID(ctx context.Context, request *proto.GetInstructionsByIDRequest) (
		response *proto.GetInstructionsResponse, err error,
	)
	Ping(ctx context.Context, request *proto.Void) (
		pong *proto.Pong, err error,
	)
}
