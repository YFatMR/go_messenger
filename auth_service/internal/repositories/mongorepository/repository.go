package mongorepository

import (
	"context"
	"errors"
	"time"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"
	"github.com/YFatMR/go_messenger/auth_service/internal/repositories"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
}

func New(collection *mongo.Collection, operationTimeout time.Duration) *AccountMongoRepository {
	if collection == nil {
		panic("Got empty collection")
	}
	return &AccountMongoRepository{
		collection:       collection,
		operationTimeout: operationTimeout,
	}
}

func (r *AccountMongoRepository) CreateAccount(ctx context.Context, credential *credential.Entity,
	userRole entities.Role) (
	*accountid.Entity, ulo.LogStash, error,
) {
	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	insertResult, err := r.collection.InsertOne(mongoOperationCtx, accountDocument{
		Login:          credential.GetLogin(),
		HashedPassword: credential.GetHashedPassword(),
		UserRole:       userRole,
	})
	if err != nil {
		return nil, ulo.FromErrorWithMsg("can't create account", err), repositories.ErrAccountCreation
	}

	accountID := accountid.New(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return accountID, nil, nil
}

func (r *AccountMongoRepository) GetTokenPayloadWithHashedPasswordByLogin(ctx context.Context, login string) (
	*tokenpayload.Entity, string, ulo.LogStash, error,
) {
	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document accountDocument
	err := r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "login", Value: login},
	}).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		logStash := ulo.FromErrorWithMsg("user credential not found", err, ulo.String("login", login))
		return nil, "", logStash, repositories.ErrAccountNotFound
	} else if err != nil {
		return nil, "", ulo.FromErrorWithMsg("database connection error", err), repositories.ErrGetToken
	}
	tokenPayload := tokenpayload.New(document.ID.Hex(), document.UserRole)
	hashedPassword := document.HashedPassword
	return tokenPayload, hashedPassword, nil, nil
}
