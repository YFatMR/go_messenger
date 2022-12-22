package services

import (
	"context"
	. "github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/user_service/internal/enities"
	"github.com/YFatMR/go_messenger/user_service/internal/metrics/prometheus"
	"go.opentelemetry.io/otel/trace"
)

type userRepository interface {
	Create(ctx context.Context, request *enities.User) (string, error)
	GetById(ctx context.Context, id string) (*enities.User, error)
}

type UserService struct {
	repository userRepository
	logger     *OtelZapLoggerWithTraceID
	tracer     trace.Tracer
}

func NewUserService(repository userRepository, logger *OtelZapLoggerWithTraceID, tracer trace.Tracer) *UserService {
	return &UserService{
		repository: repository,
		logger:     logger,
		tracer:     tracer,
	}
}

func (s *UserService) Create(ctx context.Context, request *enities.User) (string, error) {
	const endpointTag = "CreateUser"

	prometheus.RequestProcessingTotal.WithLabelValues(endpointTag).Inc()
	userID, err := s.repository.Create(ctx, request)
	if err != nil {
		prometheus.RequestProcessingErrorsTotal.WithLabelValues(endpointTag, prometheus.ServerSideErrorRequestTag).Inc()
	}
	return userID, err
}

func (s *UserService) GetById(ctx context.Context, id string) (*enities.User, error) {
	const endpointTag = "GetUserByID"

	prometheus.RequestProcessingTotal.WithLabelValues(endpointTag).Inc()
	userEntity, err := s.repository.GetById(ctx, id)
	if err != nil {
		prometheus.RequestProcessingErrorsTotal.WithLabelValues(endpointTag, prometheus.ServerSideErrorRequestTag).Inc()
	}
	return userEntity, err
}
