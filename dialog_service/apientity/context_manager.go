package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/dialog_service/entity"
)

type ContextManager interface {
	UserIDFromContext(ctx context.Context) (
		userID *entity.UserID, err error,
	)
}
