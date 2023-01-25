package decorators

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/decorators"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories"
	"go.opentelemetry.io/otel/trace"
)

type TracingRepositoryDecorator struct {
	repositories.UserRepository
	tracer       trace.Tracer
	recordErrors bool
}

func NewTracingRepositoryDecorator(repository repositories.UserRepository, tracer trace.Tracer, recordErrors bool,
) repositories.UserRepository {
	return &TracingRepositoryDecorator{
		UserRepository: repository,
		tracer:         tracer,
		recordErrors:   recordErrors,
	}
}

func (d *TracingRepositoryDecorator) Create(ctx context.Context, user *entities.User, accountID *entities.AccountID) (
	_ *entities.UserID, err error,
) {
	callback := func(ctx context.Context) (_ *entities.UserID, err error) {
		return d.UserRepository.Create(ctx, user, accountID)
	}
	return decorators.TraceCallbackWithReturnType(
		ctx, d.tracer, "repositories.UserRepository_Create", d.recordErrors, callback,
	)
}

func (d *TracingRepositoryDecorator) GetByID(ctx context.Context, userID *entities.UserID) (
	_ *entities.User, err error,
) {
	callback := func(ctx context.Context) (_ *entities.User, err error) {
		return d.UserRepository.GetByID(ctx, userID)
	}
	return decorators.TraceCallbackWithReturnType(
		ctx, d.tracer, "repositories.UserRepository_GetByID", d.recordErrors, callback,
	)
}

func (d *TracingRepositoryDecorator) DeleteByID(ctx context.Context, userID *entities.UserID) (err error) {
	callback := func(ctx context.Context) (err error) { return d.UserRepository.DeleteByID(ctx, userID) }
	return decorators.TraceCallback(ctx, d.tracer, "repositories.UserRepository_DeleteByID", d.recordErrors, callback)
}
