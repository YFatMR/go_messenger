package mongorepository

import (
	"context"
	"errors"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/userid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type userDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Surname   string             `bson:"surname,omitempty"`
	AccountID primitive.ObjectID `bson:"account_id,omitempty"`
}

type UserMongoRepository struct {
	collection       *mongo.Collection
	operationTimeout time.Duration
	logger           *loggers.OtelZapLoggerWithTraceID
}

func NewUserMongoRepository(collection *mongo.Collection, operationTimeout time.Duration,
	logger *loggers.OtelZapLoggerWithTraceID,
) *UserMongoRepository {
	return &UserMongoRepository{
		collection:       collection,
		operationTimeout: operationTimeout,
		logger:           logger,
	}

}

func (r *UserMongoRepository) Create(ctx context.Context, user *user.Entity, accountID *accountid.Entity) (
	*userid.Entity, ulo.LogStash, error,
) {
	accountMongoID, err := primitive.ObjectIDFromHex(accountID.GetID())
	if err != nil {
		return nil, ulo.ErrorMsg(ulo.Message("Got wrong id format"), ulo.Error(err)), ErrUserCreation
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	insertResult, err := r.collection.InsertOne(mongoOperationCtx, userDocument{
		Name:      user.GetName(),
		Surname:   user.GetSurname(),
		AccountID: accountMongoID,
	})
	if err != nil {
		return nil, ulo.ErrorMsg(ulo.Message("Can't insert new user"), ulo.Error(err)), err
	}

	userID := userid.New(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return userID, nil, nil
}

func (r *UserMongoRepository) GetByID(ctx context.Context, userID *userid.Entity) (
	*user.Entity, ulo.LogStash, error,
) {
	objectID, err := primitive.ObjectIDFromHex(userID.GetID())
	if err != nil {
		return nil, ulo.ErrorMsg(ulo.Message("Got wrong id format"), ulo.Error(err)), ErrGetUser
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document userDocument
	err = r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "_id", Value: objectID},
	}).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ulo.ErrorMsg(ulo.Message("User not found by id"), ulo.Error(err)), ErrUserNotFound
	} else if err != nil {
		return nil, ulo.ErrorMsg(ulo.Message("Database connection error"), ulo.Error(err)), ErrUserNotFound
	}
	user := user.New(document.Name, document.Surname)
	return user, nil, nil
}

func (r *UserMongoRepository) DeleteByID(ctx context.Context, userID *userid.Entity) (
	logstash ulo.LogStash, err error,
) {
	objectID, err := primitive.ObjectIDFromHex(userID.GetID())
	if err != nil {
		return ulo.ErrorMsg(ulo.Message("Got wrong id format"), ulo.Error(err)), ErrUserDeletion
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	deleteResult, err := r.collection.DeleteOne(mongoOperationCtx, bson.M{"_id": objectID})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ulo.ErrorMsg(ulo.Message("User not found by id"), ulo.Error(err)), ErrUserNotFound
	} else if err != nil {
		return ulo.ErrorMsg(ulo.Message("Database connection error"), ulo.Error(err)), ErrUserNotFound
	}
	r.logger.InfoContextNoExport(
		ctx,
		"Deleted user", zap.String("id", userID.GetID()),
		zap.Int64("deleted count", deleteResult.DeletedCount),
	)
	return nil, nil
}
