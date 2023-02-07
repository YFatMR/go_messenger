package interceptors

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthorizationField struct{}

var AuthorizationFieldContextKey AuthorizationField

func UnaryAuthInterceptor(authClient proto.AuthClient, grpcAuthorizationHeader string,
	grpcAuthorizationAccountIDHeader string, grpcAuthorizationUserRoleHeader string,
	logger *loggers.OtelZapLoggerWithTraceID,
) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		authorization := ctx.Value(AuthorizationField{})
		if authorization == nil {
			logger.ErrorContext(
				ctx, "gRPC method without authorization type. Please, add AuthorizationField to context.",
				zap.String("method", method),
			)
			return ErrAccessDenied
		}
		authorizationRequired := authorization.(bool)
		logger.DebugContextNoExport(
			ctx, "Authorization required",
			zap.String("method", method), zap.Bool("required", authorizationRequired),
		)

		if !authorizationRequired {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		grpcMetadata, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.ErrorContext(ctx, "Got bad context", zap.String("method", method))
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		tokens := grpcMetadata.Get(grpcAuthorizationHeader)
		if len(tokens) == 0 {
			logger.ErrorContextNoExport(ctx, "Medatada has no token", zap.String("method", method))
			return ErrAccessDenied
		} else if len(tokens) > 1 {
			logger.ErrorContextNoExport(
				ctx, "Medatada has a lot of tokens", zap.String("method", method), zap.Int("tokens count", len(tokens)),
			)
		}
		token := tokens[0]

		account, err := authClient.GetTokenPayload(ctx, &proto.Token{
			AccessToken: token,
		})
		if err != nil {
			logger.ErrorContextNoExport(ctx, "Can't get token payload", zap.String("method", method), zap.Error(err))
			return ErrAccessDenied
		}

		// add payload to metadata
		ctxMetadata := metadata.Pairs(
			grpcAuthorizationAccountIDHeader, account.GetAccountID(),
			grpcAuthorizationUserRoleHeader, account.GetUserRole(),
		)
		logger.InfoContextNoExport(
			ctx, "Authorized", zap.String("method", method),
			zap.String(grpcAuthorizationAccountIDHeader, account.GetAccountID()),
			zap.String(grpcAuthorizationUserRoleHeader, account.GetUserRole()),
		)

		ctx = metadata.NewOutgoingContext(ctx, ctxMetadata)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
