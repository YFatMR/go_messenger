package accountid

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
)

type Entity struct {
	accountID string
}

func New(accountID string) *Entity {
	return &Entity{
		accountID: accountID,
	}
}

func FromProtobuf(accountID *proto.AccountID) (*Entity, error) {
	if accountID == nil || accountID.GetID() == "" {
		return nil, entities.ErrWrongRequestFormat
	}
	return New(accountID.GetID()), nil
}

func (e *Entity) GetID() string {
	return e.accountID
}
