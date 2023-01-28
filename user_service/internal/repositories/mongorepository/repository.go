package mongorepository

import (
	"context"
	"errors"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
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

func (r *UserMongoRepository) Create(ctx context.Context, user *entities.User, accountID *entities.AccountID) (
	_ *entities.UserID, err error,
) {
	accountMongoID, err := primitive.ObjectIDFromHex(accountID.GetID())
	if err != nil {
		r.logger.ErrorContextNoExport(ctx, "Got wrong id format", zap.Error(err))
		return nil, ErrUserCreation
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	insertResult, err := r.collection.InsertOne(mongoOperationCtx, userDocument{
		Name:      user.GetName(),
		Surname:   user.GetSurname(),
		AccountID: accountMongoID,
	})
	if err != nil {
		r.logger.ErrorContext(ctx, "Can't insert new user", zap.Error(err))
		return nil, ErrUserCreation
	}

	userID := entities.NewUserID(insertResult.InsertedID.(primitive.ObjectID).Hex())
	r.logger.DebugContextNoExport(ctx, "User id response created successfully", zap.String("id", userID.GetID()))
	return userID, nil
}

func (r *UserMongoRepository) GetByID(ctx context.Context, userID *entities.UserID) (_ *entities.User, err error) {
	objectID, err := primitive.ObjectIDFromHex(userID.GetID())
	if err != nil {
		r.logger.ErrorContextNoExport(ctx, "Got wrong id format", zap.Error(err))
		return nil, ErrGetUser
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document userDocument
	err = r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "_id", Value: objectID},
	}).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		r.logger.DebugContextNoExport(ctx, "User not found (by id)", zap.String("id", userID.GetID()))
		return nil, ErrUserNotFound
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
		return nil, ErrInternalDatabase
	}
	r.logger.DebugContextNoExport(ctx, "user found", zap.String("id", userID.GetID()))

	user := entities.NewUser(document.Name, document.Surname)
	return user, nil
}

func (r *UserMongoRepository) DeleteByID(ctx context.Context, userID *entities.UserID) (err error) {
	objectID, err := primitive.ObjectIDFromHex(userID.GetID())
	if err != nil {
		r.logger.ErrorContextNoExport(ctx, "Got wrong id format", zap.Error(err))
		return ErrUserDeletion
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	deleteResult, err := r.collection.DeleteOne(mongoOperationCtx, bson.M{"_id": objectID})
	if errors.Is(err, mongo.ErrNoDocuments) {
		r.logger.InfoContextNoExport(ctx, "User not found (by id)", zap.String("id", userID.GetID()))
		return ErrUserNotFound
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
		return ErrInternalDatabase
	}
	r.logger.InfoContextNoExport(
		ctx,
		"Deleted user", zap.String("id", userID.GetID()),
		zap.Int64("deleted count", deleteResult.DeletedCount),
	)

	return nil
}
