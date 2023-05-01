package main

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/grpcclients"
	"github.com/YFatMR/go_messenger/core/pkg/pgxdb"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/dialog_service/decorator"
	"github.com/YFatMR/go_messenger/dialog_service/dialog"
	"github.com/YFatMR/go_messenger/dialog_service/grpcapi"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func DialogRepositoryFromConfig(ctx context.Context, config *cviper.CustomViper, logger *czap.Logger) (
	apientity.DialogRepository, error,
) {
	postgresURL := config.GetStringRequired("DATABASE_URL")
	connPool, err := pgxdb.Connect(ctx, postgresURL, logger)
	if err != nil {
		logger.Error("Failed to create start databse")
		return nil, err
	}

	_, err = connPool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS dialogs (
			id BIGSERIAL PRIMARY KEY,
			created_at TIMESTAMP DEFAULT NOW(),

			user_id_1 BIGINT NOT NULL,
			dialog_name_1 VARCHAR(512) NOT NULL,
			unread_messages_count_1 BIGINT DEFAULT 0,

			user_id_2 BIGINT NOT NULL,
			dialog_name_2 VARCHAR(512) NOT NULL,
			unread_messages_count_2 BIGINT DEFAULT 0,

			UNIQUE (user_id_1, user_id_2)
		);`,
	)
	if err != nil {
		logger.Error("Failed to create database dialogs tables", zap.Error(err))
		return nil, err
	}

	_, err = connPool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS messages (
			id BIGSERIAL PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),

			dialog_id BIGINT NOT NULL,
			FOREIGN KEY (dialog_id) REFERENCES dialogs (id),
			sender_id BIGINT NOT NULL,
			text VARCHAR(1000) NOT NULL
		);`,
	)
	if err != nil {
		logger.Error("Failed to create database messages tables", zap.Error(err))
		return nil, err
	}

	settings := dialog.DialogRepositorySettings{
		OperationTimeout: config.GetMillisecondsDurationRequired("DATABASE_OPERATION_TIMEOUT_MILLISECONDS"),
	}

	repository := dialog.NewPosgreRepository(settings, connPool, logger)
	// Decorators
	repository = decorator.NewLoggingDialogRepositoryDecorator(repository, logger)
	return repository, nil
}

func DialogModelFromConfig(ctx context.Context, repository apientity.DialogRepository, config *cviper.CustomViper,
	logger *czap.Logger,
) (
	apientity.DialogModel, error,
) {
	userServiceAddress := config.GetStringRequired("USER_SERVICE_ADDRESS")
	connTimeout := config.GetMillisecondsDurationRequired("USER_SERVICE_CONNECTION_TIMEOUT_MILLISECONDS")

	grpcBackoffConfig := backoff.Config{
		BaseDelay:  config.GetMillisecondsDurationRequired("GRPC_CONNECTION_BACKOFF_DELAY_MILLISECONDS"),
		Multiplier: config.GetFloat64Required("GRPC_CONNECTION_BACKOFF_MULTIPLIER"),
		Jitter:     config.GetFloat64Required("GRPC_CONNECTION_BACKOFF_JITTER"),
		MaxDelay:   config.GetMillisecondsDurationRequired("GRPC_CONNECTION_BACKOFF_MAX_DELAY_MILLISECONDS"),
	}

	grpcKeepaliveParameters := keepalive.ClientParameters{
		Time:                config.GetMillisecondsDurationRequired("GRPC_CONNECTION_KEEPALIVE_TIME_MILLISECONDS"),
		Timeout:             config.GetMillisecondsDurationRequired("GRPC_CONNECTION_KEEPALIVE_TIMEOUT_MILLISECONDS"),
		PermitWithoutStream: config.GetBoolRequired("GRPC_CONNECTION_KEEPALIVE_PERMIT_WITHOUT_STREAM"),
	}

	userClientOpts := []grpc.DialOption{
		grpc.WithKeepaliveParams(grpcKeepaliveParameters),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				Backoff: grpcBackoffConfig,
			},
		),
	}

	userClient, err := grpcclients.NewGRPCUserClient(ctx, userServiceAddress, connTimeout, userClientOpts)
	if err != nil {
		return nil, err
	}
	model := dialog.NewDialogModel(repository, userClient, logger)
	// Decorators
	model = decorator.NewLoggingDialogModelDecorator(model, logger)
	return model, nil
}

func DialogControllerFromConfig(model apientity.DialogModel, config *cviper.CustomViper,
	logger *czap.Logger,
) (
	apientity.DialogController, error,
) {
	grpcHeaders := grpcapi.HeadersFromConfig(config)
	contextManager := grpcapi.NewContextManager(grpcHeaders)
	controller := dialog.NewController(contextManager, model, logger)
	controller = decorator.NewLoggingDialogControllerDecorator(controller, logger)
	return controller, nil
}
