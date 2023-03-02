package proxy

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/front_server/apientity"
	"github.com/YFatMR/go_messenger/front_server/grpcc"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.uber.org/zap"
)

type unsafeController struct {
	userServiceClient proto.UserClient
	logger            *czap.Logger
}

func NewUnsafeController(userServiceClient proto.UserClient, logger *czap.Logger) apientity.UnsafeProxyController {
	return &unsafeController{
		userServiceClient: userServiceClient,
		logger:            logger,
	}
}

func (c *unsafeController) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (
	*proto.UserID, error,
) {
	grpcCtx := context.WithValue(ctx, grpcc.AuthorizationFieldContextKey, false)
	userID, err := c.userServiceClient.CreateUser(
		grpcCtx, &proto.CreateUserRequest{
			Credential: request.GetCredential(),
			UserData:   request.GetUserData(),
		},
	)
	if err != nil {
		c.logger.ErrorContext(ctx, "Can't create user", zap.Error(err))
		return nil, err
	}
	return userID, nil
}

func (c *unsafeController) GenerateToken(ctx context.Context, request *proto.Credential) (
	*proto.Token, error,
) {
	grpcCtx := context.WithValue(ctx, grpcc.AuthorizationFieldContextKey, false)
	return c.userServiceClient.GenerateToken(grpcCtx, request)
}

func (c *unsafeController) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
