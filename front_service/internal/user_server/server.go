package user_server

import (
	"context"
	. "core/pkg/loggers"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "protocol/pkg/proto"
)

// endpoints

type frontUserServer struct {
	proto.UnimplementedFrontUserServer
	userServerAddress string
	logger            *OtelZapLoggerWithTraceID
	tracer            trace.Tracer
}

func newFrontUserServer(userServerAddress string, logger *OtelZapLoggerWithTraceID, tracer trace.Tracer) *frontUserServer {
	return &frontUserServer{
		userServerAddress: userServerAddress,
		logger:            logger,
		tracer:            tracer,
	}
}

func (s *frontUserServer) CreateUser(ctx context.Context, request *proto.UserData) (*proto.UserId, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/SpanCreateUser", trace.WithAttributes(attribute.String("extra.key1", "extra.value")))
	defer span.End()

	s.logger.DebugContextNoExport(ctx, "called CreateUser endpoint2")
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

func (s *frontUserServer) GetUserById(ctx context.Context, request *proto.UserId) (*proto.UserData, error) {
	s.logger.DebugContextNoExport(ctx, "called GetUserById endpoint")
	conn, err := grpc.Dial(s.userServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := proto.NewUserClient(conn)
	return client.GetUserById(ctx, request)
}
