package dialog

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/dialog_service/entity"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type dialogController struct {
	contextManager apientity.ContextManager
	model          apientity.DialogModel
	logger         *czap.Logger
}

func NewController(contextManager apientity.ContextManager, model apientity.DialogModel, logger *czap.Logger,
) apientity.DialogController {
	return &dialogController{
		contextManager: contextManager,
		model:          model,
		logger:         logger,
	}
}

func (c *dialogController) CreateDialogWith(ctx context.Context, request *proto.UserID) (
	*proto.Dialog, error,
) {
	userID1, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userID2, err := entity.UserIDFromProtobuf(request)
	if err != nil {
		return nil, err
	}

	dialog, err := c.model.CreateDialog(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}
	return entity.DialogToProtobuf(dialog), nil
}

func (c *dialogController) GetDialogs(ctx context.Context, request *proto.GetDialogsRequest) (
	*proto.GetDialogsResponse, error,
) {
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if request.GetOffset() == 0 || request.GetLimit() == 0 {
		return nil, ErrParseRequest
	}

	dialogs, err := c.model.GetDialogs(ctx, userID, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}

	return &proto.GetDialogsResponse{
		Dialogs: entity.DialogsToProtobuf(dialogs),
	}, nil
}

// func (c *dialogController) CreateDialogMessage(ctx context.Context, request *proto.CreateDialogMessageRequest) (
// 	*proto.CreateDialogMessageResponse, error,
// ) {
// 	message, err := entity.DialogMessageFromProtobuf(request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = c.model.CreateDialogMessage(ctx, message)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &proto.CreateDialogMessageResponse{
// 		CreationUnixTimestamp: message.CreationUnixTimestamp,
// 	}, nil
// }

// func (c *dialogController) GetDialogMessages(ctx context.Context, request *proto.GetDialogMessagesRequest) (
// 	*proto.GetDialogMessagesResponse, error,
// ) {
// 	userID1, err := entity.UserIDFromProtobuf(request.GetMemberID1())
// 	if err != nil {
// 		return nil, err
// 	}

// 	userID2, err := entity.UserIDFromProtobuf(request.GetMemberID2())
// 	if err != nil {
// 		return nil, err
// 	}

// 	messages, err := c.model.GetDialogMessages(ctx, userID1, userID2, request.Offset, request.Limit)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &proto.GetDialogMessagesResponse{
// 		Messages: entity.DialogMessagesToProtobuf(messages),
// 	}, nil
// }

func (c *dialogController) Ping(ctx context.Context, request *proto.Void) (
	pong *proto.Pong, err error,
) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
