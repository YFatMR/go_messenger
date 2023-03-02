package tokenpayload

import "github.com/YFatMR/go_messenger/user_service/entities/userrole"

type Entity struct {
	userID   string
	userRole *userrole.Entity
}

func New(userID string, userRole *userrole.Entity) *Entity {
	return &Entity{
		userID:   userID,
		userRole: userRole,
	}
}

func (e *Entity) GetUserID() string {
	return e.userID
}

func (e *Entity) GetUserRole() *userrole.Entity {
	return e.userRole
}
