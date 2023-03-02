package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

type ContextManager interface {
	UserIDFromContext(ctx context.Context) (
		userID *entity.UserID, err error,
	)
}
