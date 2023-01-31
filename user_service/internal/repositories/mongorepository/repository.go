package mongorepository

import (
	"context"
	"errors"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
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
	*userid.Entity, logerr.Error,
) {
	accountMongoID, err := primitive.ObjectIDFromHex(accountID.GetID())
	if err != nil {
		return nil, logerr.NewError(ErrUserCreation, "Got wrong id format", logerr.Err(err))
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	insertResult, err := r.collection.InsertOne(mongoOperationCtx, userDocument{
		Name:      user.GetName(),
		Surname:   user.GetSurname(),
		AccountID: accountMongoID,
	})
	if err != nil {
		return nil, logerr.NewError(ErrUserCreation, "Can't insert new user", logerr.Err(err))
	}

	userID := userid.New(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return userID, nil
}

func (r *UserMongoRepository) GetByID(ctx context.Context, userID *userid.Entity) (*user.Entity, logerr.Error) {
	objectID, err := primitive.ObjectIDFromHex(userID.GetID())
	if err != nil {
		return nil, logerr.NewError(ErrGetUser, "Got wrong id format", logerr.Err(err))
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document userDocument
	err = r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "_id", Value: objectID},
	}).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, logerr.NewError(ErrUserNotFound, "User not found by id", logerr.Err(err))
	} else if err != nil {
		return nil, logerr.NewError(ErrUserNotFound, "Database connection error", logerr.Err(err))
	}
	user := user.New(document.Name, document.Surname)
	return user, nil
}

func (r *UserMongoRepository) DeleteByID(ctx context.Context, userID *userid.Entity) logerr.Error {
	objectID, err := primitive.ObjectIDFromHex(userID.GetID())
	if err != nil {
		return logerr.NewError(ErrUserDeletion, "Got wrong id format", logerr.Err(err))
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	deleteResult, err := r.collection.DeleteOne(mongoOperationCtx, bson.M{"_id": objectID})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return logerr.NewError(ErrUserNotFound, "User not found by id", logerr.Err(err))
	} else if err != nil {
		return logerr.NewError(ErrUserNotFound, "Database connection error", logerr.Err(err))
	}
	r.logger.InfoContextNoExport(
		ctx,
		"Deleted user", zap.String("id", userID.GetID()),
		zap.Int64("deleted count", deleteResult.DeletedCount),
	)

	return nil
}
