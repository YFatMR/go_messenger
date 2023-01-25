package decorators

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/decorators"
	"github.com/YFatMR/go_messenger/core/pkg/metrics/prometheus"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories"
)

type MetricRepositoryDecorator struct {
	repositories.UserRepository
}

func NewMetricDecorator(repository repositories.UserRepository) repositories.UserRepository {
	return &MetricRepositoryDecorator{
		UserRepository: repository,
	}
}

func (d *MetricRepositoryDecorator) Create(ctx context.Context, user *entities.User, accountID *entities.AccountID) (
	_ *entities.UserID, err error,
) {
	callback := func(ctx context.Context) (_ *entities.UserID, err error) {
		return d.UserRepository.Create(ctx, user, accountID)
	}
	return decorators.CollectMetricForDatabaseCallbackWithReturnType(ctx, prometheus.InsertOperationTag, callback)
}

func (d *MetricRepositoryDecorator) GetByID(ctx context.Context, userID *entities.UserID) (
	_ *entities.User, err error,
) {
	callback := func(ctx context.Context) (_ *entities.User, err error) { return d.UserRepository.GetByID(ctx, userID) }
	return decorators.CollectMetricForDatabaseCallbackWithReturnType(ctx, prometheus.FindOperationTag, callback)
}

func (d *MetricRepositoryDecorator) DeleteByID(ctx context.Context, userID *entities.UserID) (err error) {
	callback := func(ctx context.Context) (err error) { return d.UserRepository.DeleteByID(ctx, userID) }
	return decorators.CollectMetricForDatabaseCallback(ctx, prometheus.DeleteOperationTag, callback)
}
