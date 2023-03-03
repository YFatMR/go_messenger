package entity

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type UserID struct {
	ID string
}

func UserIDFromProtobuf(userID *proto.UserID) (*UserID, error) {
	if userID == nil || userID.ID == "" {
		return nil, ErrWrongRequestFormat
	}
	return &UserID{
		ID: userID.ID,
	}, nil
}

func UserIDToProtobuf(userID *UserID) *proto.UserID {
	return &proto.UserID{
		ID: userID.ID,
	}
}
