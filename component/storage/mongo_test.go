package storage

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ihezebin/go-template-ddd/domain/entity"
)

func TestMongo(t *testing.T) {
	ctx := context.Background()
	err := InitMongoClient(ctx, "mongodb://root:root@localhost:27017/go-template-ddd?authSource=admin")
	if err != nil {
		t.Fatal(err)
	}

	collection := MongoDatabase().Collection("example")
	_, err = collection.InsertOne(ctx, &entity.Example{
		Id:       primitive.NewObjectID().Hex(),
		Username: "admin",
		Password: "123456",
		Email:    "6wqz8@example.com",
		Salt:     "123456",
	})
	if err != nil {
		t.Fatal(err)
	}

	examples := make([]*entity.Example, 0)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	err = cursor.All(ctx, &examples)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(examples)
}
