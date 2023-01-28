package mongodb

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func Connect(ctx context.Context, settings *cviper.DatabaseSettings, logger *loggers.OtelZapLoggerWithTraceID) (
	_ *mongo.Collection, err error,
) {
	getCollection := func() (*mongo.Collection, error) {
		client, err := mongo.Connect(
			ctx,
			options.Client().ApplyURI(settings.GetURI()),
			options.Client().SetConnectTimeout(settings.GetConnectionTimeout()),
		)
		if err != nil {
			return nil, err
		}

		logger.Info("Starting Ping mongodb")
		err = client.Ping(ctx, nil)
		if err != nil {
			logger.Info("mongodb Ping failed", zap.Error(err))
			return nil, err
		}

		logger.Info("mongodb Ping successfully finished")
		collection := client.Database(settings.GetDatabaseName()).Collection(settings.GetCollectionName())
		return collection, nil
	}

	logger.Info("Connecting to database...")
	for i := 0; i < settings.GetStartupReconnectionCount(); i++ {
		collection, err := getCollection()
		if err == nil {
			logger.Info("Successfully connected to database")
			return collection, nil
		}
		time.Sleep(settings.GetSturtupReconnectionInterval())
	}
	return nil, err
}
