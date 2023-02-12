package userserver

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/front_server/internal/interceptors"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/status"
)

type FrontServer struct {
	proto.UnimplementedFrontServer
	authServiceClient    proto.AuthClient
	userServiceClient    proto.UserClient
	sandboxServiceClient proto.SandboxClient
	logger               *loggers.OtelZapLoggerWithTraceID
	tracer               trace.Tracer
	backoffConfig        backoff.Config
}

func NewFrontServer(authServiceClient proto.AuthClient,
	userServiceClient proto.UserClient, sandboxServiceClient proto.SandboxClient,
	logger *loggers.OtelZapLoggerWithTraceID,
	tracer trace.Tracer, backoffConfig backoff.Config,
) FrontServer {
	return FrontServer{
		authServiceClient:    authServiceClient,
		userServiceClient:    userServiceClient,
		sandboxServiceClient: sandboxServiceClient,
		logger:               logger,
		tracer:               tracer,
		backoffConfig:        backoffConfig,
	}
}

func (s *FrontServer) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (
	_ *proto.UserID, err error,
) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/CreateUserSpan")
	defer span.End()

	s.logger.InfoContextNoExport(ctx, "Creating account...", zap.String("login", request.GetCredential().GetLogin()))
	accountID, err := s.authServiceClient.CreateAccount(ctx, request.GetCredential())
	s.logger.InfoContextNoExport(ctx, "gRPC call CreateAccount", zap.String("status", status.Code(err).String()))

	if err != nil {
		s.logger.ErrorContext(
			ctx, "Can't create account", zap.Error(err), zap.String("login", request.GetCredential().GetLogin()),
		)
		return nil, err
	}
	s.logger.InfoContextNoExport(
		ctx, "Account created successfully", zap.String("login", request.GetCredential().GetLogin()),
	)

	grpcCtx := context.WithValue(ctx, interceptors.AuthorizationFieldContextKey, false)
	s.logger.InfoContextNoExport(ctx, "Creating user data...")
	userID, err := s.userServiceClient.CreateUser(
		grpcCtx, &proto.CreateUserDataRequest{
			AccountID: accountID,
			UserData:  request.GetUserData(),
		},
	)
	s.logger.InfoContextNoExport(ctx, "gRPC call CreateUser", zap.String("status", status.Code(err).String()))

	if err != nil {
		s.logger.ErrorContext(ctx, "Can't create user", zap.Error(err))
		return nil, err
	}
	s.logger.InfoContextNoExport(ctx, "User data created successfully")

	// TODO: delete account if err to create user data.
	return userID, nil
}

func (s *FrontServer) GetToken(ctx context.Context, request *proto.Credential) (*proto.Token, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/GenerateTokenSpan")
	defer span.End()

	s.logger.DebugContextNoExport(ctx, "called GetToken endpoint")
	return s.authServiceClient.GetToken(ctx, request)
}

func (s *FrontServer) GetUserByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/GetUserByIDSpan")
	defer span.End()

	s.logger.DebugContextNoExport(ctx, "called GetUserByID endpoint")
	grpcCtx := context.WithValue(ctx, interceptors.AuthorizationFieldContextKey, true)
	return s.userServiceClient.GetUserByID(grpcCtx, request)
}

func (s *FrontServer) Execute(ctx context.Context, request *proto.Program) (
	*proto.ProgramResult, error,
) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "/Execute")
	defer span.End()

	s.logger.DebugContextNoExport(ctx, "called Execute endpoint")
	grpcCtx := context.WithValue(ctx, interceptors.AuthorizationFieldContextKey, true)
	return s.sandboxServiceClient.Execute(grpcCtx, request)
}

func (s *FrontServer) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, error) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
