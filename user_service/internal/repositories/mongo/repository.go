package mongo

import (
	"context"
	. "core/pkg/loggers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"user_server/internal/enities"
	"user_server/internal/repositories"
)

type UserMongoRepository struct {
	collection *mongo.Collection
	logger     *OtelZapLoggerWithTraceID
}

func NewUserMongoRepository(collection *mongo.Collection, logger *OtelZapLoggerWithTraceID) *UserMongoRepository {
	return &UserMongoRepository{
		collection: collection,
		logger:     logger,
	}
}

type userDocument struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Surname string             `bson:"surname"`
}

func (r *UserMongoRepository) Create(ctx context.Context, request *enities.User) (string, error) {
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

func (r *UserMongoRepository) GetById(ctx context.Context, id string) (*enities.User, error) {
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
