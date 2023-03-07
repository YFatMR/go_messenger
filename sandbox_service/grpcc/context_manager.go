package grpcc

import (
	"context"

	"github.com/YFatMR/go_messenger/sandbox_service/apientity"
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
	"google.golang.org/grpc/metadata"
)

type contextManager struct {
	headers Headers
}

func NewContextManager(headers Headers) apientity.ContextManager {
	return &contextManager{
		headers: headers,
	}
}

func (m contextManager) UserIDFromContext(ctx context.Context) (
	*entity.UserID, error,
) {
	grpcMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrNoMetadata
	}
	userIDs := grpcMetadata.Get(m.headers.UserID)
	if len(userIDs) != 1 {
		return nil, ErrUnexpectedMetadataAccountIDCount
	}
	return &entity.UserID{ID: userIDs[0]}, nil
}
