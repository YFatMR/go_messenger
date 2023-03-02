package credential

import (
	"github.com/YFatMR/go_messenger/user_service/entities/userrole"
)

type Entity struct {
	login          string
	hashedPassword string
	role           *userrole.Entity
}

func New(login string, hashedPassword string, role *userrole.Entity) *Entity {
	return &Entity{
		login:          login,
		hashedPassword: hashedPassword,
		role:           role,
	}
}

func (e *Entity) GetLogin() string {
	return e.login
}

func (e *Entity) GetHashedPassword() string {
	return e.hashedPassword
}

func (e *Entity) GetRole() *userrole.Entity {
	return e.role
}
