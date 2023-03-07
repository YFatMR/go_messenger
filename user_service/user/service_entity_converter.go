package user

import (
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"github.com/YFatMR/go_messenger/user_service/entity"
)

func TokenPayloadFromAccount(account *entity.Account) jwtmanager.TokenPayload {
	return jwtmanager.TokenPayload{
		UserID:   account.UserID,
		UserRole: account.Role.Name,
	}
}
