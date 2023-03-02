package services

import (
	"context"

	"github.com/YFatMR/go_messenger/user_service/entities/token"
	"github.com/YFatMR/go_messenger/user_service/entities/unsafecredential"
	"github.com/YFatMR/go_messenger/user_service/entities/user"
	"github.com/YFatMR/go_messenger/user_service/entities/userid"
)

type UserService interface {
	Create(ctx context.Context, user *user.Entity, unsafeCredential *unsafecredential.Entity) (
		userID *userid.Entity, err error,
	)
	GetByID(ctx context.Context, userID *userid.Entity) (
		user *user.Entity, err error,
	)
	DeleteByID(ctx context.Context, userID *userid.Entity) (
		err error,
	)
	GenerateToken(ctx context.Context, unsafeCredential *unsafecredential.Entity) (
		token *token.Entity, err error,
	)
}
