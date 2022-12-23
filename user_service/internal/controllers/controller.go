package controllers

import (
	"context"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/enities"
	"go.opentelemetry.io/otel/trace"
)

type userService interface {
	Create(ctx context.Context, request *enities.User) (string, error)
	GetById(ctx context.Context, id string) (*enities.User, error)
}

type UserController struct {
	service userService
	logger  *loggers.OtelZapLoggerWithTraceID
	tracer  trace.Tracer
}

func NewUserController(service userService, logger *loggers.OtelZapLoggerWithTraceID, tracer trace.Tracer) *UserController {
	return &UserController{
		service: service,
		logger:  logger,
		tracer:  tracer,
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
