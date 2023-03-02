package unsafecredential

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/entities"
	"github.com/YFatMR/go_messenger/user_service/entities/userrole"
)

type Entity struct {
	login    string
	password string
	role     *userrole.Entity
}

func New(login string, password string, role *userrole.Entity) *Entity {
	return &Entity{
		login:    login,
		password: password,
		role:     role,
	}
}

func FromProtobuf(credential *proto.Credential) (
	*Entity, error,
) {
	if credential == nil || credential.GetLogin() == "" || credential.GetPassword() == "" || credential.GetRole() == "" {
		return nil, entities.ErrWrongRequestFormat
	}
	role, err := userrole.FromString(credential.GetRole())
	if err != nil {
		return nil, err
	}
	return New(credential.GetLogin(), credential.GetPassword(), role), nil
}

func (e *Entity) GetLogin() string {
	return e.login
}

func (e *Entity) GetPassword() string {
	return e.password
}

func (e *Entity) GetRole() *userrole.Entity {
	return e.role
}
