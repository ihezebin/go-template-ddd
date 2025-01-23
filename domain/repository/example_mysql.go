package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"

	"github.com/ihezebin/go-template-ddd/component/storage"
	"github.com/ihezebin/go-template-ddd/domain/entity"
)

type exampleMysqlRepository struct {
	db *gorm.DB
}

func (e *exampleMysqlRepository) InsertOne(ctx context.Context, example *entity.Example) error {
	example.Id = primitive.NewObjectID().Hex()
	err := e.db.Create(example).Error
	if err != nil {
		return err
	}

	return nil
}

func (e *exampleMysqlRepository) FindByUsername(ctx context.Context, username string) (example *entity.Example, err error) {
	example = &entity.Example{}
	err = e.db.Where("username = ?", username).First(example).Error
	if err != nil {
		return nil, err
	}

	return example, nil
}

func (e *exampleMysqlRepository) FindByEmail(ctx context.Context, email string) (example *entity.Example, err error) {
	example = &entity.Example{}
	err = e.db.Where("email = ?", email).First(example).Error
	if err != nil {
		return nil, err
	}

	return example, nil
}

func NewExampleMysqlRepository() ExampleRepository {
	return &exampleMysqlRepository{
		db: storage.MySQLDatabase(),
	}
}

var _ ExampleRepository = (*exampleMysqlRepository)(nil)
