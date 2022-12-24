package userserver

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
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

func (s *frontUserServer) GetUserByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	s.logger.DebugContextNoExport(ctx, "called GetUserByID endpoint")
	conn, err := grpc.Dial(s.userServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := proto.NewUserClient(conn)
	return client.GetUserByID(ctx, request)
}
