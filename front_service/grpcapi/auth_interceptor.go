package grpcapi

import (
	"context"
	"strconv"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// jwtManager := jwtmanager.FromConfig(config, logger)

type RequireAuthorizationField struct{}

type AuthorizationField struct{}

type GRPCHeaders struct {
	AuthorizationHeader string
	UserIDHeader        string
	UserRoleHeader      string
}

func NewGRPCHeaders(authorizationHeader string, userIDHeader string, userRoleHeader string) GRPCHeaders {
	return GRPCHeaders{
		AuthorizationHeader: authorizationHeader,
		UserIDHeader:        userIDHeader,
		UserRoleHeader:      userRoleHeader,
	}
}

func GRPCHeadersFromConfig(config *cviper.CustomViper) GRPCHeaders {
	return GRPCHeaders{
		AuthorizationHeader: config.GetStringRequired("GRPC_AUTHORIZARION_HEADER"),
		UserIDHeader:        config.GetStringRequired("GRPC_USER_ID_HEADER"),
		UserRoleHeader:      config.GetStringRequired("GRPC_USER_ROLE_HEADER"),
	}
}

func AccessTokenFromContext(ctx context.Context, authHeader string, logger *czap.Logger) (
	string, error,
) {
	grpcMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.ErrorContext(ctx, "Got bad context for auth")
		return "", ErrAccessDenied
	}

	tokens := grpcMetadata.Get(authHeader)
	if len(tokens) == 0 {
		logger.ErrorContext(ctx, "Medatada has no token")
		return "", ErrAccessDenied
	} else if len(tokens) > 1 {
		logger.ErrorContext(
			ctx, "Medatada has a lot of tokens", zap.Int("tokens count", len(tokens)),
		)
		return "", ErrAccessDenied
	}
	return tokens[0], nil
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
		authorization := ctx.Value(RequireAuthorizationField{})
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

		accessToken := ctx.Value(AuthorizationField{})
		claims, err := jwtManager.VerifyToken(ctx, accessToken.(string))
		if err != nil {
			return err
		}

		// add payload to metadata
		ctxMetadata := metadata.Pairs(
			headers.UserIDHeader, strconv.FormatUint(claims.UserID, 10),
			headers.UserRoleHeader, claims.UserRole,
		)

		ctx = metadata.NewOutgoingContext(ctx, ctxMetadata)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
