package services

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/metrics/prometheus"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
	"go.opentelemetry.io/otel/trace"
)

type userRepository interface {
	Create(ctx context.Context, user *entities.User, accountID *entities.AccountID) (*entities.UserID, error)
	GetByID(ctx context.Context, userID *entities.UserID) (*entities.User, error)
	DeleteByID(ctx context.Context, userID *entities.UserID) error
}

type UserService struct {
	userRepository userRepository
	logger         *loggers.OtelZapLoggerWithTraceID
	tracer         trace.Tracer
}

func NewUserService(repository userRepository, logger *loggers.OtelZapLoggerWithTraceID, tracer trace.Tracer,
) *UserService {
	return &UserService{
		userRepository: repository,
		logger:         logger,
		tracer:         tracer,
	}
}

func (s *UserService) Create(ctx context.Context, user *entities.User, accountID *entities.AccountID) (
	_ *entities.UserID, err error,
) {
	const endpointTag = "GetAccountByToken"
	defer prometheus.CollectServiceRequestMetrics(endpointTag, err)

	return s.userRepository.Create(ctx, user, accountID)
}

func (s *UserService) GetByID(ctx context.Context, userID *entities.UserID) (_ *entities.User, err error) {
	const endpointTag = "GetUserByID"
	defer prometheus.CollectServiceRequestMetrics(endpointTag, err)

	return s.userRepository.GetByID(ctx, userID)
}

func (s *UserService) DeleteByID(ctx context.Context, userID *entities.UserID) (err error) {
	const endpointTag = "DeleteUserByID"
	defer prometheus.CollectServiceRequestMetrics(endpointTag, err)

	return s.userRepository.DeleteByID(ctx, userID)
}
