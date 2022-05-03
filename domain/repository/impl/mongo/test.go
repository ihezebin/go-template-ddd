package mongo

import (
	"context"
	"github.com/whereabouts/sdk/db/mongoc"
	"github.com/whereabouts/web-template-ddd/component/storage"
	"github.com/whereabouts/web-template-ddd/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 实现TestRepository接口
type testRepository struct {
	mongoc.Model
}

func NewTestRepository() *testRepository {
	return &testRepository{Model: mongoc.NewAutoTimeModel(storage.GetMongoCli(), "ddd", "test").SetSoftDelete(true)}
}

func (repo *testRepository) Register(ctx context.Context, test *entity.Test) error {
	if test.Id.IsZero() {
		test.Id = primitive.NewObjectID()
	}
	_, err := repo.InsertOne(ctx, test)
	return err
}
