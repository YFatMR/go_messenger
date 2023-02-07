package tokenpayload

import "github.com/YFatMR/go_messenger/auth_service/internal/entities"

type Entity struct {
	accountID string
	userRole  entities.Role
}

func New(accountID string, userRole entities.Role) *Entity {
	return &Entity{
		accountID: accountID,
		userRole:  userRole,
	}
}

func (e *Entity) GetAccountID() string {
	return e.accountID
}

func (e *Entity) GetUserRole() entities.Role {
	return e.userRole
}
