package main

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"github.com/YFatMR/go_messenger/core/pkg/pgxdb"
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
	operationTimeout := config.GetMillisecondsDurationRequired("DATABASE_OPERATION_TIMEOUT_MILLISECONDS")
	postgresURL := config.GetStringRequired("DATABASE_URL")
	connPool, err := pgxdb.Connect(ctx, postgresURL, logger)
	if err != nil {
		logger.Error("Failed to create start databse")
		return nil, err
	}
	collectDatabaseQueryMetrics := config.GetBoolRequired("ENABLE_DATABASE_QUERY_METRICS")
	traceDatabaseQuery := config.GetBoolRequired("ENABLE_DATABASE_QUERY_TRACING")

	_, err = connPool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			created_at TIMESTAMP DEFAULT NOW(),

			email VARCHAR(256) NOT NULL,
			hashed_password VARCHAR(256) NOT NULL,
			role VARCHAR(256) NOT NULL,
			nickname VARCHAR(128) NOT NULL,
			name VARCHAR(256) NOT NULL,
			surname VARCHAR(256) NOT NULL,

			UNIQUE (email)
		);`,
	)
	if err != nil {
		logger.Error("Failed to create database tables", zap.Error(err))
		return nil, err
	}

	repository := user.NewPosgreSQLRepository(connPool, operationTimeout, logger)
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
