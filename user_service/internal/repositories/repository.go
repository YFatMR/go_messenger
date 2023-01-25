package repositories

import (
	"context"

	"github.com/YFatMR/go_messenger/user_service/internal/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User, accountID *entities.AccountID) (*entities.UserID, error)
	GetByID(ctx context.Context, userID *entities.UserID) (*entities.User, error)
	DeleteByID(ctx context.Context, userID *entities.UserID) error
}
