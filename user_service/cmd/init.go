package main

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"github.com/YFatMR/go_messenger/core/pkg/mongodb"
	"github.com/YFatMR/go_messenger/user_service/apientity"
	"github.com/YFatMR/go_messenger/user_service/decorator"
	"github.com/YFatMR/go_messenger/user_service/user"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func UserRepositoryFromConfig(ctx context.Context, config *cviper.CustomViper, logger *czap.Logger,
	tracer trace.Tracer,
) (
	apientity.UserRepository, error,
) {
	databaseSettings := cviper.NewDatabaseSettingsFromConfig(config)
	mongoCollection, err := mongodb.Connect(ctx, databaseSettings, logger)
	if err != nil {
		logger.Error("Can't establish connection with database", zap.Error(err))
		return nil, err
	}

	collectDatabaseQueryMetrics := config.GetBoolRequired("ENABLE_DATABASE_QUERY_METRICS")
	traceDatabaseQuery := config.GetBoolRequired("ENABLE_DATABASE_QUERY_TRACING")

	repository := user.NewMongoRepository(mongoCollection, databaseSettings.GetOperationTimeout(), logger)
	repository = decorator.NewLoggingUserRepositoryDecorator(repository, logger)
	if collectDatabaseQueryMetrics {
		repository = decorator.NewPrometheusMetricsUserRepositoryDecorator(repository)
	}
	if traceDatabaseQuery {
		// TODO: make as config option (?)
		recordTraceErrors := true
		repository = decorator.NewOpentelemetryTracingUserRepositoryDecorator(
			repository, tracer, recordTraceErrors,
		)
	}
	return repository, nil
}

func UserServiceFromConfig(config *cviper.CustomViper, logger *czap.Logger, repository apientity.UserRepository,
	tracer trace.Tracer,
) apientity.UserService {
	passwordManager := user.DefaultPasswordManager()
	jwtManager := jwtmanager.FromConfig(config, logger)
	service := user.NewService(repository, passwordManager, jwtManager, logger)
	// TODO: make config options
	if true {
		service = decorator.NewLoggingUserServiceDecorator(service, logger)
	}
	// TODO: make config options
	if true {
		service = decorator.NewPrometheusMetricsUserServiceDecorator(service)
	}
	// TODO: make config options
	if true {
		recordTraceErrors := true
		service = decorator.NewOpentelemetryTracingUserServiceDecorator(service, tracer, recordTraceErrors)
	}
	return service
}

func UserControllerFromService(service apientity.UserService, logger *czap.Logger) apientity.UserController {
	controller := user.NewController(service, logger)
	// TODO: make config options
	if true {
		controller = decorator.NewLoggingUserControllerDecorator(controller, logger)
	}
	return controller
}
