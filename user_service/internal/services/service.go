package services

import (
	"context"
	"errors"
	"user_server/internal/enities"
	"user_server/internal/repositories"
)

type userRepository interface {
	Create(ctx context.Context, request *enities.User) (string, error)
	GetById(ctx context.Context, id string) (*enities.User, error)
}

type UserService struct {
	repository userRepository
}

func NewUserService(repository userRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Create(ctx context.Context, request *enities.User) (string, error) {
	return s.repository.Create(ctx, request)
}

func (s *UserService) GetById(ctx context.Context, id string) *enities.User {
	p, err := s.repository.GetById(ctx, id)
	if errors.Is(err, repositories.UserNotFoundErr) {
		return nil
	} else if err != nil {
		//
		return nil
	}
	return p
}
