package userid

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/entities"
)

type Entity struct {
	userID string
}

func New(userID string) *Entity {
	return &Entity{
		userID: userID,
	}
}

func FromProtobuf(userID *proto.UserID) (*Entity, error) {
	if userID == nil || userID.GetID() == "" {
		return nil, entities.ErrWrongRequestFormat
	}
	return New(userID.GetID()), nil
}

func (e *Entity) GetID() string {
	return e.userID
}
