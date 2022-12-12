package mongo

import (
	"context"
	"core/pkg/loggers"
	recipe "core/pkg/recipes/go/mongo"
	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"testing"
	"user_server/internal/enities"
)

var testDatabase *mongo.Database

func TestMain(m *testing.M) {
	const dockerDeletionTimeoutSeconds uint = 60
	mongoDatabaseName := "TEST_DATABASE"
	mongoClient := recipe.NewMongoClient(dockerDeletionTimeoutSeconds)

	// setup global variables
	testDatabase = mongoClient.Database(mongoDatabaseName)

	// Run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}

func NewDatabaseCollection(database *mongo.Database) *mongo.Collection {
	return database.Collection(uuid.New().String())
}

func NewUserMongoCollectionWithDrop(t *testing.T, database *mongo.Database) (*mongo.Collection, func(context.Context, *mongo.Collection)) {
	return NewDatabaseCollection(database), func(ctx context.Context, collection *mongo.Collection) {
		err := collection.Drop(ctx)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestUserCreation(t *testing.T) {
	// initialize repository
	collection, dropCollection := NewUserMongoCollectionWithDrop(t, testDatabase)
	defer dropCollection(context.Background(), collection)
	repository := NewUserMongoRepository(collection, loggers.NewOtelZapLoggerWithTraceID(otelzap.New(zap.NewNop())))

	// start test
	userData := enities.NewUser("Ivan", "Petrov")
	_, err := repository.Create(context.Background(), userData)
	if err != nil {
		t.Fatalf("User creation failed with error: %s", err)
	}
}

func TestFindCreatedUser(t *testing.T) {
	// initialize repository
	collection, dropCollection := NewUserMongoCollectionWithDrop(t, testDatabase)
	defer dropCollection(context.Background(), collection)
	repository := NewUserMongoRepository(collection, loggers.NewOtelZapLoggerWithTraceID(otelzap.New(zap.NewNop())))

	// start test
	userData := enities.NewUser("Sergey", "Satnav")
	userId, err := repository.Create(context.Background(), userData)
	if err != nil {
		t.Fatalf("User creation failed with error: %s", err)
	}

	responseUserData, err := repository.GetById(context.Background(), userId)
	if err != nil {
		t.Fatalf("User search failed with error: %s", err)
	}
	if responseUserData == nil {
		t.Fatalf("User with id %s not exist", userId)
	}

	if *userData != *responseUserData {
		t.Fatalf("Created and found different users %s %s", userData, responseUserData)
	}
}
