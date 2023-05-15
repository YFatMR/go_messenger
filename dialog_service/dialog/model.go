package dialog

import (
	"context"
	"regexp"

	"github.com/YFatMR/go_messenger/core/pkg/ckafka"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/dialog_service/entity"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.uber.org/zap"
)

type dialogModel struct {
	repository        apientity.DialogRepository
	userServiceClient proto.UserClient
	kafkaClient       apientity.KafkaClient
	logger            *czap.Logger
}

func NewDialogModel(repository apientity.DialogRepository, userServiceClient proto.UserClient,
	kafkaClient apientity.KafkaClient, logger *czap.Logger,
) apientity.DialogModel {
	return &dialogModel{
		repository:        repository,
		userServiceClient: userServiceClient,
		kafkaClient:       kafkaClient,
		logger:            logger,
	}
}

func (m *dialogModel) isUserDialogMember(ctx context.Context, dialogID *entity.DialogID, userID *entity.UserID) (
	bool, error,
) {
	members, err := m.repository.GetDialogMembers(ctx, dialogID)
	if err != nil {
		return false, err
	}
	for _, member := range members {
		if member.ID == userID.ID {
			return true, nil
		}
	}
	return false, nil
}

func (m *dialogModel) CreateDialog(ctx context.Context, userID1 *entity.UserID, userID2 *entity.UserID) (
	*entity.Dialog, error,
) {
	getUserData := func(userID *entity.UserID) (*entity.UserData, error) {
		userDataPb, err := m.userServiceClient.GetUserByID(ctx, entity.UserIDToProtobuf(userID))
		if err != nil {
			return nil, err
		}
		return entity.UserDataFromProtobuf(userDataPb)
	}

	userData1, err := getUserData(userID1)
	if err != nil {
		return nil, err
	}
	userData2, err := getUserData(userID2)
	if err != nil {
		return nil, err
	}

	return m.repository.CreateDialog(ctx, userID1, userData1, userID2, userData2)
}

func (m *dialogModel) GetDialog(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID) (
	*entity.Dialog, error,
) {
	return m.repository.GetDialog(ctx, userID, dialogID)
}

func (m *dialogModel) GetDialogs(ctx context.Context, userID *entity.UserID, offset uint64, limit uint64) (
	[]*entity.Dialog, error,
) {
	return m.repository.GetDialogs(ctx, userID, offset, limit)
}

func (m *dialogModel) CreateDialogMessage(ctx context.Context, request *entity.CreateDialogMessageRequest) (
	*entity.DialogMessage, error,
) {
	members, err := m.repository.GetDialogMembers(ctx, &request.DialogID)
	if err != nil {
		return nil, err
	}
	if members[0].ID != request.SenderID.ID && members[1].ID != request.SenderID.ID {
		return nil, ErrFobidden
	}

	messageURLs := func(messageText string) []string {
		re := regexp.MustCompile(`https?://[^\s]+`)
		return re.FindAllString(messageText, -1)
	}(request.Text)
	message, err := m.repository.CreateDialogMessageWithURLs(ctx, request, messageURLs)
	if err != nil {
		return nil, err
	}

	// async writing
	go m.kafkaClient.WriteNewDialogMessage(context.TODO(), &ckafka.DialogMessage{
		MessageID: ckafka.MessageID{
			ID: message.MessageID.ID,
		},
		SenderID: ckafka.UserID{
			ID: message.SenderID.ID,
		},
		ReciverID: ckafka.UserID{
			ID: func() uint64 {
				if members[0].ID == message.SenderID.ID {
					return members[1].ID
				}
				return members[0].ID
			}(),
		},
		DialogID: ckafka.DialogID{
			ID: request.DialogID.ID,
		},
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
		Type:      message.Type,
	})
	return message, nil
}

