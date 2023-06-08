package main

import (
	"github.com/YFatMR/go_messenger/bots_service/internal/apientity"
	"github.com/YFatMR/go_messenger/bots_service/internal/domain"
	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/sashabaranov/go-openai"
)

func BotsModelFromConfig(config *cviper.CustomViper, logger *czap.Logger) apientity.BotsModel {
	openaiClient := openai.NewClient(config.GetStringRequired("OPENAI_API_KEY"))
	openaiSettings := domain.OpenAISettings{
		MaxTokens: config.GetIntRequired("OPENAI_MAX_TOKENS_COUNT"),
	}
	return domain.NewBotsModel(openaiClient, openaiSettings, logger)
}

func BotsControllerFromConfig(config *cviper.CustomViper, logger *czap.Logger) apientity.BotsController {
	model := BotsModelFromConfig(config, logger)
	return domain.NewBotsController(model, logger)
}
