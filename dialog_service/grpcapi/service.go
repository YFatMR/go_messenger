package grpcapi

import (
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"golang.org/x/net/context"
)

type Server struct {
	proto.UnimplementedDialogServiceServer
	controller apientity.DialogController
}

func NewServer(controller apientity.DialogController) Server {
	return Server{
		controller: controller,
	}
}

func (s *Server) CreateDialogWith(ctx context.Context, request *proto.UserID) (
	*proto.Dialog, error,
) {
	return s.controller.CreateDialogWith(ctx, request)
}

func (s *Server) GetDialogByID(ctx context.Context, request *proto.DialogID) (
	*proto.Dialog, error,
) {
	return s.controller.GetDialogByID(ctx, request)
}

func (s *Server) GetDialogs(ctx context.Context, request *proto.GetDialogsRequest) (
	*proto.GetDialogsResponse, error,
) {
	return s.controller.GetDialogs(ctx, request)
}

func (s *Server) CreateDialogMessage(ctx context.Context, request *proto.CreateDialogMessageRequest) (
	*proto.CreateDialogMessageResponse, error,
) {
	return s.controller.CreateDialogMessage(ctx, request)
}

func (s *Server) GetDialogMessages(ctx context.Context, request *proto.GetDialogMessagesRequest) (
	*proto.GetDialogMessagesResponse, error,
) {
	return s.controller.GetDialogMessages(ctx, request)
}

func (s *Server) ReadAllMessagesBeforeAndIncl(ctx context.Context, request *proto.ReadAllMessagesBeforeRequest) (
	*proto.Void, error,
) {
	return s.controller.ReadAllMessagesBeforeAndIncl(ctx, request)
}

func (s *Server) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	return s.controller.Ping(ctx, request)
}
