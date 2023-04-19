package user_test

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	recipe "github.com/YFatMR/go_messenger/core/pkg/recipes/go/mongo"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoConfigPathFlag string

// TODO: rename mongo vars to common. Like database
func init() {
	flag.StringVar(&mongoConfigPathFlag, "mongo_config_path", "", "Path to mongodb configuration")
}

// func newMockUserRepository(collection *mongo.Collection) apientity.UserRepository {
// 	nopLogger := czap.New(
// 		*otelzap.New(zap.NewNop()),
// 		czap.Settings{
// 			LogTraceID: false,
// 		},
// 	)
// 	databaseOperationTimeout := time.Millisecond * 800
// 	// TODO change to pgx
// 	return user.NewPosgreSQLRepository(nil, databaseOperationTimeout, nopLogger)
// }

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

// func (s *MongoRepositoryTestSuite) TestUserCreation() {
// 	ctx := context.Background()
// 	require := s.Require()

// 	// Initialize repository
// 	randomCollection := s.database.Collection(uuid.NewString())
// 	defer require.NoError(randomCollection.Drop(ctx))
// 	repository := newMockUserMongoRepository(randomCollection)

// 	// Start test
// 	userData := &entity.User{Nickname: "nick", Name: "Ivan", Surname: "Petrov"}
// 	credential := &entity.Credential{Email: "email", HashedPassword: "hashedPassword", Role: entity.RoleUser}
// 	userID, err := repository.Create(context.Background(), userData, credential)
// 	require.NoError(err)
// 	require.NotNil(userID)
// }

// func (s *MongoRepositoryTestSuite) TestFindCreatedUser() {
// 	ctx := context.Background()
// 	require := s.Require()

// 	// Initialize repository
// 	randomCollection := s.database.Collection(uuid.NewString())
// 	defer require.NoError(randomCollection.Drop(ctx))
// 	repository := newMockUserMongoRepository(randomCollection)

// 	// Start test
// 	userData := &entity.User{Nickname: "nick", Name: "Ivan", Surname: "Petrov"}
// 	credential := &entity.Credential{Email: "email", HashedPassword: "hashedPassword", Role: entity.RoleUser}
// 	userID, err := repository.Create(context.Background(), userData, credential)
// 	require.NoError(err, "Can't create user")
// 	require.NotNil(userID)

// 	responseUserData, err := repository.GetByID(context.Background(), userID)
// 	require.NoError(err)
// 	require.NotNil(responseUserData)

// 	usersSame := userData.Name == responseUserData.Name &&
// 		userData.Surname == responseUserData.Surname
// 	require.True(usersSame, "Created and found users are different")
// }
