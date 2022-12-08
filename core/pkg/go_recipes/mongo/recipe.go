package mongo

import (
	"context"
	"fmt"
	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(dockerDeletionTimeoutSeconds uint) *mongo.Client {
	mongoUsername := "root"
	mongoPassword := "password"
	mongoDockerTag := "6.0"
	mongoEndpointUrl := "localhost"

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
	return client
}
