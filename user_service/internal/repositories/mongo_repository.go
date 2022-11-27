package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"user_server/internal/enities"
)

type UserMongoRepository struct {
	collection *mongo.Collection
}

func NewUserMongoRepository(collection *mongo.Collection) *UserMongoRepository {
	return &UserMongoRepository{
		collection: collection,
	}
}

type userDocument struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Surname string             `bson:"surname"`
}

// CRUD - 1 request - 1 methods
func (r *UserMongoRepository) Create(ctx context.Context, request *enities.User) (string, error) {
	insertResult, err := r.collection.InsertOne(ctx, userDocument{
		Name:    request.GetName(),
		Surname: request.GetSurname(),
	})
	fmt.Println("insert result", insertResult)
	return insertResult.InsertedID.(primitive.ObjectID).Hex(), err
}

func (r *UserMongoRepository) GetById(ctx context.Context, id string) (*enities.User, error) {
	var document userDocument
	fmt.Println("call repository GetById")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// TODO: not panic ->
		panic(err)
	}
	err = r.collection.FindOne(ctx, bson.D{
		{"_id", objectId},
	}).Decode(&document)
	if err == nil { // TODO: handle err == mongo.ErrNoDocuments
		return enities.NewUser(document.Name, document.Surname), nil
	} else if err == mongo.ErrNoDocuments {
		fmt.Println("Unexist user with id", id)
		return nil, UserNotFoundErr
	}
	return nil, err
}
