package test

import (
	"context"
	"github.com/ihezebin/go-template-ddd/component/storage"
	"github.com/ihezebin/go-template-ddd/domain/entity"
	"github.com/ihezebin/sdk/model/mongoc"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 实现Repository接口
type repoMongo struct {
	mongoc.Model
}

func NewRepoMongo(db string) *repoMongo {
	return &repoMongo{
		Model: mongoc.NewModelAutoTime(storage.GetMongoCli(), db, "test").SetSoftDelete(true),
	}
}

func (repo *repoMongo) Register(ctx context.Context, test *entity.Test) error {
	if test.Id.IsZero() {
		test.Id = primitive.NewObjectID()
	}
	_, err := repo.InsertOne(ctx, test)
	return err
}
