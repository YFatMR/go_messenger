package mongo

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	recipe "github.com/YFatMR/go_messenger/core/pkg/recipes/go/mongo"
	"github.com/YFatMR/go_messenger/user_service/internal/enities"
	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	testDatabase        *mongo.Database
	mongoConfigPathFlag string
)

func init() {
	flag.StringVar(&mongoConfigPathFlag, "mongo_config_path", "", "Path to mongodb configuration")
}

func TestMain(m *testing.M) {
	flag.Parse()

	ctx := context.Background()

	mongoConfig := cviper.New()
	mongoConfig.SetConfigFile(mongoConfigPathFlag)
	database, dockerContainerPurge := recipe.NewMongoTestDatabase(mongoConfig)

	// Setup global variable
	testDatabase = database

	// Run tests
	exitCode := m.Run()

	// Clenup database
	if err := testDatabase.Drop(ctx); err != nil {
		panic(err)
	}

	// Clenup container
	if err := dockerContainerPurge(); err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func dropCollection(ctx context.Context, t *testing.T, collection *mongo.Collection) {
	t.Helper()

	err := collection.Drop(ctx)
	if err != nil {
		t.Error(err)
	}
}

func newMockUserMongoRepository(collection *mongo.Collection) *UserMongoRepository {
	tracer := otel.Tracer("fake")
	nopLogger := loggers.NewOtelZapLoggerWithTraceID(otelzap.New(zap.NewNop()))
	return NewUserMongoRepository(collection, nopLogger, tracer)
}

func TestUserCreation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Initialize repository
	randomCollection := testDatabase.Collection(uuid.New().String())
	defer dropCollection(ctx, t, randomCollection)
	repository := newMockUserMongoRepository(randomCollection)

	// Start test
	userData := enities.NewUser("Ivan", "Petrov")
	_, err := repository.Create(context.Background(), userData)
	if err != nil {
		t.Fatalf("User creation failed with error: %s", err)
	}
}

func TestFindCreatedUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Initialize repository
	randomCollection := testDatabase.Collection(uuid.New().String())
	defer dropCollection(ctx, t, randomCollection)
	repository := newMockUserMongoRepository(randomCollection)

	// Start test
	userData := enities.NewUser("Sergey", "Satnav")
	UserID, err := repository.Create(context.Background(), userData)
	if err != nil {
		t.Fatalf("User creation failed with error: %s", err)
	}

	responseUserData, err := repository.GetByID(context.Background(), UserID)
	if err != nil {
		t.Fatalf("User search failed with error: %s", err)
	}
	if responseUserData == nil {
		t.Fatalf("User with id %s not exist", UserID)
	}

	if *userData != *responseUserData {
		t.Fatalf("Created and found different users %s %s", userData, responseUserData)
	}
}
