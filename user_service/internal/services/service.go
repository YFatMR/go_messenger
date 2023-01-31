package services

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/userid"
)

type UserService interface {
	Create(ctx context.Context, user *user.Entity, accountID *accountid.Entity) (
		userID *userid.Entity, lerr logerr.Error,
	)
	GetByID(ctx context.Context, userID *userid.Entity) (user *user.Entity, lerr logerr.Error)
	DeleteByID(ctx context.Context, userID *userid.Entity) (lerr logerr.Error)
}
