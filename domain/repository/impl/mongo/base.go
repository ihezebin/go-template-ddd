package mongo

import (
	"context"
	"github.com/whereabouts/sdk/db/mongoc"
	"github.com/whereabouts/web-template-ddd/component/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Base Rewrite the methods provided in mongoc.Base to support modification, new automatic update time fields and soft delete.
// Soft delete is off by default; And the time is stored in the form of a timestamp.
// 重写 mongoc.Base 中提供的方法, 使其支持修改、新增自动更新时间字段和软删除, 默认关闭软删除, 且时间以时间戳形式存储
type Base struct {
	kernel             *mongoc.Base
	softDelete         bool
	createTimeFieldKey string
	updateTimeFieldKey string
	deleteTimeFieldKey string
}

func NewMongoBase(db string, collection string) *Base {
	return &Base{
		kernel:             mongoc.NewBaseModel(storage.GetMongoCli(), db, collection),
		softDelete:         false,
		createTimeFieldKey: "create_time",
		updateTimeFieldKey: "update_time",
		deleteTimeFieldKey: "delete_time",
	}
}

// Kernel If the rewriting part cannot meet the needs, use the original core Base
// 重写部分不能满足需求的则使用原核心 mongoc.Base
func (base *Base) Kernel() *mongoc.Base {
	return base.kernel
}

func (base *Base) SetSoftDelete(softDelete bool) *Base {
	base.softDelete = softDelete
	return base
}

func (base *Base) FindOne(ctx context.Context, filter bson.M, result interface{}, opts ...*options.FindOneOptions) error {
	if base.softDelete {
		filter[base.deleteTimeFieldKey] = bson.M{
			"$eq": 0,
		}
	}
	return base.Kernel().FindOne(ctx, filter, result, opts...)
}

func (base *Base) Find(ctx context.Context, filter bson.M, results interface{}, opts ...*options.FindOptions) error {
	if base.softDelete {
		filter[base.deleteTimeFieldKey] = bson.M{
			"$eq": 0,
		}
	}
	return base.Kernel().Find(ctx, filter, results, opts...)
}

func (base *Base) Count(ctx context.Context, filter bson.M, opts ...*options.CountOptions) (int64, error) {
	if base.softDelete {
		filter[base.deleteTimeFieldKey] = bson.M{
			"$eq": 0,
		}
	}
	return base.Kernel().Count(ctx, filter, opts...)
}

func (base *Base) DeleteOne(ctx context.Context, filter bson.M) error {
	if base.softDelete {
		filter[base.deleteTimeFieldKey] = bson.M{
			"$eq": 0,
		}

		now := time.Now().Unix()
		_, err := base.Kernel().UpdateOne(ctx, filter, bson.M{
			"$set": bson.M{
				base.deleteTimeFieldKey: now,
				base.updateTimeFieldKey: now,
			},
		})
		return err
	}

	_, err := base.Kernel().DeleteOne(ctx, filter)
	return err
}

func (base *Base) DeleteMany(ctx context.Context, filter bson.M) error {
	if base.softDelete {
		filter[base.deleteTimeFieldKey] = bson.M{
			"$eq": 0,
		}

		now := time.Now().Unix()
		_, err := base.Kernel().UpdateMany(ctx, filter, bson.M{
			"$set": bson.M{
				base.deleteTimeFieldKey: now,
				base.updateTimeFieldKey: now,
			},
		})
		return err
	}

	_, err := base.Kernel().DeleteMany(ctx, filter)
	return err
}

func (base *Base) UpdateOne(ctx context.Context, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if base.softDelete {
		filter[base.deleteTimeFieldKey] = bson.M{
			"$eq": 0,
		}
	}

	if _, ok := update["$set"]; ok {
		setter, err := base.convert2Doc(ctx, update["$set"])
		if err != nil {
			return nil, err
		}
		setter[base.updateTimeFieldKey] = time.Now().Unix()
		update["$set"] = setter
	}

	return base.Kernel().UpdateOne(ctx, filter, update, opts...)
}

func (base *Base) UpdateMany(ctx context.Context, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if base.softDelete {
		filter[base.deleteTimeFieldKey] = bson.M{
			"$eq": 0,
		}
	}

	if _, ok := update["$set"]; ok {
		setter, err := base.convert2Doc(ctx, update["$set"])
		if err != nil {
			return nil, err
		}
		setter[base.updateTimeFieldKey] = time.Now().Unix()
		update["$set"] = setter
	}

	return base.Kernel().UpdateMany(ctx, filter, update, opts...)
}

func (base *Base) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	doc, err := base.convert2Doc(ctx, document)
	if err != nil {
		return nil, err
	}

	if base.softDelete {
		doc[base.deleteTimeFieldKey] = int64(0)
	}

	now := time.Now().Unix()
	doc[base.createTimeFieldKey] = now
	doc[base.updateTimeFieldKey] = now
	return base.Kernel().InsertOne(ctx, doc, opts...)
}

func (base *Base) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	if base.softDelete {
		for i := range documents {
			doc, err := base.convert2Doc(ctx, documents[i])
			if err != nil {
				return nil, err
			}
			doc[base.deleteTimeFieldKey] = int64(0)
			documents[i] = doc
		}
	}

	for i := range documents {
		doc, err := base.convert2Doc(ctx, documents[i])
		if err != nil {
			return nil, err
		}
		now := time.Now().Unix()
		doc[base.createTimeFieldKey] = now
		doc[base.updateTimeFieldKey] = now
		documents[i] = doc
	}

	return base.Kernel().InsertMany(ctx, documents, opts...)
}

func (base *Base) convert2Doc(ctx context.Context, document interface{}) (bson.M, error) {
	var doc bson.M
	data, err := bson.Marshal(document)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(data, &doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
