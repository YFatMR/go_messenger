package controllers

import (
	"context"
	"go.uber.org/zap"
	proto "protocol/pkg/proto"
	"user_server/internal/enities"
)

type userService interface {
	Create(ctx context.Context, request *enities.User) (string, error)
	GetById(ctx context.Context, id string) (*enities.User, error)
}

type UserController struct {
	service userService
	logger  *zap.Logger
}

func NewUserController(service userService, logger *zap.Logger) *UserController {
	return &UserController{
		service: service,
		logger:  logger,
	}
}

func (s *UserController) Create(ctx context.Context, request *proto.UserData) (*proto.UserId, error) {
	user := enities.NewUser(request.Name, request.Surname)
	insertedId, err := s.service.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &proto.UserId{
		Id: insertedId,
	}, nil
}

func (s *UserController) GetById(ctx context.Context, request *proto.UserId) (*proto.UserData, error) {
	userId := request.GetId()
	userData, err := s.service.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &proto.UserData{
		Name:    userData.GetName(),
		Surname: userData.GetSurname(),
	}, nil
}
