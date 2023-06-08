package grpcapi

import (
	"context"

	"github.com/YFatMR/go_messenger/bots_service/internal/apientity"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type Server struct {
	proto.UnimplementedBotsServiceServer
	controller apientity.BotsController
}

func NewServer(controller apientity.BotsController) *Server {
	return &Server{
		controller: controller,
	}
}

func (s *Server) GetBotMessageCompletion(ctx context.Context, request *proto.BotMessage) (
	*proto.BotMessage, error,
) {
	return s.controller.GetBotMessageCompletion(ctx, request)
}

func (s *Server) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	return s.controller.Ping(ctx, request)
}
