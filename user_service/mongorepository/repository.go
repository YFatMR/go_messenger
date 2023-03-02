package mongorepository

import (
	"context"
	"errors"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/user_service/entities/account"
	"github.com/YFatMR/go_messenger/user_service/entities/credential"
	"github.com/YFatMR/go_messenger/user_service/entities/user"
	"github.com/YFatMR/go_messenger/user_service/entities/userid"
	"github.com/YFatMR/go_messenger/user_service/entities/userrole"
	"github.com/YFatMR/go_messenger/user_service/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type userDocument struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Login          string             `bson:"login,omitempty"`
	HashedPassword string             `bson:"hashed_password,omitempty"`
	UserRole       string             `bson:"user_role,omitempty"`
	Nickname       string             `bson:"nickname,omitempty"`
	Name           string             `bson:"name,omitempty"`
	Surname        string             `bson:"surname,omitempty"`
}

type userMongoRepository struct {
	collection       *mongo.Collection
	operationTimeout time.Duration
	logger           *czap.Logger
}

func NewUserMongoRepository(collection *mongo.Collection, operationTimeout time.Duration,
	logger *czap.Logger,
) repositories.UserRepository {
	return &userMongoRepository{
		collection:       collection,
		operationTimeout: operationTimeout,
		logger:           logger,
	}
}

func (r *userMongoRepository) Create(ctx context.Context, user *user.Entity, credential *credential.Entity) (
	*userid.Entity, error,
) {
	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	insertResult, err := r.collection.InsertOne(mongoOperationCtx, userDocument{
		Login:          credential.GetLogin(),
		HashedPassword: credential.GetHashedPassword(),
		UserRole:       credential.GetRole().GetName(),
		Nickname:       user.GetNickname(),
		Name:           user.GetName(),
		Surname:        user.GetSurname(),
	})
	if err != nil {
		r.logger.ErrorContext(ctx, "Can't insert new user", zap.Error(err))
		return nil, ErrUserCreation
	}

	userID := userid.New(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return userID, nil
}

func (r *userMongoRepository) GetByID(ctx context.Context, userID *userid.Entity) (
	*user.Entity, error,
) {
	objectID, err := primitive.ObjectIDFromHex(userID.GetID())
	if err != nil {
		r.logger.ErrorContext(ctx, "Got wrong id format", zap.Error(err))
		return nil, ErrGetUser
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document userDocument
	err = r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "_id", Value: objectID},
	}).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		r.logger.ErrorContext(ctx, "Can't insert new user", zap.Error(err))
		return nil, ErrUserNotFound
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
		return nil, ErrUserNotFound
	}
	user := user.New(document.Nickname, document.Name, document.Surname)
	return user, nil
}

func (r *userMongoRepository) DeleteByID(ctx context.Context, userID *userid.Entity) error {
	objectID, err := primitive.ObjectIDFromHex(userID.GetID())
	if err != nil {
		r.logger.ErrorContext(ctx, "Got wrong id format", zap.Error(err))
		return ErrUserDeletion
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	_, err = r.collection.DeleteOne(mongoOperationCtx, bson.M{"_id": objectID})
	if errors.Is(err, mongo.ErrNoDocuments) {
		r.logger.ErrorContext(ctx, "User not found by id", zap.Error(err))
		return ErrUserNotFound
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
		return ErrUserNotFound
	}
	return nil
}

func (r *userMongoRepository) GetAccountByLogin(ctx context.Context, login string) (
	*account.Entity, error,
) {
	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document userDocument
	err := r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "login", Value: login},
	}).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		r.logger.ErrorContext(ctx, "user credential not found", zap.Error(err), zap.String("login", login))
		return nil, ErrUserNotFound
	} else if err != nil {
		r.logger.ErrorContext(ctx, "database connection error", zap.Error(err))
		return nil, ErrGetToken
	}

	role, err := userrole.FromString(document.UserRole)
	if err != nil {
		r.logger.FatalContext(ctx, "database has incorrect user role", zap.String("login", login))
		return nil, err
	}

	result := account.New(
		document.ID.Hex(), document.Login, document.HashedPassword, role,
		document.Nickname, document.Name, document.Surname,
	)
	return result, nil
}
