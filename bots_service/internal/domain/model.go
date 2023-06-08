package domain

import (
	"context"

	"github.com/YFatMR/go_messenger/bots_service/internal/apientity"
	"github.com/YFatMR/go_messenger/bots_service/internal/entity"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type OpenAISettings struct {
	MaxTokens int
}

type BotsModel struct {
	openaiClient   *openai.Client
	openaiSettings OpenAISettings
	logger         *czap.Logger
}

func NewBotsModel(openaiClient *openai.Client, openaiSettings OpenAISettings, logger *czap.Logger) apientity.BotsModel {
	return &BotsModel{
		openaiClient:   openaiClient,
		openaiSettings: openaiSettings,
		logger:         logger,
	}
}

func (m *BotsModel) GetBotMessageCompletion(ctx context.Context, message *entity.BotMessage) (
	*entity.BotMessage, error,
) {
	req := openai.CompletionRequest{
		Model:     openai.GPT3Ada,
		MaxTokens: m.openaiSettings.MaxTokens,
		Prompt:    message.Text,
	}
	response, err := m.openaiClient.CreateCompletion(ctx, req)
	if err != nil {
		m.logger.ErrorContext(ctx, "Can not create cgatGPT response", zap.Error(err))
		return nil, err
	}
	return &entity.BotMessage{
		Text: response.Choices[0].Text,
	}, nil
}
