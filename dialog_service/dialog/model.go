package dialog

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/dialog_service/entity"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type dialogModel struct {
	repository        apientity.DialogRepository
	userServiceClient proto.UserClient
	logger            *czap.Logger
}

func NewDialogModel(repository apientity.DialogRepository, userServiceClient proto.UserClient, logger *czap.Logger,
) apientity.DialogModel {
	return &dialogModel{
		repository:        repository,
		userServiceClient: userServiceClient,
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

func (m *dialogModel) GetDialogs(ctx context.Context, userID *entity.UserID, offset uint64, limit uint64) (
	[]*entity.Dialog, error,
) {
	return m.repository.GetDialogs(ctx, userID, offset, limit)
}

func (m *dialogModel) CreateDialogMessage(ctx context.Context, dialogID *entity.DialogID,
	message *entity.DialogMessage,
) (
	*entity.DialogMessage, error,
) {
	return m.repository.CreateDialogMessage(ctx, dialogID, message)
}

func (m *dialogModel) GetDialogMessages(ctx context.Context, dialogID *entity.DialogID,
	offset uint64, limit uint64,
) (
	[]*entity.DialogMessage, error,
) {
	return m.repository.GetDialogMessages(ctx, dialogID, offset, limit)
}
