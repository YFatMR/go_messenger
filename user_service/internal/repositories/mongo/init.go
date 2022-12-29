package mongo

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomMongoSettings struct {
	uri               string
	databaseName      string
	collectionName    string
	connectionTimeout time.Duration
}

func NewMongoSettings(mongoURI string, databaseName string, collectionName string,
	connectionTimeout time.Duration,
) *CustomMongoSettings {
	return &CustomMongoSettings{
		uri:               mongoURI,
		databaseName:      databaseName,
		collectionName:    collectionName,
		connectionTimeout: connectionTimeout,
	}
}

func NewMongoCollection(ctx context.Context, settings *CustomMongoSettings,
	logger *loggers.OtelZapLoggerWithTraceID,
) (*mongo.Collection, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, settings.connectionTimeout)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(settings.uri))
	if err != nil {
		panic(err)
	}

	logger.Info("Starting Ping mongodb")
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	logger.Info("mongodb Ping successfully finished")

	collection := client.Database(settings.databaseName).Collection(settings.collectionName)
	return collection, cancel
}
