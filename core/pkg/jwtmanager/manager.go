package jwtmanager

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type TokenPayload struct {
	UserID   uint64
	UserRole string
}

type TokenClaims struct {
	jwt.RegisteredClaims
	TokenPayload
}

type Manager interface {
	GenerateAccessToken(ctx context.Context, payload TokenPayload) (
		accessToken string, err error,
	)
	VerifyToken(ctx context.Context, accessToken string) (
		tokenClaims *TokenClaims, err error,
	)
}

type manager struct {
	secretKey               string
	tokenExpirationDuration time.Duration
	signingMethod           jwt.SigningMethod
	logger                  *czap.Logger
}

func New(secretKey string, tokenExpirationDuration time.Duration, logger *czap.Logger) Manager {
	return &manager{
		secretKey:               secretKey,
		tokenExpirationDuration: tokenExpirationDuration,
		signingMethod:           jwt.SigningMethodHS256,
		logger:                  logger,
	}
}

func FromConfig(config *cviper.CustomViper, logger *czap.Logger) Manager {
	authTokenSecretKey := config.GetStringRequired("AUTH_TOKEN_SECRET_KEY")
	authTokenExpirationDuration := config.GetSecondsDurationRequired("AUTH_TOKEN_EXPIRATION_SECONDS")
	return New(authTokenSecretKey, authTokenExpirationDuration, logger)
}

func (m *manager) GenerateAccessToken(ctx context.Context, payload TokenPayload) (
	/*accessToken*/ string, error,
) {
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenExpirationDuration)),
		},
		TokenPayload: payload,
	}

	resultToken := jwt.NewWithClaims(m.signingMethod, claims)
	accessToken, err := resultToken.SignedString([]byte(m.secretKey))
	if err != nil {
		m.logger.ErrorContext(ctx, "can't generate signed string", zap.Error(err))
		return "", ErrTokenGenerationFailed
	}
	return accessToken, nil
}

// Check token expiration withount direct checks.
func (m *manager) VerifyToken(ctx context.Context, accessToken string) (
	*TokenClaims, error,
) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrInvalidAccessToken
			}
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		m.logger.ErrorContext(ctx, "invalid access token got", zap.Error(err))
		return nil, ErrInvalidAccessToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		m.logger.ErrorContext(ctx, "invalid token claims", zap.Error(err))
		return nil, ErrInvalidAccessToken
	}
	return claims, nil
}
