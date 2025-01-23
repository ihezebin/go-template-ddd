package storage

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
)

var elasticsearchClient *elasticsearch.TypedClient

func ElasticsearchClient() *elasticsearch.TypedClient {
	return elasticsearchClient
}

func InitElasticsearchClient(ctx context.Context, url string) error {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return err
	}

	elasticsearchClient = client
	return nil
}
