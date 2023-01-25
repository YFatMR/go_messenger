package decorators

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/decorators"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories"
)

type LoggingRepositoryDecorator struct {
	repositories.UserRepository
	logger *loggers.OtelZapLoggerWithTraceID
}

func NewLoggingRepositoryDecorator(repository repositories.UserRepository, logger *loggers.OtelZapLoggerWithTraceID,
) repositories.UserRepository {
	return &LoggingRepositoryDecorator{
		UserRepository: repository,
		logger:         logger,
	}
}

func (d *LoggingRepositoryDecorator) Create(ctx context.Context, user *entities.User, accountID *entities.AccountID) (
	_ *entities.UserID, err error,
) {
	callback := func() (*entities.UserID, error) { return d.UserRepository.Create(ctx, user, accountID) }
	return decorators.LogCallbackErrorWithReturnType(
		ctx, d.logger, "Create", "UserRepository", callback,
	)
}

func (d *LoggingRepositoryDecorator) GetByID(ctx context.Context, userID *entities.UserID) (
	_ *entities.User, err error,
) {
	callback := func() (_ *entities.User, err error) { return d.UserRepository.GetByID(ctx, userID) }
	return decorators.LogCallbackErrorWithReturnType(
		ctx, d.logger, "GetByID", "UserRepository", callback,
	)
}

func (d *LoggingRepositoryDecorator) DeleteByID(ctx context.Context, userID *entities.UserID) (err error) {
	callback := func() error { return d.UserRepository.DeleteByID(ctx, userID) }
	return decorators.LogCallbackError(
		ctx, d.logger, "DeleteByID", "UserRepository", callback,
	)
}
