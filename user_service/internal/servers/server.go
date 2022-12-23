package servers

import (
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

type userController interface {
	Create(ctx context.Context, request *proto.UserData) (*proto.UserId, error)
	GetById(ctx context.Context, request *proto.UserId) (*proto.UserData, error)
}

type GRPCUserServer struct {
	proto.UnimplementedUserServer
	controller userController
	logger     *loggers.OtelZapLoggerWithTraceID
	tracer     trace.Tracer
}

func NewGRPCUserServer(controller userController, logger *loggers.OtelZapLoggerWithTraceID, tracer trace.Tracer) *GRPCUserServer {
	return &GRPCUserServer{
		controller: controller,
		logger:     logger,
		tracer:     tracer,
	}
}

func (s *GRPCUserServer) CreateUser(ctx context.Context, request *proto.UserData) (*proto.UserId, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/GRPCUserServer CreateUser")
	defer span.End()

	return s.controller.Create(ctx, request)
}

func (s *GRPCUserServer) GetUserById(ctx context.Context, request *proto.UserId) (*proto.UserData, error) {
	return s.controller.GetById(ctx, request)
}
