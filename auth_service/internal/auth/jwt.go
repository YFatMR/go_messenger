package auth

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/golang-jwt/jwt/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	AccountID string
	UserRole  entities.Role
}

func (c *TokenClaims) GetTokenPayload() *entities.TokenPayload {
	return entities.NewTokenPayload(c.AccountID, c.UserRole)
}

type JWTManager struct {
	secretKey               string
	tokenExpirationDuration time.Duration
	logger                  *loggers.OtelZapLoggerWithTraceID
	tracer                  trace.Tracer
	signingMethod           jwt.SigningMethod
}

func NewJWTManager(secretKey string, tokenExpirationDuration time.Duration,
	logger *loggers.OtelZapLoggerWithTraceID, tracer trace.Tracer,
) *JWTManager {
	return &JWTManager{
		secretKey:               secretKey,
		tokenExpirationDuration: tokenExpirationDuration,
		logger:                  logger,
		tracer:                  tracer,
		signingMethod:           jwt.SigningMethodHS256,
	}
}

func (m *JWTManager) GenerateToken(ctx context.Context, payload *entities.TokenPayload) (*entities.Token, error) {
	m.logger.DebugContextNoExport(ctx, "Generating token ...")
	if payload == nil {
		m.logger.ErrorContext(ctx, "null payload got")
		return nil, ErrTokenGenerationFailed
	}
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenExpirationDuration)),
		},
		AccountID: payload.GetAccountID(),
		UserRole:  payload.GetUserRole(),
	}

	token := jwt.NewWithClaims(m.signingMethod, claims)
	accessToken, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		m.logger.ErrorContext(ctx, "can't generate signed string", zap.Error(err))
		return nil, ErrTokenGenerationFailed
	}
	return entities.NewToken(accessToken), nil
}

// Check token expiration withount direct checks.
func (m *JWTManager) VerifyToken(ctx context.Context, accessToken *entities.Token) (*TokenClaims, error) {
	m.logger.DebugContextNoExport(ctx, "Verefying token...")
	if accessToken == nil {
		m.logger.ErrorContext(ctx, "null accessToken got")
		return nil, ErrInvalidAccessToken
	}
	token, err := jwt.ParseWithClaims(
		accessToken.GetAccessToken(),
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				m.logger.ErrorContext(ctx, "unexpected sign in method")
				return nil, ErrInvalidAccessToken
			}
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		m.logger.ErrorContext(ctx, "invalid access token got", zap.Error(err))
		return nil, ErrInvalidAccessToken
	}
	m.logger.DebugContextNoExport(ctx, "Token claims parsed successfully")

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		m.logger.ErrorContext(ctx, "invalid token claims")
		return nil, ErrInvalidAccessToken
	}
	m.logger.DebugContextNoExport(ctx, "Token verified")

	return claims, nil
}
