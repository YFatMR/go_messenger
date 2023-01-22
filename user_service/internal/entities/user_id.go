package entities

import proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type UserID struct {
	userID string
}

func newUserID(userID string) *UserID {
	return &UserID{
		userID: userID,
	}
}

func NewUserIDFromProtobuf(userID *proto.UserID) (*UserID, error) {
	if userID == nil || userID.GetID() == "" {
		return nil, ErrWrongRequestFormat
	}
	return newUserID(userID.GetID()), nil
}

func NewUserIDFromRawDatabaseDocument(userID string) *UserID {
	return newUserID(userID)
}

func (u *UserID) GetID() string {
	return u.userID
}
