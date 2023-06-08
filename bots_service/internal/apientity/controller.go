package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type BotsController interface {
	GetBotMessageCompletion(ctx context.Context, request *proto.BotMessage) (
		response *proto.BotMessage, err error,
	)
	Ping(ctx context.Context, request *proto.Void) (
		response *proto.Pong, err error,
	)
}
