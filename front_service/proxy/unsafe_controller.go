package proxy

// import (
// 	"context"

// 	"github.com/YFatMR/go_messenger/core/pkg/czap"
// 	"github.com/YFatMR/go_messenger/front_server/apientity"
// 	"github.com/YFatMR/go_messenger/front_server/grpcapi"
// 	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
// 	"go.uber.org/zap"
// )

// type unsafeController struct {
// 	userServiceClient proto.UserClient
// 	logger            *czap.Logger
// }

// func NewUnsafeController(userServiceClient proto.UserClient, logger *czap.Logger) apientity.UnsafeProxyController {
// 	return &unsafeController{
// 		userServiceClient: userServiceClient,
// 		logger:            logger,
// 	}
// }

// func (c *unsafeController) CreateUser(ctx context.Context, request *proto.CreateUserFrontRequest) (
// 	*proto.UserID, error,
// ) {
// 	grpcCtx := context.WithValue(ctx, grpcapi.AuthorizationFieldContextKey, false)
// 	userID, err := c.userServiceClient.CreateUser(
// 		grpcCtx, &proto.CreateUserRequest{
// 			Credential: &proto.Credential{
// 				Email:    request.GetCredential().GetEmail(),
// 				Password: request.GetCredential().GetPassword(),
// 				Role:     "user",
// 			},
// 			UserData: request.GetUserData(),
// 		},
// 	)
// 	if err != nil {
// 		c.logger.ErrorContext(ctx, "Can't create user", zap.Error(err))
// 		return nil, err
// 	}
// 	return userID, nil
// }

// func (c *unsafeController) GenerateToken(ctx context.Context, request *proto.PublicCredential) (
// 	*proto.Token, error,
// ) {
// 	grpcCtx := context.WithValue(ctx, grpcapi.AuthorizationFieldContextKey, false)
// 	return c.userServiceClient.GenerateToken(grpcCtx, &proto.Credential{
// 		Email:    request.GetEmail(),
// 		Password: request.GetPassword(),
// 		Role:     "user",
// 	})
// }

// func (c *unsafeController) Ping(ctx context.Context, request *proto.Void) (
// 	*proto.Pong, error,
// ) {
// 	return &proto.Pong{
// 		Message: "pong",
// 	}, nil
// }
