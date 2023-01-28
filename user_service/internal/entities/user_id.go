package entities

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type UserID struct {
	userID string
}

func NewUserID(userID string) *UserID {
	return &UserID{
		userID: userID,
	}
}

func NewUserIDFromProtobuf(userID *proto.UserID) (*UserID, error) {
	if userID == nil || userID.GetID() == "" {
		return nil, ErrWrongRequestFormat
	}
	return NewUserID(userID.GetID()), nil
}

func (u *UserID) GetID() string {
	return u.userID
}
