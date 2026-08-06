package repository

import (
	"context"

	"github.com/ihezebin/go-template-ddd/component/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GenericRepository[T any] struct {
	coll *mongo.Collection
}

func NewGenericRepository[T any](name string) *GenericRepository[T] {
	db := storage.MongoDatabase()
	return &GenericRepository[T]{coll: db.Collection(name)}
}

func (r *GenericRepository[T]) FindAll(ctx context.Context) ([]*T, error) {
	cur, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	items := make([]*T, 0)
	if err = cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GenericRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var item T
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GenericRepository[T]) Exists(ctx context.Context, id string) (bool, error) {
	n, err := r.coll.CountDocuments(ctx, bson.M{"_id": id}, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *GenericRepository[T]) InsertIfNotExists(ctx context.Context, id string, doc interface{}) (bool, error) {
	exists, err := r.Exists(ctx, id)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}
	_, err = r.coll.InsertOne(ctx, doc)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *GenericRepository[T]) Upsert(ctx context.Context, id string, doc interface{}) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": id}, doc, options.Replace().SetUpsert(true))
	return err
}

func (r *GenericRepository[T]) Delete(ctx context.Context, id string) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *GenericRepository[T]) UpdateFields(ctx context.Context, id string, fields bson.M) error {
	if len(fields) == 0 {
		return nil
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": fields})
	return err
}

func (r *GenericRepository[T]) FindPage(ctx context.Context, filter bson.M, page, pageSize int64, sort bson.D) ([]*T, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if filter == nil {
		filter = bson.M{}
	}
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().
		SetSkip((page - 1) * pageSize).
		SetLimit(pageSize)
	if len(sort) > 0 {
		opts.SetSort(sort)
	}
	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	items := make([]*T, 0)
	if err = cur.All(ctx, &items); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *GenericRepository[T]) Insert(ctx context.Context, doc any) error {
	_, err := r.coll.InsertOne(ctx, doc)
	return err
}
