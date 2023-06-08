package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/bots_service/internal/entity"
)

type BotsModel interface {
	GetBotMessageCompletion(ctx context.Context, message *entity.BotMessage) (
		response *entity.BotMessage, err error,
	)
}
