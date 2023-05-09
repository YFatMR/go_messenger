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
	CreateDialogMessage(ctx context.Context, request *proto.CreateDialogMessageRequest) (
		response *proto.CreateDialogMessageResponse, err error,
	)
	GetDialogMessages(ctx context.Context, request *proto.GetDialogMessagesRequest) (
		response *proto.GetDialogMessagesResponse, err error,
	)
	ReadAllMessagesBeforeAndIncl(ctx context.Context, request *proto.ReadAllMessagesBeforeRequest) (
		void *proto.Void, err error,
	)
	Ping(ctx context.Context, request *proto.Void) (
		pong *proto.Pong, err error,
	)
}
