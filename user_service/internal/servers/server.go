package servers

import (
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

type userController interface {
	Create(ctx context.Context, request *proto.CreateUserDataRequest) (*proto.UserID, error)
	GetByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error)
	DeleteByID(ctx context.Context, request *proto.UserID) (*proto.Void, error)
	Ping(ctx context.Context, request *proto.Void) (*proto.Pong, error)
}

type GRPCUserServer struct {
	proto.UnimplementedUserServer
	controller userController
	logger     *loggers.OtelZapLoggerWithTraceID
	tracer     trace.Tracer
}

func NewGRPCUserServer(controller userController, logger *loggers.OtelZapLoggerWithTraceID,
	tracer trace.Tracer,
) GRPCUserServer {
	return GRPCUserServer{
		controller: controller,
		logger:     logger,
		tracer:     tracer,
	}
}

func (s *GRPCUserServer) CreateUser(ctx context.Context, request *proto.CreateUserDataRequest) (*proto.UserID, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/ServerCreateUserSpan")
	defer span.End()

	return s.controller.Create(ctx, request)
}

func (s *GRPCUserServer) GetUserByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/ServerGetUserByIDSpan")
	defer span.End()

	return s.controller.GetByID(ctx, request)
}

func (s *GRPCUserServer) DeleteUserByID(ctx context.Context, request *proto.UserID) (*proto.Void, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/ServerDeleteUserByIDSpan")
	defer span.End()

	return s.controller.DeleteByID(ctx, request)
}

func (s *GRPCUserServer) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, error) {
	return s.controller.Ping(ctx, request)
}
