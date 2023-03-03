package user

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/apientity"
	"github.com/YFatMR/go_messenger/user_service/entity"
	"go.uber.org/zap"
)

type userController struct {
	service apientity.UserService
	logger  *czap.Logger
}

func NewController(service apientity.UserService, logger *czap.Logger) apientity.UserController {
	return &userController{
		service: service,
		logger:  logger,
	}
}

func (c *userController) Create(ctx context.Context, request *proto.CreateUserRequest) (
	*proto.UserID, error,
) {
	user, err := entity.UserFromProtobuf(request.GetUserData())
	if err != nil {
		return nil, ErrWrongRequestFormat
	}

	unsafeCredential, err := entity.UnsafeCredentialFromProtobuf(request.GetCredential())
	if err != nil {
		return nil, ErrWrongRequestFormat
	}

	userID, err := c.service.Create(ctx, user, unsafeCredential)
	if err != nil {
		return nil, err
	}
	return entity.UserIDToProtobuf(userID), nil
}

func (c *userController) GetByID(ctx context.Context, request *proto.UserID) (
	*proto.UserData, error,
) {
	userID, err := entity.UserIDFromProtobuf(request)
	if err != nil {
		c.logger.ErrorContext(ctx, "wrong format for userID data", zap.Error(err))
		return nil, ErrWrongRequestFormat
	}

	user, err := c.service.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return entity.UserToProtobuf(user), nil
}

func (c *userController) DeleteByID(ctx context.Context, request *proto.UserID) (
	*proto.Void, error,
) {
	userID, err := entity.UserIDFromProtobuf(request)
	if err != nil {
		c.logger.ErrorContext(ctx, "wrong format for userID data", zap.Error(err))
		return nil, ErrWrongRequestFormat
	}

	err = c.service.DeleteByID(ctx, userID)
	return entity.VoidProtobuf(), err
}

func (c *userController) GenerateToken(ctx context.Context, request *proto.Credential) (
	*proto.Token, error,
) {
	unsafeCredential, err := entity.UnsafeCredentialFromProtobuf(request)
	if err != nil {
		c.logger.ErrorContext(ctx, "wrong format for credential", zap.Error(err))
		return nil, ErrWrongRequestFormat
	}
	token, err := c.service.GenerateToken(ctx, unsafeCredential)
	if err != nil {
		return nil, err
	}
	return entity.TokenToProtobuf(token), nil
}

func (c *userController) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	return entity.PongProtobuf(), nil
}
