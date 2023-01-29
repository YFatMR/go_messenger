package userservice

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	accountid "github.com/YFatMR/go_messenger/user_service/internal/entities/account_id"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	userid "github.com/YFatMR/go_messenger/user_service/internal/entities/user_id"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories"
)

type UserService struct {
	userRepository repositories.UserRepository
}

func New(repository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: repository,
	}
}

func (s *UserService) Create(ctx context.Context, user *user.Entity, accountID *accountid.Entity) (
	*userid.Entity, cerrors.Error,
) {
	return s.userRepository.Create(ctx, user, accountID)
}

func (s *UserService) GetByID(ctx context.Context, userID *userid.Entity) (*user.Entity, cerrors.Error) {
	return s.userRepository.GetByID(ctx, userID)
}

func (s *UserService) DeleteByID(ctx context.Context, userID *userid.Entity) cerrors.Error {
	return s.userRepository.DeleteByID(ctx, userID)
}
