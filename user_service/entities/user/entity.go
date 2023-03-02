package user

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/entities"
)

type Entity struct {
	nickname string
	name     string
	surname  string
}

func New(nickname string, name string, surname string) *Entity {
	return &Entity{
		name:     name,
		surname:  surname,
		nickname: nickname,
	}
}

func FromProtobuf(user *proto.UserData) (*Entity, error) {
	if user == nil || user.GetName() == "" || user.GetSurname() == "" || user.GetNickname() == "" {
		return nil, entities.ErrWrongRequestFormat
	}
	return New(user.GetName(), user.GetSurname(), user.GetNickname()), nil
}

func (e *Entity) GetNickname() string {
	return e.name
}

func (e *Entity) GetName() string {
	return e.name
}

func (e *Entity) GetSurname() string {
	return e.surname
}
