package mongorepository

import (
	"context"
	"errors"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	accountid "github.com/YFatMR/go_messenger/user_service/internal/entities/account_id"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	userid "github.com/YFatMR/go_messenger/user_service/internal/entities/user_id"
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
	*userid.Entity, cerrors.Error,
) {
	accountMongoID, err := primitive.ObjectIDFromHex(accountID.GetID())
	if err != nil {
		return nil, cerrors.New("Got wrong id format", err, ErrUserCreation)
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	insertResult, err := r.collection.InsertOne(mongoOperationCtx, userDocument{
		Name:      user.GetName(),
		Surname:   user.GetSurname(),
		AccountID: accountMongoID,
	})
	if err != nil {
		return nil, cerrors.New("Can't insert new user", err, ErrUserCreation)
	}

	userID := userid.New(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return userID, nil
}

func (r *UserMongoRepository) GetByID(ctx context.Context, userID *userid.Entity) (*user.Entity, cerrors.Error) {
	objectID, err := primitive.ObjectIDFromHex(userID.GetID())
	if err != nil {
		return nil, cerrors.New("Got wrong id format", err, ErrUserCreation)
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document userDocument
	err = r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "_id", Value: objectID},
	}).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, cerrors.New("User not found by id", err, ErrUserNotFound)
	} else if err != nil {
		return nil, cerrors.New("Database connection error", err, ErrUserNotFound)
	}
	user := user.New(document.Name, document.Surname)
	return user, nil
}

func (r *UserMongoRepository) DeleteByID(ctx context.Context, userID *userid.Entity) cerrors.Error {
	objectID, err := primitive.ObjectIDFromHex(userID.GetID())
	if err != nil {
		return cerrors.New("Got wrong id format", err, ErrUserDeletion)
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	deleteResult, err := r.collection.DeleteOne(mongoOperationCtx, bson.M{"_id": objectID})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return cerrors.New("User not found by id", err, ErrUserNotFound)
	} else if err != nil {
		return cerrors.New("Database connection error", err, ErrUserNotFound)
	}
	r.logger.InfoContextNoExport(
		ctx,
		"Deleted user", zap.String("id", userID.GetID()),
		zap.Int64("deleted count", deleteResult.DeletedCount),
	)

	return nil
}
