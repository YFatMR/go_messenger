package mongorepository

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	recipe "github.com/YFatMR/go_messenger/core/pkg/recipes/go/mongo"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.mongodb.org/mongo-driver/mongo"
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
	nopLogger := loggers.NewOtelZapLoggerWithTraceID(otelzap.New(zap.NewNop()))
	databaseOperationTimeout := time.Millisecond * 800
	return NewUserMongoRepository(collection, databaseOperationTimeout, nopLogger)
}

func TestUserCreation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Initialize repository
	randomCollection := testDatabase.Collection(uuid.NewString())
	defer dropCollection(ctx, t, randomCollection)
	repository := newMockUserMongoRepository(randomCollection)

	// Start test
	userData := entities.NewUser("Ivan", "Petrov")
	accountID := entities.NewAccountID("63c6f759bbe1022255a6b9b5")
	_, err := repository.Create(context.Background(), userData, accountID)
	if err != nil {
		t.Fatalf("User creation failed with error: %s", err)
	}
}

func TestFindCreatedUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Initialize repository
	randomCollection := testDatabase.Collection(uuid.NewString())
	defer dropCollection(ctx, t, randomCollection)
	repository := newMockUserMongoRepository(randomCollection)

	// Start test
	userData := entities.NewUser("Ivan1", "Petrov1")
	accountID := entities.NewAccountID("53c6f759bbe1022255a6b9b5")
	userID, err := repository.Create(context.Background(), userData, accountID)
	if err != nil {
		t.Fatalf("User creation failed with error: %s", err)
	}

	responseUserData, err := repository.GetByID(context.Background(), userID)
	if err != nil {
		t.Fatalf("User search failed with error: %s", err)
	}
	if responseUserData == nil {
		t.Fatalf("User with id %s not exist", userID)
	}

	if userData.GetName() != responseUserData.GetName() || userData.GetSurname() != responseUserData.GetSurname() {
		t.Fatalf("Created and found different users %s %s", userData, responseUserData)
	}
}