func (m *dialogModel) CreateDialogMessageWithCode(ctx context.Context,
	request *entity.CreateDialogMessageWithCodeRequest) (
	*entity.DialogMessage, error,
) {
	members, err := m.repository.GetDialogMembers(ctx, &request.DialogID)
	if err != nil {
		return nil, err
	}
	if members[0].ID != request.SenderID.ID && members[1].ID != request.SenderID.ID {
		return nil, ErrFobidden
	}

	message, err := m.repository.CreateDialogMessageWithCode(ctx, request)
	if err != nil {
		return nil, err
	}

	// async writing
	go m.kafkaClient.WriteNewDialogMessage(context.TODO(), &ckafka.DialogMessage{
		MessageID: ckafka.MessageID{
			ID: message.MessageID.ID,
		},
		SenderID: ckafka.UserID{
			ID: message.SenderID.ID,
		},
		ReciverID: ckafka.UserID{
			ID: func() uint64 {
				if members[0].ID == message.SenderID.ID {
					return members[1].ID
				}
				return members[0].ID
			}(),
		},
		DialogID: ckafka.DialogID{
			ID: request.DialogID.ID,
		},
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
		Type:      message.Type,
	})
	return message, nil
}

func (m *dialogModel) GetDialogMessages(ctx context.Context, dialogID *entity.DialogID,
	messageID *entity.MessageID, limit uint64, offsetType entity.DialogMessagesOffserType,
) (
	[]*entity.DialogMessage, error,
) {
	switch offsetType {
	case entity.DIALOG_MESSAGE_OFFSET_BEFORE:
		return m.repository.GetDialogMessagesBefore(ctx, dialogID, messageID, limit)
	case entity.DIALOG_MESSAGE_OFFSET_BEFORE_INCLUDE:
		return m.repository.GetDialogMessagesBeforeAndInclude(ctx, dialogID, messageID, limit)
	case entity.DIALOG_MESSAGE_OFFSET_AFTER:
		return m.repository.GetDialogMessagesAfter(ctx, dialogID, messageID, limit)
	default:
		return m.repository.GetDialogMessagesAfterAndInclude(ctx, dialogID, messageID, limit)
	}
}

func (m *dialogModel) ReadAllMessagesBeforeAndIncl(ctx context.Context, userID *entity.UserID,
	dialogID *entity.DialogID, messageID *entity.MessageID,
) error {
	if err := m.repository.ReadAllMessagesBeforeAndIncl(ctx, userID, dialogID, messageID); err != nil {
		return err
	}

	// async writing
	go func() {
		message, err := m.repository.GetDialogMessageByID(context.TODO(), dialogID, messageID)
		if err != nil {
			m.logger.Error("ReadAllMessagesBeforeAndIncl async op failed", zap.Error(err))
			return
		}

		m.kafkaClient.WriteNewViewedMessage(context.TODO(), &ckafka.ViewedMessage{
			MessageID: ckafka.MessageID{
				ID: messageID.ID,
			},
			MessageCreatedAt: message.CreatedAt,
			SenderID: ckafka.UserID{
				ID: userID.ID,
			},
			ReciverID: ckafka.UserID{
				ID: message.SenderID.ID,
			},
			DialogID: ckafka.DialogID{
				ID: dialogID.ID,
			},
		})
	}()
	return nil
}

func (m *dialogModel) CreateInstruction(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID,
	instructionTitle string, instructionText string,
) (
	*entity.InstructionID, error,
) {
	isMember, err := m.isUserDialogMember(ctx, dialogID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrFobidden
	}
	return m.repository.CreateInstruction(ctx, userID, dialogID, instructionTitle, instructionText)
}

func (m *dialogModel) GetInstructions(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID,
	limit uint64,
) (
	instructions []*entity.Instruction, err error,
) {
	isMember, err := m.isUserDialogMember(ctx, dialogID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrFobidden
	}
	return m.repository.GetInstructions(ctx, dialogID, limit)
}

func (m *dialogModel) GetInstructionsByID(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID,
	instructionID *entity.InstructionID, offsetType entity.InstructionOffserType, limit uint64,
) (
	instructions []*entity.Instruction, err error,
) {
	isMember, err := m.isUserDialogMember(ctx, dialogID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrFobidden
	}
	return m.repository.GetInstructionsBefore(ctx, dialogID, instructionID, limit)
}
