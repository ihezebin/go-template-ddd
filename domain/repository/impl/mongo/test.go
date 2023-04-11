package mongo

import (
	"context"
	"github.com/ihezebin/web-template-ddd/component/storage"
	"github.com/ihezebin/web-template-ddd/domain/entity"
	"github.com/whereabouts/sdk/db/mongoc"
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
