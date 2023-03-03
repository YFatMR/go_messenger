package user_test

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	recipe "github.com/YFatMR/go_messenger/core/pkg/recipes/go/mongo"
	"github.com/YFatMR/go_messenger/user_service/apientity"
	"github.com/YFatMR/go_messenger/user_service/entity"
	"github.com/YFatMR/go_messenger/user_service/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var mongoConfigPathFlag string

func init() {
	flag.StringVar(&mongoConfigPathFlag, "mongo_config_path", "", "Path to mongodb configuration")
}

func newMockUserMongoRepository(collection *mongo.Collection) apientity.UserRepository {
	nopLogger := czap.New(
		*otelzap.New(zap.NewNop()),
		czap.Settings{
			LogTraceID: false,
		},
	)
	databaseOperationTimeout := time.Millisecond * 800
	return user.NewMongoRepository(collection, databaseOperationTimeout, nopLogger)
}

type MongoRepositoryTestSuite struct {
	database *mongo.Database
	suite.Suite
}

func TestMongoRepositoryTestSuite(t *testing.T) {
	flag.Parse()

	ctx := context.Background()

	mongoConfig := cviper.New()
	mongoConfig.SetConfigFile(mongoConfigPathFlag)
	if err := mongoConfig.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		panic(err)
	}

	database, dockerContainerPurge := recipe.NewMongoTestDatabase(mongoConfig)

	suite.Run(t, &MongoRepositoryTestSuite{
		database: database,
	})

	if err := database.Drop(ctx); err != nil {
		panic(err)
	}

	// Clenup container
	if err := dockerContainerPurge(); err != nil {
		panic(err)
	}
}

func (s *MongoRepositoryTestSuite) TestUserCreation() {
	ctx := context.Background()
	require := s.Require()

	// Initialize repository
	randomCollection := s.database.Collection(uuid.NewString())
	defer require.NoError(randomCollection.Drop(ctx))
	repository := newMockUserMongoRepository(randomCollection)

	// Start test
	userData := &entity.User{Nickname: "nick", Name: "Ivan", Surname: "Petrov"}
	credential := &entity.Credential{Login: "login", HashedPassword: "hashedPassword", Role: &entity.RoleUser}
	userID, err := repository.Create(context.Background(), userData, credential)
	require.NoError(err)
	require.NotNil(userID)
}

func (s *MongoRepositoryTestSuite) TestFindCreatedUser() {
	ctx := context.Background()
	require := s.Require()

	// Initialize repository
	randomCollection := s.database.Collection(uuid.NewString())
	defer require.NoError(randomCollection.Drop(ctx))
	repository := newMockUserMongoRepository(randomCollection)

	// Start test
	userData := &entity.User{Nickname: "nick", Name: "Ivan", Surname: "Petrov"}
	credential := &entity.Credential{Login: "login", HashedPassword: "hashedPassword", Role: &entity.RoleUser}
	userID, err := repository.Create(context.Background(), userData, credential)
	require.NoError(err, "Can't create user")
	require.NotNil(userID)

	responseUserData, err := repository.GetByID(context.Background(), userID)
	require.NoError(err)
	require.NotNil(responseUserData)

	usersSame := userData.Name == responseUserData.Name &&
		userData.Surname == responseUserData.Surname
	require.True(usersSame, "Created and found users are different")
}
