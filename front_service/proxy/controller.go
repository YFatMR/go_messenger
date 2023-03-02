package proxy

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/front_server/apientity"
	"github.com/YFatMR/go_messenger/front_server/grpcc"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type controller struct {
	userServiceClient    proto.UserClient
	sandboxServiceClient proto.SandboxClient
	logger               *czap.Logger
}

func NewController(userServiceClient proto.UserClient, sandboxServiceClient proto.SandboxClient, logger *czap.Logger,
) apientity.ProxyController {
	return &controller{
		userServiceClient:    userServiceClient,
		sandboxServiceClient: sandboxServiceClient,
		logger:               logger,
	}
}

func (c *controller) GetUserByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	grpcCtx := context.WithValue(ctx, grpcc.AuthorizationFieldContextKey, true)
	return c.userServiceClient.GetUserByID(grpcCtx, request)
}

func (c *controller) GetProgramByID(ctx context.Context, request *proto.ProgramID) (
	program *proto.Program, err error,
) {
	grpcCtx := context.WithValue(ctx, grpcc.AuthorizationFieldContextKey, true)
	return c.sandboxServiceClient.GetProgramByID(grpcCtx, request)
}

func (c *controller) CreateProgram(ctx context.Context, request *proto.ProgramSource) (
	programID *proto.ProgramID, err error,
) {
	grpcCtx := context.WithValue(ctx, grpcc.AuthorizationFieldContextKey, true)
	return c.sandboxServiceClient.CreateProgram(grpcCtx, request)
}

func (c *controller) UpdateProgramSource(ctx context.Context, request *proto.UpdateProgramSourceRequest) (
	void *proto.Void, err error,
) {
	grpcCtx := context.WithValue(ctx, grpcc.AuthorizationFieldContextKey, true)
	return c.sandboxServiceClient.UpdateProgramSource(grpcCtx, request)
}

func (c *controller) RunProgram(ctx context.Context, request *proto.ProgramID) (
	void *proto.Void, err error,
) {
	grpcCtx := context.WithValue(ctx, grpcc.AuthorizationFieldContextKey, true)
	return c.sandboxServiceClient.RunProgram(grpcCtx, request)
}

func (c *controller) LintProgram(ctx context.Context, request *proto.ProgramID) (
	void *proto.Void, err error,
) {
	grpcCtx := context.WithValue(ctx, grpcc.AuthorizationFieldContextKey, true)
	return c.sandboxServiceClient.LintProgram(grpcCtx, request)
}
