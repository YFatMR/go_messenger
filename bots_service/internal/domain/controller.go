package domain

import (
	"context"

	"github.com/YFatMR/go_messenger/bots_service/internal/apientity"
	"github.com/YFatMR/go_messenger/bots_service/internal/entity"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type BotsController struct {
	model  apientity.BotsModel
	logger *czap.Logger
}

func NewBotsController(model apientity.BotsModel, logger *czap.Logger) apientity.BotsController {
	return &BotsController{
		model:  model,
		logger: logger,
	}
}

func (c *BotsController) GetBotMessageCompletion(ctx context.Context, request *proto.BotMessage) (
	*proto.BotMessage, error,
) {
	inMessage, err := entity.GPTMessageFromProtobuf(request)
	if err != nil {
		return nil, err
	}

	outMessage, err := c.model.GetBotMessageCompletion(ctx, inMessage)
	if err != nil {
		return nil, err
	}
	return entity.GPTMessageToProtobuf(outMessage), nil
}

func (c *BotsController) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
