package controllers

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
	"go.opentelemetry.io/otel/trace"
)

type userService interface {
	Create(context.Context, *entities.User, *entities.AccountID) (*entities.UserID, error)
	GetByID(context.Context, *entities.UserID) (*entities.User, error)
	DeleteByID(context.Context, *entities.UserID) error
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

func (c *UserController) Create(ctx context.Context, request *proto.CreateUserDataRequest) (*proto.UserID, error) {
	user, err := entities.NewUserFromProtobuf(request.GetUserData())
	if err != nil {
		c.logger.DebugContextNoExport(ctx, "Wrong format for user data")
		return nil, ErrWrongRequestFormat
	}
	c.logger.DebugContextNoExport(ctx, "User data parsed successfully")

	accountID, err := entities.NewAccountIDFromProtobuf(request.GetAccountID())
	if err != nil {
		c.logger.DebugContextNoExport(ctx, "Wrong format for account id")
		return nil, ErrWrongRequestFormat
	}
	c.logger.DebugContextNoExport(ctx, "Account id parsed successfully")

	insertedID, err := c.service.Create(ctx, user, accountID)
	if err != nil {
		return nil, err
	}

	return &proto.UserID{
		ID: insertedID.GetID(),
	}, nil
}

func (c *UserController) GetByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	userID, err := entities.NewUserIDFromProtobuf(request)
	if err != nil {
		c.logger.DebugContextNoExport(ctx, "Wrong format for userID data")
		return nil, ErrWrongRequestFormat
	}
	c.logger.DebugContextNoExport(ctx, "User id parsed successfully")

	userData, err := c.service.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &proto.UserData{
		Name:    userData.GetName(),
		Surname: userData.GetSurname(),
	}, nil
}

func (c *UserController) DeleteByID(ctx context.Context, request *proto.UserID) (*proto.Void, error) {
	userID, err := entities.NewUserIDFromProtobuf(request)
	if err != nil {
		c.logger.DebugContextNoExport(ctx, "Wrong format for userID data")
		return nil, err
	}
	c.logger.DebugContextNoExport(ctx, "User id parsed successfully")

	err = c.service.DeleteByID(ctx, userID)
	return &proto.Void{}, err
}
