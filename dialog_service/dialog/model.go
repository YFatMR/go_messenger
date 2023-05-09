package dialog

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/ckafka"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/dialog_service/entity"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
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

func (m *dialogModel) CreateDialogMessage(ctx context.Context, dialogID *entity.DialogID,
	inMessage *entity.DialogMessage,
) (
	*entity.DialogMessage, error,
) {
	members, err := m.repository.GetDialogMembers(ctx, dialogID)
	if err != nil {
		return nil, err
	}
	if members[0].ID != inMessage.SenderID.ID && members[1].ID != inMessage.SenderID.ID {
		return nil, ErrFobidden
	}

	message, err := m.repository.CreateDialogMessage(ctx, dialogID, inMessage)
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
			ID: dialogID.ID,
		},
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
	})
	return message, nil
}

func (m *dialogModel) GetDialogMessages(ctx context.Context, dialogID *entity.DialogID,
	messageID *entity.MessageID, limit uint64, offsetType entity.OffserType,
) (
	[]*entity.DialogMessage, error,
) {
	return m.repository.GetDialogMessages(ctx, dialogID, messageID, limit, offsetType)
}

func (m *dialogModel) ReadAllMessagesBeforeAndIncl(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID,
	messageID *entity.MessageID,
) error {
	return m.repository.ReadAllMessagesBeforeAndIncl(ctx, userID, dialogID, messageID)
}
