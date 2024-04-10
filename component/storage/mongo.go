package storage

import (
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoDatabase *mongo.Database

func MongoStorageClient() *mongo.Client {
	return mongoClient
}

func MongoStorageDatabase() *mongo.Database {
	return mongoDatabase
}

func InitMongoStorageClient(ctx context.Context, dsn string) error {
	u, err := url.Parse(dsn)
	if err != nil {
		return errors.Wrap(err, "mongo dsn parse error")
	}

	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" {
		return errors.New("mongo db name is empty")
	}

	option := options.Client().ApplyURI(dsn)
	if err = option.Validate(); err != nil {
		return errors.Wrap(err, "mongo dsn validate error")
	}
	client, err := mongo.Connect(ctx, option)
	if err != nil {
		return errors.Wrap(err, "mongo connect error")
	}

	mongoClient = client
	mongoDatabase = client.Database(dbName)

	return nil
}
