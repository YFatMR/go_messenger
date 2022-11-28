package services

import (
	"context"
	"go.uber.org/zap"
	"user_server/internal/enities"
)

type userRepository interface {
	Create(ctx context.Context, request *enities.User) (string, error)
	GetById(ctx context.Context, id string) (*enities.User, error)
}

type UserService struct {
	repository userRepository
	logger     *zap.Logger
}

func NewUserService(repository userRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repository: repository,
		logger:     logger,
	}
}

func (s *UserService) Create(ctx context.Context, request *enities.User) (string, error) {
	return s.repository.Create(ctx, request)
}

func (s *UserService) GetById(ctx context.Context, id string) (*enities.User, error) {
	return s.repository.GetById(ctx, id)
}
