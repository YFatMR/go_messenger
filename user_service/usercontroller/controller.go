package usercontroller

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/controllers"
	"github.com/YFatMR/go_messenger/user_service/entities/unsafecredential"
	"github.com/YFatMR/go_messenger/user_service/entities/user"
	"github.com/YFatMR/go_messenger/user_service/entities/userid"
	"github.com/YFatMR/go_messenger/user_service/services"
	"go.uber.org/zap"
)

type userController struct {
	service services.UserService
	logger  *czap.Logger
}

func New(service services.UserService, logger *czap.Logger) controllers.UserController {
	return &userController{
		service: service,
		logger:  logger,
	}
}

func (c *userController) Create(ctx context.Context, request *proto.CreateUserRequest) (
	*proto.UserID, error,
) {
	user, err := user.FromProtobuf(request.GetUserData())
	if err != nil {
		return nil, ErrWrongRequestFormat
	}

	unsafeCredential, err := unsafecredential.FromProtobuf(request.GetCredential())
	if err != nil {
		return nil, ErrWrongRequestFormat
	}

	insertedID, err := c.service.Create(ctx, user, unsafeCredential)
	if err != nil {
		return nil, err
	}

	return &proto.UserID{
		ID: insertedID.GetID(),
	}, nil
}

func (c *userController) GetByID(ctx context.Context, request *proto.UserID) (
	*proto.UserData, error,
) {
	userID, err := userid.FromProtobuf(request)
	if err != nil {
		c.logger.ErrorContext(ctx, "wrong format for userID data", zap.Error(err))
		return nil, ErrWrongRequestFormat
	}

	userData, err := c.service.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &proto.UserData{
		Name:    userData.GetName(),
		Surname: userData.GetSurname(),
	}, nil
}

func (c *userController) DeleteByID(ctx context.Context, request *proto.UserID) (
	*proto.Void, error,
) {
	userID, err := userid.FromProtobuf(request)
	if err != nil {
		c.logger.ErrorContext(ctx, "wrong format for userID data", zap.Error(err))
		return nil, ErrWrongRequestFormat
	}

	err = c.service.DeleteByID(ctx, userID)
	return &proto.Void{}, err
}

func (c *userController) GenerateToken(ctx context.Context, request *proto.Credential) (
	*proto.Token, error,
) {
	unsafeCredential, err := unsafecredential.FromProtobuf(request)
	if err != nil {
		c.logger.ErrorContext(ctx, "wrong format for credential", zap.Error(err))
		return nil, ErrWrongRequestFormat
	}
	token, err := c.service.GenerateToken(ctx, unsafeCredential)
	if err != nil {
		return nil, err
	}
	return &proto.Token{
		AccessToken: token.GetAccessToken(),
	}, nil
}

func (c *userController) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
