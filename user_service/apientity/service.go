package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/user_service/entity"
)

type UserService interface {
	Create(ctx context.Context, user *entity.User, unsafeCredential *entity.UnsafeCredential) (
		userID *entity.UserID, err error,
	)
	GetByID(ctx context.Context, userID *entity.UserID) (
		user *entity.User, err error,
	)
	DeleteByID(ctx context.Context, userID *entity.UserID) (
		err error,
	)
	GenerateToken(ctx context.Context, unsafeCredential *entity.UnsafeCredential) (
		token *entity.Token, err error,
	)
	UpdateUserData(ctx context.Context, userID *entity.UserID, request *entity.UpdateUserRequest) (
		err error,
	)
	GetUsersByPrefix(ctx context.Context, selfID *entity.UserID, prefix string, limit uint64) (
		usersData []*entity.UserWithID, err error,
	)
}
