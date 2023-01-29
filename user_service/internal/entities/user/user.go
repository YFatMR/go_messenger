package user

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
)

type Entity struct {
	name    string
	surname string
}

func New(name string, surname string) *Entity {
	return &Entity{
		name:    name,
		surname: surname,
	}
}

func FromProtobuf(user *proto.UserData) (*Entity, error) {
	if user == nil || user.GetName() == "" || user.GetSurname() == "" {
		return nil, entities.ErrWrongRequestFormat
	}
	return New(user.GetName(), user.GetSurname()), nil
}

func (e *Entity) GetName() string {
	return e.name
}

func (e *Entity) GetSurname() string {
	return e.surname
}
