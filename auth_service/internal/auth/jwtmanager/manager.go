package jwtmanager

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/auth_service/internal/auth"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/token"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"
	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	AccountID string
	UserRole  entities.Role
}

func (c *TokenClaims) GetTokenPayload() *tokenpayload.Entity {
	return tokenpayload.New(c.AccountID, c.UserRole)
}

type Manager struct {
	secretKey               string
	tokenExpirationDuration time.Duration
	signingMethod           jwt.SigningMethod
}

func New(secretKey string, tokenExpirationDuration time.Duration) *Manager {
	return &Manager{
		secretKey:               secretKey,
		tokenExpirationDuration: tokenExpirationDuration,
		signingMethod:           jwt.SigningMethodHS256,
	}
}

func (m *Manager) GenerateToken(ctx context.Context, payload *tokenpayload.Entity) (*token.Entity, cerrors.Error) {
	if payload == nil {
		return nil, cerrors.New("null payload got", auth.ErrTokenGenerationFailed, auth.ErrTokenGenerationFailed)
	}
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenExpirationDuration)),
		},
		AccountID: payload.GetAccountID(),
		UserRole:  payload.GetUserRole(),
	}

	resultToken := jwt.NewWithClaims(m.signingMethod, claims)
	accessToken, err := resultToken.SignedString([]byte(m.secretKey))
	if err != nil {
		return nil, cerrors.New("can't generate signed string", err, auth.ErrTokenGenerationFailed)
	}
	return token.New(accessToken), nil
}

// Check token expiration withount direct checks.
func (m *Manager) VerifyToken(ctx context.Context, accessToken *token.Entity) (*TokenClaims, cerrors.Error) {
	if accessToken == nil {
		return nil, cerrors.New("null accessToken got", auth.ErrInvalidAccessToken, auth.ErrInvalidAccessToken)
	}
	token, err := jwt.ParseWithClaims(
		accessToken.GetAccessToken(),
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, auth.ErrInvalidAccessToken
			}
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		return nil, cerrors.New("invalid access token got", err, auth.ErrInvalidAccessToken)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, cerrors.New("invalid token claims", auth.ErrInvalidAccessToken, auth.ErrInvalidAccessToken)
	}
	return claims, nil
}
