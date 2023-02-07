package userservice

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/userid"
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
	*userid.Entity, ulo.LogStash, error,
) {
	return s.userRepository.Create(ctx, user, accountID)
}

func (s *UserService) GetByID(ctx context.Context, userID *userid.Entity) (
	*user.Entity, ulo.LogStash, error,
) {
	return s.userRepository.GetByID(ctx, userID)
}

func (s *UserService) DeleteByID(ctx context.Context, userID *userid.Entity) (
	ulo.LogStash, error,
) {
	return s.userRepository.DeleteByID(ctx, userID)
}
