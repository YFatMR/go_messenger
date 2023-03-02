package repositories

import (
	"context"

	"github.com/YFatMR/go_messenger/user_service/entities/account"
	"github.com/YFatMR/go_messenger/user_service/entities/credential"
	"github.com/YFatMR/go_messenger/user_service/entities/user"
	"github.com/YFatMR/go_messenger/user_service/entities/userid"
)

type UserRepository interface {
	Create(ctx context.Context, user *user.Entity, credential *credential.Entity) (
		userID *userid.Entity, err error,
	)
	GetByID(ctx context.Context, userID *userid.Entity) (
		user *user.Entity, err error,
	)
	DeleteByID(ctx context.Context, userID *userid.Entity) (
		err error,
	)
	GetAccountByLogin(ctx context.Context, login string) (
		account *account.Entity, err error,
	)
}
