package grpcapi

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/grpcctx"
	"github.com/YFatMR/go_messenger/user_service/apientity"
	"github.com/YFatMR/go_messenger/user_service/entity"
)

type contextManager struct {
	headers Headers
	grpcctx.ContextManager
}

func NewContextManager(headers Headers) apientity.ContextManager {
	return &contextManager{
		headers: headers,
		ContextManager: grpcctx.ContextManager{
			UserIDHeader: headers.UserID,
		},
	}
}

func (m *contextManager) UserIDFromContext(ctx context.Context) (
	*entity.UserID, error,
) {
	userID, err := m.ContextManager.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &entity.UserID{ID: userID}, nil
}
