package mongo

import (
	"context"
	"fmt"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	dockertest "github.com/ory/dockertest/v3" // alias for golangci-lint
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoTestDatabase(mongoConfig *cviper.CustomViper) (*mongo.Database, func() error) {
	if err := mongoConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	mongoUsername := mongoConfig.GetStringRequired("MONGO_INITDB_ROOT_USERNAME")
	mongoPassword := mongoConfig.GetStringRequired("MONGO_INITDB_ROOT_PASSWORD")
	mongoDockerTag := mongoConfig.GetStringRequired("MONGO_DOCKER_TAG")
	mongoTestDatabaseName := mongoConfig.GetStringRequired("MONGO_TEST_DATABASE_NAME")

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
		port := resource.GetPort("27017/tcp")
		mongoURI := fmt.Sprintf("mongodb://%s:%s@localhost:%s", mongoUsername, mongoPassword, port)
		var err error
		ctx := context.Background()
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			return err
		}
		return client.Ping(ctx, nil)
	}); err != nil {
		panic(err)
	}
	return client.Database(mongoTestDatabaseName), func() error {
		return pool.Purge(resource)
	}
}
