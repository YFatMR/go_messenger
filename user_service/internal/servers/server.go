package servers

import (
	. "core/pkg/loggers"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
	proto "protocol/pkg/proto"
)

type userController interface {
	Create(ctx context.Context, request *proto.UserData) (*proto.UserId, error)
	GetById(ctx context.Context, request *proto.UserId) (*proto.UserData, error)
}

type GRPCUserServer struct {
	proto.UnimplementedUserServer
	controller userController
	logger     *OtelZapLoggerWithTraceID
	tracer     trace.Tracer
}

func NewGRPCUserServer(controller userController, logger *OtelZapLoggerWithTraceID, tracer trace.Tracer) *GRPCUserServer {
	return &GRPCUserServer{
		controller: controller,
		logger:     logger,
		tracer:     tracer,
	}
}

func (s *GRPCUserServer) CreateUser(ctx context.Context, request *proto.UserData) (*proto.UserId, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/CreateUser - grpc")
	defer span.End()

	return s.controller.Create(ctx, request)
}

func (s *GRPCUserServer) GetUserById(ctx context.Context, request *proto.UserId) (*proto.UserData, error) {
	return s.controller.GetById(ctx, request)
}
