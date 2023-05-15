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

func (s *Server) CreateDialogMessageWithCode(ctx context.Context, request *proto.CreateDialogMessageWithCodeRequest) (
	*proto.CreateDialogMessageResponse, error,
) {
	return s.controller.CreateDialogMessageWithCode(ctx, request)
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

func (s *Server) ReadAllMessagesBeforeAndInclude(ctx context.Context, request *proto.ReadAllMessagesBeforeRequest) (
	*proto.Void, error,
) {
	return s.controller.ReadAllMessagesBeforeAndInclude(ctx, request)
}

func (s *Server) CreateInstruction(ctx context.Context, request *proto.CreateInstructionRequest) (
	*proto.InstructionID, error,
) {
	return s.controller.CreateInstruction(ctx, request)
}

func (s *Server) GetInstructions(ctx context.Context, request *proto.GetInstructionsRequest) (
	*proto.GetInstructionsResponse, error,
) {
	return s.controller.GetInstructions(ctx, request)
}

func (s *Server) GetInstructionsByID(ctx context.Context, request *proto.GetInstructionsByIDRequest) (
	*proto.GetInstructionsResponse, error,
) {
	return s.controller.GetInstructionsByID(ctx, request)
}

func (s *Server) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	return s.controller.Ping(ctx, request)
}
