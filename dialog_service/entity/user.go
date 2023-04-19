package entity

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type UserID struct {
	ID uint64
}

func UserIDFromProtobuf(userID *proto.UserID) (*UserID, error) {
	if userID == nil || userID.ID == 0 {
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

type UserData struct {
	Nickname string
	Name     string
	Surname  string
}

func UserDataFromProtobuf(user *proto.UserData) (*UserData, error) {
	if user == nil || user.Name == "" || user.Surname == "" || user.Nickname == "" {
		return nil, ErrWrongRequestFormat
	}
	return &UserData{
		Name:     user.Name,
		Surname:  user.Surname,
		Nickname: user.Nickname,
	}, nil
}
