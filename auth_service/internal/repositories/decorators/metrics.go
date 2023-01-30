package decorators

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/repositories"
	"github.com/YFatMR/go_messenger/core/pkg/decorators"
	"github.com/YFatMR/go_messenger/core/pkg/metrics/prometheus"
)

type MetricRepositoryDecorator struct {
	repositories.AccountRepository
}

func NewMetricDecorator(repository repositories.AccountRepository) repositories.AccountRepository {
	return &MetricRepositoryDecorator{
		AccountRepository: repository,
	}
}

func (d *MetricRepositoryDecorator) CreateAccount(ctx context.Context, credential *credential.Entity,
	role entities.Role) (
	_ *accountid.Entity, err error,
) {
	callback := func(ctx context.Context) (_ *accountid.Entity, err error) {
		return d.AccountRepository.CreateAccount(ctx, credential, role)
	}
	return decorators.CollectMetricForDatabaseCallbackWithReturnType(ctx, prometheus.InsertOperationTag, callback)
}

func (d *MetricRepositoryDecorator) GetTokenPayloadWithHashedPasswordByLogin(ctx context.Context, login string) (
	_ *tokenpayload.Entity, hashedPassword string, err error,
) {
	callback := func(ctx context.Context) (_ *tokenpayload.Entity, hashedPassword string, err error) {
		return d.AccountRepository.GetTokenPayloadWithHashedPasswordByLogin(ctx, login)
	}
	return decorators.CollectMetricForDatabaseCallbackWithTwoReturnType(ctx, prometheus.FindOperationTag, callback)
}
