package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/user_service/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User, credential *entity.Credential) (
		userID *entity.UserID, err error,
	)
	GetByID(ctx context.Context, userID *entity.UserID) (
		user *entity.User, err error,
	)
	DeleteByID(ctx context.Context, userID *entity.UserID) (
		err error,
	)
	GetAccountByEmail(ctx context.Context, email string) (
		account *entity.Account, err error,
	)
	UpdateUserData(ctx context.Context, userID *entity.UserID, request *entity.UpdateUserRequest) (
		err error,
	)
	GetUsersByPrefix(ctx context.Context, selfID *entity.UserID, prefix string, limit uint64) (
		resp []*entity.UserWithID, err error,
	)
}
