package grpcc

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// jwtManager := jwtmanager.FromConfig(config, logger)

type AuthorizationField struct{}

var AuthorizationFieldContextKey AuthorizationField

type GRPCHeaders struct {
	authorizationHeader string
	userIDHeader        string
	userRoleHeader      string
}

func NewGRPCHeaders(authorizationHeader string, userIDHeader string, userRoleHeader string) GRPCHeaders {
	return GRPCHeaders{
		authorizationHeader: authorizationHeader,
		userIDHeader:        userIDHeader,
		userRoleHeader:      userRoleHeader,
	}
}

func GRPCHeadersFromConfig(config *cviper.CustomViper) GRPCHeaders {
	return GRPCHeaders{
		authorizationHeader: config.GetStringRequired("GRPC_AUTHORIZARION_HEADER"),
		userIDHeader:        config.GetStringRequired("GRPC_USER_ID_HEADER"),
		userRoleHeader:      config.GetStringRequired("GRPC_USER_ROLE_HEADER"),
	}
}

func UnaryAuthInterceptor(jwtManager jwtmanager.Manager, headers GRPCHeaders,
	logger *czap.Logger,
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

		if !authorizationRequired {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		grpcMetadata, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.ErrorContext(ctx, "Got bad context", zap.String("method", method))
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		tokens := grpcMetadata.Get(headers.authorizationHeader)
		if len(tokens) == 0 {
			logger.ErrorContext(ctx, "Medatada has no token", zap.String("method", method))
			return ErrAccessDenied
		} else if len(tokens) > 1 {
			logger.ErrorContext(
				ctx, "Medatada has a lot of tokens", zap.String("method", method), zap.Int("tokens count", len(tokens)),
			)
		}
		accessToken := tokens[0]

		claims, err := jwtManager.VerifyToken(ctx, accessToken)
		if err != nil {
			return err
		}

		// add payload to metadata
		ctxMetadata := metadata.Pairs(
			headers.userIDHeader, claims.UserID,
			headers.userRoleHeader, claims.UserRole,
		)

		ctx = metadata.NewOutgoingContext(ctx, ctxMetadata)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
