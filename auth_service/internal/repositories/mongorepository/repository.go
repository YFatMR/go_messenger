package mongorepository

import (
	"context"
	"errors"
	"time"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/metrics/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
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
	tracer           trace.Tracer
}

func NewAccountMongoRepository(collection *mongo.Collection, operationTimeout time.Duration,
	logger *loggers.OtelZapLoggerWithTraceID, tracer trace.Tracer,
) *AccountMongoRepository {
	return &AccountMongoRepository{
		collection:       collection,
		operationTimeout: operationTimeout,
		logger:           logger,
		tracer:           tracer,
	}
}

func (r *AccountMongoRepository) CreateAccount(ctx context.Context, credential *entities.Credential,
	userRole entities.Role) (
	_ *entities.AccountID, err error,
) {
	// metrics
	startTime := time.Now()
	defer prometheus.CollectDatabaseQueryMetrics(startTime, prometheus.InsertOperationTag, err)

	// process database insertion
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

	accountID := entities.NewAccountIDFromRawDatabaseDocument(insertResult.InsertedID.(primitive.ObjectID).Hex())
	r.logger.DebugContextNoExport(ctx, "Account created", zap.String("id", accountID.GetID()))
	return accountID, nil
}

func (r *AccountMongoRepository) GetTokenPayloadWithHashedPasswordByLogin(ctx context.Context, login string) (
	_ *entities.TokenPayload, _ string, err error,
) {
	// metrics
	startTime := time.Now()
	defer prometheus.CollectDatabaseQueryMetrics(startTime, prometheus.FindOperationTag, err)

	// process database search
	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document accountDocument
	err = r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "login", Value: login},
	}).Decode(&document)
	if err == nil {
		tokenPayload := entities.NewTokenPayloadFromRawDatabaseDocument(document.ID.Hex(), document.UserRole)
		hashedPassword := document.HashedPassword
		return tokenPayload, hashedPassword, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		r.logger.DebugContextNoExport(ctx, "User credential not found (by login)", zap.String("login", login))
		return nil, "", ErrAccountNotFound
	}
	r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
	return nil, "", err
}
