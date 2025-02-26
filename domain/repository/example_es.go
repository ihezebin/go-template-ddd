package repository

import (
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ihezebin/go-template-ddd/component/storage"
	"github.com/ihezebin/go-template-ddd/domain/entity"
)

const exampleEsIndexName = "example"

type exampleEsRepository struct {
	client *elasticsearch.TypedClient
}

func NewExampleEsRepository() ExampleRepository {
	return &exampleEsRepository{
		client: storage.ElasticsearchClient(),
	}
}

var _ ExampleRepository = (*exampleEsRepository)(nil)

func (e *exampleEsRepository) InsertOne(ctx context.Context, example *entity.Example) error {
	if example.Id == "" {
		example.Id = primitive.NewObjectID().Hex()
	}
	_, err := e.client.Index(exampleEsIndexName).Id(example.Id).Document(example).Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "insert es example error")
	}

	return nil
}

func (e *exampleEsRepository) FindByUsername(ctx context.Context, username string) (*entity.Example, error) {

	query := types.NewQuery()
	query.Match = map[string]types.MatchQuery{
		"username": {
			Query: username,
		},
	}

	res, err := e.client.Search().Index(exampleEsIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "find es example by username error")
	}

	if res.Hits.Total.Value == 0 {
		return nil, nil
	}

	example := entity.Example{}
	err = json.Unmarshal(res.Hits.Hits[0].Source_, &example)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal es example error")
	}

	return &example, nil
}

func (e *exampleEsRepository) FindByEmail(ctx context.Context, email string) (*entity.Example, error) {
	query := types.NewQuery()
	query.Match = map[string]types.MatchQuery{
		"email": {
			Query: email,
		},
	}

	res, err := e.client.Search().Index(exampleEsIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "find es example by email error")
	}

	if res.Hits.Total.Value == 0 {
		return nil, nil
	}

	example := entity.Example{}
	err = json.Unmarshal(res.Hits.Hits[0].Source_, &example)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal es example error")
	}

	return &example, nil
}
