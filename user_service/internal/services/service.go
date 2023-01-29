package services

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	accountid "github.com/YFatMR/go_messenger/user_service/internal/entities/account_id"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	userid "github.com/YFatMR/go_messenger/user_service/internal/entities/user_id"
)

type UserService interface {
	Create(ctx context.Context, user *user.Entity, accountID *accountid.Entity) (
		userID *userid.Entity, cerr cerrors.Error,
	)
	GetByID(ctx context.Context, userID *userid.Entity) (user *user.Entity, cerr cerrors.Error)
	DeleteByID(ctx context.Context, userID *userid.Entity) (cerr cerrors.Error)
}
