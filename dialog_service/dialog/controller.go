package dialog

import (
	"context"
	"fmt"

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

	newRequest, err := entity.CreateDialogMessageRequestFromProtobuf(request)
	if err != nil {
		return nil, err
	}
	newRequest.SenderID = *senderID

	message, err := c.model.CreateDialogMessage(ctx, newRequest)
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

func (c *dialogController) CreateDialogMessageWithCode(ctx context.Context,
	request *proto.CreateDialogMessageWithCodeRequest) (
	*proto.CreateDialogMessageResponse, error,
) {
	senderID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	newRequest, err := entity.CreateDialogMessageWithCodeRequestFromProtobuf(request)
	if err != nil {
		return nil, err
	}
	newRequest.SenderID = *senderID

	message, err := c.model.CreateDialogMessageWithCode(ctx, newRequest)
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

	offsetType := entity.DialogMessagesOffserTypeFromProtobuf(request.GetOffsetType())

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

func (c *dialogController) ReadMessage(ctx context.Context, request *proto.ReadMessageRequest) (
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

	err = c.model.ReadMessage(ctx, senderID, dialogID, messageID)
	if err != nil {
		return nil, err
	}
	return &proto.Void{}, nil
}

func (c *dialogController) CreateInstruction(ctx context.Context, request *proto.CreateInstructionRequest) (
	*proto.InstructionID, error,
) {
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dialogID, err := entity.DialogIDFromProtobuf(request.GetDialogID())
	if err != nil {
		return nil, err
	}

	if request.GetTitle() == "" {
		return nil, fmt.Errorf("empty title")
	}

	if request.GetText() == "" {
		return nil, fmt.Errorf("empty text")
	}

	instructionID, err := c.model.CreateInstruction(ctx, userID, dialogID, request.Title, request.Text)
	if err != nil {
		return nil, err
	}
	return entity.InstructionIDToProtobuf(instructionID), nil
}

func (c *dialogController) GetInstructions(ctx context.Context, request *proto.GetInstructionsRequest) (
	*proto.GetInstructionsResponse, error,
) {
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dialogID, err := entity.DialogIDFromProtobuf(request.GetDialogID())
	if err != nil {
		return nil, err
	}

	if request.GetLimit() <= 0 {
		return nil, fmt.Errorf("incorrect limit params")
	}

	instructions, err := c.model.GetInstructions(ctx, userID, dialogID, request.Limit)
	if err != nil {
		return nil, err
	}
	return &proto.GetInstructionsResponse{
		Instructions: entity.InstructionsToProtobuf(instructions),
	}, nil
}

func (c *dialogController) GetInstructionsByID(ctx context.Context, request *proto.GetInstructionsByIDRequest) (
	*proto.GetInstructionsResponse, error,
) {
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dialogID, err := entity.DialogIDFromProtobuf(request.GetDialogID())
	if err != nil {
		return nil, err
	}

	instructionID, err := entity.InstructionIDFromProtobuf(request.GetInstructionID())
	if err != nil {
		return nil, err
	}
	offsetType := entity.InstructionOffserTypeFromProtobuf(request.GetOffsetType())
	if request.GetLimit() <= 0 {
		return nil, fmt.Errorf("incorrect limit params")
	}

	instructions, err := c.model.GetInstructionsByID(ctx, userID, dialogID, instructionID, offsetType, request.Limit)
	if err != nil {
		return nil, err
	}
	return &proto.GetInstructionsResponse{
		Instructions: entity.InstructionsToProtobuf(instructions),
	}, nil
}

func (c *dialogController) GetDialogLinks(ctx context.Context, request *proto.GetDialogLinksRequest) (
	*proto.GetDialogLinksResponse, error,
) {
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dialogID, err := entity.DialogIDFromProtobuf(request.GetDialogID())
	if err != nil {
		return nil, err
	}

	if request.GetLimit() <= 0 {
		return nil, fmt.Errorf("incorrect limit params")
	}

	links, err := c.model.GetLinks(ctx, userID, dialogID, request.Limit)
	if err != nil {
		return nil, err
	}
	return &proto.GetDialogLinksResponse{
		Links: entity.LinksToProtobuf(links),
	}, nil
}

func (c *dialogController) GetDialogLinksByID(ctx context.Context, request *proto.GetDialogLinksByIDRequest) (
	*proto.GetDialogLinksResponse, error,
) {
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dialogID, err := entity.DialogIDFromProtobuf(request.GetDialogID())
	if err != nil {
		return nil, err
	}

	linkID, err := entity.LinkIDFromProtobuf(request.GetLinkID())
	if err != nil {
		return nil, err
	}

	offsetType := entity.LinkOffserTypeFromProtobuf(request.OffsetType)

	if request.GetLimit() <= 0 {
		return nil, fmt.Errorf("incorrect limit params")
	}

	links, err := c.model.GetLinksByID(ctx, userID, dialogID, linkID, offsetType, request.Limit)
	if err != nil {
		return nil, err
	}
	return &proto.GetDialogLinksResponse{
		Links: entity.LinksToProtobuf(links),
	}, nil
}

func (c *dialogController) GetDialogMembers(ctx context.Context, request *proto.DialogID) (
	*proto.GetDialogMembersResponse, error,
) {
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dialogID, err := entity.DialogIDFromProtobuf(request)
	if err != nil {
		return nil, err
	}

	selfID, memberID, err := c.model.GetDialogMembers(ctx, userID, dialogID)
	if err != nil {
		return nil, err
	}

	return &proto.GetDialogMembersResponse{
		SelfID: &proto.UserID{
			ID: selfID.ID,
		},
		MemberID: &proto.UserID{
			ID: memberID.ID,
		},
	}, nil
}

func (c *dialogController) GetUnreadDialogMessagesCount(ctx context.Context, request *proto.DialogID) (
	*proto.GetUnreadDialogMessagesCountResponse, error,
) {
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dialogID, err := entity.DialogIDFromProtobuf(request)
	if err != nil {
		return nil, err
	}

	unreadMessagesCount, err := c.model.GetUnreadDialogMessagesCount(ctx, userID, dialogID)
	if err != nil {
		return nil, err
	}

	return &proto.GetUnreadDialogMessagesCountResponse{
		Count: unreadMessagesCount,
	}, nil
}

func (c *dialogController) Ping(ctx context.Context, request *proto.Void) (
	pong *proto.Pong, err error,
) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
