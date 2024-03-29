package sandbox

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/apientity"
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
	"go.uber.org/zap"
)

type controller struct {
	sandboxService apientity.SandboxService
	contextManager apientity.ContextManager
	logger         *czap.Logger
}

func NewController(sandboxService apientity.SandboxService, contextManager apientity.ContextManager,
	logger *czap.Logger,
) apientity.SandboxController {
	return &controller{
		sandboxService: sandboxService,
		contextManager: contextManager,
		logger:         logger,
	}
}

func (c *controller) GetProgramByID(ctx context.Context, request *proto.ProgramID) (
	*proto.Program, error,
) {
	programID, err := entity.ProgramIDFromProtobuf(request)
	if err != nil {
		return nil, err
	}
	program, err := c.sandboxService.GetProgramByID(ctx, programID)
	if err != nil {
		return nil, err
	}
	result := entity.ProgramToProtobuf(program)
	c.logger.Debug("program  src result", zap.String("source", result.Source.SourceCode))
	if program == nil {
		c.logger.Error("Got null program!")
	}
	return result, nil
}

func (c *controller) CreateProgram(ctx context.Context, request *proto.ProgramSource) (
	*proto.ProgramID, error,
) {
	programSource, err := entity.ProgramSourceFromProtobuf(request)
	if err != nil {
		return nil, err
	}
	programID, err := c.sandboxService.CreateProgram(ctx, programSource)
	if err != nil {
		return nil, err
	}
	return entity.ProgramIDToProtobuf(programID), nil
}

func (c *controller) UpdateProgramSource(ctx context.Context,
	request *proto.UpdateProgramSourceRequest,
) (
	*proto.Void, error,
) {
	programSource, err := entity.ProgramSourceFromProtobuf(request.GetProgramSource())
	if err != nil {
		return nil, err
	}
	programID, err := entity.ProgramIDFromProtobuf(request.GetProgramID())
	if err != nil {
		return nil, err
	}
	err = c.sandboxService.UpdateProgramSource(ctx, programID, programSource)
	if err != nil {
		return nil, err
	}
	return entity.VoidProtobuf(), nil
}

func (c *controller) RunProgram(ctx context.Context, request *proto.ProgramID) (
	*proto.Void, error,
) {
	programID, err := entity.ProgramIDFromProtobuf(request)
	if err != nil {
		return nil, err
	}
	userID, err := c.contextManager.UserIDFromContext(ctx)
	if err != nil {
		c.logger.ErrorContext(ctx, "Can not extract user ID from metedata", zap.Error(err))
		return entity.VoidProtobuf(), ErrNoMetadataKey
	}
	err = c.sandboxService.RunProgram(ctx, programID, userID)
	if err != nil {
		return nil, err
	}
	return entity.VoidProtobuf(), nil
}

func (c *controller) LintProgram(ctx context.Context, request *proto.ProgramID) (
	*proto.Void, error,
) {
	// TODO: implement
	// programID, err := entity.ProgramIDFromProtobuf(request)
	// if err != nil {
	// 	return nil, err
	// }
	// err = c.sandboxService.LintProgram(ctx, programID)
	// if err != nil {
	// 	return nil, err
	// }
	return entity.VoidProtobuf(), nil
}

func (c *controller) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
