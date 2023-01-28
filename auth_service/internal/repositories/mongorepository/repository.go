package mongorepository

import (
	"context"
	"errors"
	"time"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// TODO: Distinct Login.
type accountDocument struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Login          string             `bson:"login,omitempty"`
	HashedPassword string             `bson:"hashed_password,omitempty"`
	UserRole       entities.Role      `bson:"user_role,omitempty"`
}

type AccountMongoRepository struct {
	collection       *mongo.Collection
	operationTimeout time.Duration
	logger           *loggers.OtelZapLoggerWithTraceID
}

func NewAccountMongoRepository(collection *mongo.Collection, operationTimeout time.Duration,
	logger *loggers.OtelZapLoggerWithTraceID,
) *AccountMongoRepository {
	return &AccountMongoRepository{
		collection:       collection,
		operationTimeout: operationTimeout,
		logger:           logger,
	}
}

func (r *AccountMongoRepository) CreateAccount(ctx context.Context, credential *entities.Credential,
	userRole entities.Role) (
	_ *entities.AccountID, err error,
) {
	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	insertResult, err := r.collection.InsertOne(mongoOperationCtx, accountDocument{
		Login:          credential.GetLogin(),
		HashedPassword: credential.GetHashedPassword(),
		UserRole:       userRole,
	})
	if err != nil {
		r.logger.ErrorContext(ctx, "can't create account", zap.Error(err))
		return nil, err
	}

	accountID := entities.NewAccountID(insertResult.InsertedID.(primitive.ObjectID).Hex())
	r.logger.DebugContextNoExport(ctx, "Account created", zap.String("id", accountID.GetID()))
	return accountID, nil
}

func (r *AccountMongoRepository) GetTokenPayloadWithHashedPasswordByLogin(ctx context.Context, login string) (
	_ *entities.TokenPayload, _ string, err error,
) {
	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document accountDocument
	err = r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "login", Value: login},
	}).Decode(&document)
	if err == nil {
		tokenPayload := entities.NewTokenPayload(document.ID.Hex(), document.UserRole)
		hashedPassword := document.HashedPassword
		return tokenPayload, hashedPassword, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		r.logger.DebugContextNoExport(ctx, "User credential not found (by login)", zap.String("login", login))
		return nil, "", ErrAccountNotFound
	}
	r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
	return nil, "", err
}
