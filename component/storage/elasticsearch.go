package storage

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/ihezebin/soup/logger"
	"github.com/pkg/errors"
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

	pingResp, err := client.Info().Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "ping es error")
	}

	logger.Infof(ctx, "ping es info: %+v", pingResp)

	if err := initEsMapping(ctx, client); err != nil {
		return errors.Wrapf(err, "init es mapping error")
	}

	elasticsearchClient = client
	return nil
}

// Deprecated: pingEsWithApi 使用 esapi 的方式 ping es
func pingEsWithApi(ctx context.Context, client *elasticsearch.TypedClient) error {
	pingRes, err := esapi.PingRequest{
		Pretty: true,
		Human:  true,
	}.Do(ctx, client)
	if err != nil {
		return errors.Wrapf(err, "ping es error")
	}
	defer pingRes.Body.Close()

	if pingRes.IsError() {
		return errors.New("ping es error")
	}

	return nil
}
