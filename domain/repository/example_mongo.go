package repository

import (
	"context"

	"github.com/ihezebin/go-template-ddd/component/storage"
	"github.com/ihezebin/go-template-ddd/domain/entity"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type exampleMongoRepository struct {
	coll *mongo.Collection
}

const exampleCollectionName = "example"

func NewExampleMongoRepository() ExampleRepository {
	return &exampleMongoRepository{
		coll: storage.MongoStorageDatabase().Collection(exampleCollectionName),
	}
}

var _ ExampleRepository = (*exampleMongoRepository)(nil)

func (e *exampleMongoRepository) InsertOne(ctx context.Context, example *entity.Example) error {
	example.Id = primitive.NewObjectID().Hex()
	_, err := e.coll.InsertOne(ctx, example)
	if err != nil {
		return err
	}

	return nil
}

func (e *exampleMongoRepository) FindByUsername(ctx context.Context, username string) (example *entity.Example, err error) {
	example = &entity.Example{}
	err = e.coll.FindOne(ctx, bson.M{"username": username}).Decode(example)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return example, nil
}

func (e *exampleMongoRepository) FindByEmail(ctx context.Context, email string) (example *entity.Example, err error) {
	example = &entity.Example{}
	err = e.coll.FindOne(ctx, bson.M{"email": email}).Decode(example)
	if err != nil {
		return nil, err
	}

	return example, nil
}
