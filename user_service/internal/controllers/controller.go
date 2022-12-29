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
	GetByID(ctx context.Context, id string) (*enities.User, error)
}

type UserController struct {
	service userService
	logger  *loggers.OtelZapLoggerWithTraceID
	tracer  trace.Tracer
}

func NewUserController(service userService, logger *loggers.OtelZapLoggerWithTraceID,
	tracer trace.Tracer,
) *UserController {
	return &UserController{
		service: service,
		logger:  logger,
		tracer:  tracer,
	}
}

func (s *UserController) Create(ctx context.Context, request *proto.UserData) (*proto.UserID, error) {
	user := enities.NewUser(request.Name, request.Surname)
	insertedID, err := s.service.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &proto.UserID{
		ID: insertedID,
	}, nil
}

func (s *UserController) GetByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	UserID := request.GetID()
	userData, err := s.service.GetByID(ctx, UserID)
	if err != nil {
		return nil, err
	}
	return &proto.UserData{
		Name:    userData.GetName(),
		Surname: userData.GetSurname(),
	}, nil
}
