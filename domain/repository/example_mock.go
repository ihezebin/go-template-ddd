package repository

import (
	"context"

	"github.com/ihezebin/go-template-ddd/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type exampleMockRepository struct {
}

func NewExampleMockRepository() ExampleRepository {
	return &exampleMockRepository{}
}

var _ ExampleRepository = (*exampleMockRepository)(nil)

var examples = make([]*entity.Example, 0)

func (e *exampleMockRepository) InsertOne(ctx context.Context, example *entity.Example) error {
	example.Id = primitive.NewObjectID().Hex()
	examples = append(examples, example)
	return nil
}

func (e *exampleMockRepository) FindByUsername(ctx context.Context, username string) (example *entity.Example, err error) {
	for _, v := range examples {
		if v.Username == username {
			return v, nil
		}
	}
	return nil, nil
}

func (e *exampleMockRepository) FindByEmail(ctx context.Context, email string) (example *entity.Example, err error) {
	for _, v := range examples {
		if v.Email == email {
			return v, nil
		}
	}
	return nil, nil
}
