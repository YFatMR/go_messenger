package entity

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type User struct {
	Nickname string
	Name     string
	Surname  string
}

func UserFromProtobuf(user *proto.UserData) (*User, error) {
	if user == nil || user.Name == "" || user.Surname == "" || user.Nickname == "" {
		return nil, ErrWrongRequestFormat
	}
	return &User{
		Name:     user.Name,
		Surname:  user.Surname,
		Nickname: user.Nickname,
	}, nil
}

func UserToProtobuf(user *User) *proto.UserData {
	return &proto.UserData{
		Name:    user.Name,
		Surname: user.Surname,
	}
}
