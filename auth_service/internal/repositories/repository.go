package repositories

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
)

type AccountRepository interface {
	CreateAccount(context.Context, *entities.Credential, entities.Role) (_ *entities.AccountID, err error)
	GetTokenPayloadWithHashedPasswordByLogin(context.Context, string) (
		_ *entities.TokenPayload, hashedPassword string, err error,
	)
}
