package account

import "github.com/YFatMR/go_messenger/user_service/entities/userrole"

type Entity struct {
	userID         string
	login          string
	hashedPassword string
	role           *userrole.Entity
	nickname       string
	name           string
	surname        string
}

func New(userID string, login string, hashedPassword string, role *userrole.Entity, nickname string,
	name string, surname string,
) *Entity {
	return &Entity{
		userID:         userID,
		login:          login,
		hashedPassword: hashedPassword,
		role:           role,
		nickname:       nickname,
		name:           name,
		surname:        surname,
	}
}

func (e *Entity) GetUserID() string {
	return e.userID
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

func (e *Entity) GetNickname() string {
	return e.nickname
}

func (e *Entity) GetName() string {
	return e.name
}

func (e *Entity) GetSurname() string {
	return e.surname
}
