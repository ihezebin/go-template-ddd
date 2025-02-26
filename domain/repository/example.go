package repository

import (
	"context"

	"github.com/ihezebin/go-template-ddd/domain/entity"
)

type ExampleRepository interface {
	InsertOne(ctx context.Context, example *entity.Example) error
	FindByUsername(ctx context.Context, username string) (example *entity.Example, err error)
	FindByEmail(ctx context.Context, email string) (example *entity.Example, err error)
}

var exampleRepo ExampleRepository

func GetExampleRepository() ExampleRepository {
	return exampleRepo
}

func SetExampleRepository(repo ExampleRepository) {
	exampleRepo = repo
}

var exampleEsRepo ExampleRepository

func GetExampleEsRepository() ExampleRepository {
	return exampleEsRepo
}

func SetExampleEsRepository(repo ExampleRepository) {
	exampleEsRepo = repo
}
