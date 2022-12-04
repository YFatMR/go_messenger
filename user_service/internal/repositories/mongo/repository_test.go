package mongo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"os"
	"testing"
	"user_server/internal/enities"
)

var testDatabase *mongo.Database

func TestMain(m *testing.M) {
	mongoUsername := "root"
	mongoPassword := "password"
	mongoDockerTag := "6.0"
	mongoEndpointUrl := "localhost"
	mongoDatabaseName := "TEST_DATABASE"
	const dockerDeletionTimeoutSeconds uint = 60

	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}
	environmentVariables := []string{
		"MONGO_INITDB_ROOT_USERNAME=" + mongoUsername,
		"MONGO_INITDB_ROOT_PASSWORD=" + mongoPassword,
	}
	resource, err := pool.Run("mongo", mongoDockerTag, environmentVariables)
	if err != nil {
		panic(err)
	}

	var client *mongo.Client
	if err = pool.Retry(func() error {
		mongoUri := fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUsername, mongoPassword, mongoEndpointUrl, resource.GetPort("27017/tcp"))
		var err error
		ctx := context.Background()
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
		if err != nil {
			return err
		}
		return client.Ping(ctx, nil)
	}); err != nil {
		panic(err)
	}

	// setup global variables
	testDatabase = client.Database(mongoDatabaseName)

	// set docker deletion timeout
	err = resource.Expire(dockerDeletionTimeoutSeconds)
	if err != nil {
		panic(err)
	}

	// Run tests
	exitCode := m.Run()

	if err != nil {
		panic(err)
	}
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
	repository := NewUserMongoRepository(collection, zap.NewNop())

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
	repository := NewUserMongoRepository(collection, zap.NewNop())

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
