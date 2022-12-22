package mongo

import (
	"context"
	. "github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/user_service/internal/enities"
	"github.com/YFatMR/go_messenger/user_service/internal/metrics/prometheus"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

type UserMongoRepository struct {
	collection *mongo.Collection
	logger     *OtelZapLoggerWithTraceID
	tracer     trace.Tracer
}

func NewUserMongoRepository(collection *mongo.Collection, logger *OtelZapLoggerWithTraceID, tracer trace.Tracer) *UserMongoRepository {
	return &UserMongoRepository{
		collection: collection,
		logger:     logger,
		tracer:     tracer,
	}
}

type userDocument struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Surname string             `bson:"surname"`
}

func (r *UserMongoRepository) Create(ctx context.Context, request *enities.User) (insertedID string, err error) {
	// metrics
	startTime := time.Now()
	defer collectDatabaseQueryMetrics(startTime, prometheus.InsertOperationTag, &err)

	// process database insertion
	insertResult, err := r.collection.InsertOne(ctx, userDocument{
		Name:    request.GetName(),
		Surname: request.GetSurname(),
	})
	if err != nil {
		return "", err
	}
	r.logger.DebugContextNoExport(ctx, "Insert result", zap.String("id", insertResult.InsertedID.(primitive.ObjectID).Hex()))
	return insertResult.InsertedID.(primitive.ObjectID).Hex(), err
}

func (r *UserMongoRepository) GetById(ctx context.Context, id string) (foundUser *enities.User, err error) { // metrics
	// metrics
	startTime := time.Now()
	defer collectDatabaseQueryMetrics(startTime, prometheus.FindOperationTag, &err)

	// process database search
	var document userDocument
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repositories.WrongUserIdFormatErr
	}
	err = r.collection.FindOne(ctx, bson.D{
		{"_id", objectId},
	}).Decode(&document)
	if err == nil { // TODO: handle err == mongo.ErrNoDocuments
		return enities.NewUser(document.Name, document.Surname), nil
	} else if err == mongo.ErrNoDocuments {
		r.logger.DebugContextNoExport(ctx, "User not found", zap.String("id", id))
		return nil, repositories.UserNotFoundErr
	}
	r.logger.ErrorContext(ctx, "Database connection error", zap.String("error", err.Error()))
	return nil, err
}
