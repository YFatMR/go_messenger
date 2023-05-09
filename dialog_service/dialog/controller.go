package dialog

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/dialog_service/entity"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	return entity.DialogToProtobuf(dialog, userID1), nil
}

func (c *dialogController) GetDialogByID(ctx context.Context, request *proto.DialogID) (
	*proto.Dialog, error,
) {
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dialogID, err := entity.DialogIDFromProtobuf(request)
	if err != nil {
		return nil, err
	}

	dialog, err := c.model.GetDialog(ctx, userID, dialogID)
	if err != nil {
		return nil, err
	}
	return entity.DialogToProtobuf(dialog, userID), nil
}

func (c *dialogController) GetDialogs(ctx context.Context, request *proto.GetDialogsRequest) (
	*proto.GetDialogsResponse, error,
) {
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if request.GetLimit() == 0 {
		return nil, ErrParseRequest
	}

	dialogs, err := c.model.GetDialogs(ctx, userID, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}
	return &proto.GetDialogsResponse{
		Dialogs: entity.DialogsToProtobuf(dialogs, userID),
	}, nil
}

func (c *dialogController) CreateDialogMessage(ctx context.Context, request *proto.CreateDialogMessageRequest) (
	*proto.CreateDialogMessageResponse, error,
) {
	senderID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	message, err := entity.DialogMessageFromProtobuf(request)
	if err != nil {
		return nil, err
	}
	message.SenderID = *senderID

	dialogID, err := entity.DialogIDFromProtobuf(request.GetDialogID())
	if err != nil {
		return nil, err
	}

	message, err = c.model.CreateDialogMessage(ctx, dialogID, message)
	if err != nil {
		return nil, err
	}
	return &proto.CreateDialogMessageResponse{
		CreatedAt: timestamppb.New(message.CreatedAt),
		MessageID: &proto.MessageID{
			ID: message.MessageID.ID,
		},
	}, nil
}

func (c *dialogController) GetDialogMessages(ctx context.Context, request *proto.GetDialogMessagesRequest) (
	*proto.GetDialogMessagesResponse, error,
) {
	senderID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dialogID, err := entity.DialogIDFromProtobuf(request.GetDialogID())
	if err != nil {
		return nil, err
	}

	messageID, err := entity.MessageIDFromProtobuf(request.GetMessageID())
	if err != nil {
		return nil, err
	}

	offsetType := entity.OffserTypeFromProtobuf(request.GetOffsetType())

	if request.GetLimit() == 0 {
		return nil, ErrWrongRequestFormat
	}

	messages, err := c.model.GetDialogMessages(
		ctx, dialogID, messageID, request.Limit, offsetType,
	)
	if err != nil {
		return nil, err
	}
	return &proto.GetDialogMessagesResponse{
		Messages: entity.DialogMessagesToProtobuf(messages, senderID),
	}, nil
}

func (c *dialogController) ReadAllMessagesBeforeAndIncl(ctx context.Context, request *proto.ReadAllMessagesBeforeRequest) (
	*proto.Void, error,
) {
	senderID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dialogID, err := entity.DialogIDFromProtobuf(request.GetDialogID())
	if err != nil {
		return nil, err
	}

	messageID, err := entity.MessageIDFromProtobuf(request.GetMessageID())
	if err != nil {
		return nil, err
	}

	err = c.model.ReadAllMessagesBeforeAndIncl(ctx, senderID, dialogID, messageID)
	if err != nil {
		return nil, err
	}
	return &proto.Void{}, nil
}

func (c *dialogController) Ping(ctx context.Context, request *proto.Void) (
	pong *proto.Pong, err error,
) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
