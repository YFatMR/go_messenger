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

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
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
	*accountid.Entity, logerr.Error,
) {
	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	insertResult, err := r.collection.InsertOne(mongoOperationCtx, accountDocument{
		Login:          credential.GetLogin(),
		HashedPassword: credential.GetHashedPassword(),
		UserRole:       userRole,
	})
	if err != nil {
		return nil, logerr.NewError(repositories.ErrAccountCreation, "can't create account", logerr.Err(err))
	}

	accountID := accountid.New(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return accountID, nil
}

func (r *AccountMongoRepository) GetTokenPayloadWithHashedPasswordByLogin(ctx context.Context, login string) (
	*tokenpayload.Entity, string, logerr.Error,
) {
	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document accountDocument
	err := r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "login", Value: login},
	}).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, "", logerr.NewError(
			repositories.ErrAccountNotFound,
			"user credential not found", logerr.String("login", login), logerr.Err(err),
		)
	} else if err != nil {
		return nil, "", logerr.NewError(repositories.ErrGetToken, "database connection error", logerr.Err(err))
	}
	tokenPayload := tokenpayload.New(document.ID.Hex(), document.UserRole)
	hashedPassword := document.HashedPassword
	return tokenPayload, hashedPassword, nil
}
