package userserver

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// endpoints

type frontUserServer struct {
	proto.UnimplementedFrontUserServer
	userServerAddress string
	logger            *loggers.OtelZapLoggerWithTraceID
	tracer            trace.Tracer
}

func newFrontUserServer(userServerAddress string, logger *loggers.OtelZapLoggerWithTraceID,
	tracer trace.Tracer,
) *frontUserServer {
	return &frontUserServer{
		userServerAddress: userServerAddress,
		logger:            logger,
		tracer:            tracer,
	}
}

func (s *frontUserServer) CreateUser(ctx context.Context, request *proto.UserData) (*proto.UserID, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/CreateUserSpan")
	defer span.End()

	s.logger.DebugContextNoExport(ctx, "called CreateUser endpoint")
	conn, err := grpc.Dial(
		s.userServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	client := proto.NewUserClient(conn)
	return client.CreateUser(ctx, request)
}

func (s *frontUserServer) GetUserByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/GetUserByIDSpan")
	defer span.End()

	s.logger.DebugContextNoExport(ctx, "called GetUserByID endpoint")
	conn, err := grpc.Dial(
		s.userServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	client := proto.NewUserClient(conn)
	return client.GetUserByID(ctx, request)
}
