package mongo

import (
	"context"
	"github.com/whereabouts/web-template-ddd/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 实现TestRepository接口
type testRepository struct {
	base *Base
}

func NewTestRepository() *testRepository {
	return &testRepository{base: NewMongoBase("ddd", "test").SetSoftDelete(true)}
}

func (repo *testRepository) Register(ctx context.Context, test *entity.Test) error {
	if test.Id.IsZero() {
		test.Id = primitive.NewObjectID()
	}
	_, err := repo.base.InsertOne(ctx, test)
	return err
}
