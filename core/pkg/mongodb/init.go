package mongodb

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoSettings struct {
	mongoURI          string
	databaseName      string
	collectionName    string
	connectionTimeout time.Duration
	logger            *loggers.OtelZapLoggerWithTraceID
}

func NewMongoSettings(mongoURI string, databaseName string, collectionName string,
	connectionTimeout time.Duration, logger *loggers.OtelZapLoggerWithTraceID,
) *MongoSettings {
	return &MongoSettings{
		mongoURI:          mongoURI,
		databaseName:      databaseName,
		collectionName:    collectionName,
		connectionTimeout: connectionTimeout,
		logger:            logger,
	}
}

func Connect(ctx context.Context, reconnectionCount int, reconnectInterval time.Duration, settings *MongoSettings) (
	_ *mongo.Collection, err error,
) {
	getCollection := func() (*mongo.Collection, error) {
		client, err := mongo.Connect(
			ctx,
			options.Client().ApplyURI(settings.mongoURI),
			options.Client().SetConnectTimeout(settings.connectionTimeout),
		)
		if err != nil {
			return nil, err
		}

		settings.logger.Info("Starting Ping mongodb")
		err = client.Ping(ctx, nil)
		if err != nil {
			settings.logger.Info("mongodb Ping failed", zap.Error(err))
			return nil, err
		}

		settings.logger.Info("mongodb Ping successfully finished")
		collection := client.Database(settings.databaseName).Collection(settings.collectionName)
		return collection, nil
	}

	for i := 0; i < reconnectionCount; i++ {
		collection, err := getCollection()
		if err == nil {
			return collection, nil
		}
		time.Sleep(reconnectInterval)
	}
	return nil, err
}
