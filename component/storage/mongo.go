package storage

import (
	"context"
	"github.com/whereabouts/sdk/db/mongoc"
)

var gMongoCli *mongoc.Client

func InitMongo(ctx context.Context, config mongoc.Config) (err error) {
	gMongoCli, err = mongoc.NewClient(ctx, config)
	return
}

func GetMongoCli() *mongoc.Client {
	return gMongoCli
}
