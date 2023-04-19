package grpcctx

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"
)

type ContextManager struct {
	UserIDHeader string
}

func NewContextManager(userIDHeader string) *ContextManager {
	return &ContextManager{
		UserIDHeader: userIDHeader,
	}
}

func (m *ContextManager) UserIDFromContext(ctx context.Context) (
	uint64, error,
) {
	grpcMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, ErrNoMetadata
	}
	userIDs := grpcMetadata.Get(m.UserIDHeader)
	if len(userIDs) != 1 {
		return 0, ErrUnexpectedMetadataAccountIDCount
	}
	return strconv.ParseUint(userIDs[0], 10, 64)
}
