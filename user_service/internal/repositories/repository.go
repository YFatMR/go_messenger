package repositories

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/userid"
)

type UserRepository interface {
	Create(ctx context.Context, user *user.Entity, accountID *accountid.Entity) (
		userID *userid.Entity, cerr cerrors.Error,
	)
	GetByID(ctx context.Context, userID *userid.Entity) (user *user.Entity, cerr cerrors.Error)
	DeleteByID(ctx context.Context, userID *userid.Entity) (cerr cerrors.Error)
}
