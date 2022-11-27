package controllers

import (
	"context"
	proto "protocol/pkg/proto"
	"user_server/internal/enities"
)

type userService interface {
	Create(ctx context.Context, request *enities.User) (string, error)
	GetById(ctx context.Context, id string) *enities.User
}

type UserController struct {
	service userService
}

func NewUserController(service userService) *UserController {
	return &UserController{
		service: service,
	}
}

func (s *UserController) Create(ctx context.Context, request *proto.UserData) (*proto.UserId, error) {
	user := enities.NewUser(request.Name, request.Surname)
	insertedId, err := s.service.Create(ctx, user)
	if err != nil {
		panic(err)
	}
	return &proto.UserId{
		Id: insertedId,
	}, nil
}

func (s *UserController) GetById(ctx context.Context, request *proto.UserId) (*proto.UserData, error) {
	userId := request.GetId()
	userData := s.service.GetById(ctx, userId)
	if userData == nil {
		// TODO: unexist user
		return nil, nil
	}
	return &proto.UserData{
		Name:    userData.GetName(),
		Surname: userData.GetSurname(),
	}, nil
}
