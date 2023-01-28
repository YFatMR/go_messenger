package servers

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.opentelemetry.io/otel/trace"
)

type accountController interface {
	CreateAccount(context.Context, *proto.Credential) (*proto.AccountID, error)
	GetToken(context.Context, *proto.Credential) (*proto.Token, error)
	GetTokenPayload(ctx context.Context, request *proto.Token) (*proto.TokenPayload, error)
	Ping(ctx context.Context, request *proto.Void) (*proto.Pong, error)
}

type GRPCAuthServer struct {
	proto.UnimplementedAuthServer
	accountController accountController
	logger            *loggers.OtelZapLoggerWithTraceID
	tracer            trace.Tracer
}

func NewGRPCAuthServer(accountController accountController, logger *loggers.OtelZapLoggerWithTraceID,
	tracer trace.Tracer,
) GRPCAuthServer {
	return GRPCAuthServer{
		accountController: accountController,
		logger:            logger,
		tracer:            tracer,
	}
}

func (s *GRPCAuthServer) CreateAccount(ctx context.Context, request *proto.Credential) (
	*proto.AccountID, error,
) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/ServerCreateAccountSpan")
	defer span.End()

	return s.accountController.CreateAccount(ctx, request)
}

func (s *GRPCAuthServer) GetToken(ctx context.Context, request *proto.Credential) (*proto.Token, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/ServerCreateTokenSpan")
	defer span.End()

	return s.accountController.GetToken(ctx, request)
}

func (s *GRPCAuthServer) GetTokenPayload(ctx context.Context, request *proto.Token) (*proto.TokenPayload, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/ServerGetAccountByTokenSpan")
	defer span.End()

	return s.accountController.GetTokenPayload(ctx, request)
}

func (s *GRPCAuthServer) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, error) {
	return s.accountController.Ping(ctx, request)
}
