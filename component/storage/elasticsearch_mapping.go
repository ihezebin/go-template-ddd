package storage

import (
	"bytes"
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/pkg/errors"
)

func initEsMapping(ctx context.Context, client *elasticsearch.TypedClient) error {
	// example
	if err := initEsExampleMapping(ctx, client); err != nil {
		return errors.Wrapf(err, "init es example mapping error")
	}

	// others...

	return nil
}

/*
PUT /example

	{
	  "settings": {
	    "analysis": {
	      "analyzer": {
	        "text_analyzer": {
	          "tokenizer": "ik_smart",
	          "filter": "pinyin_filter"
	        },
	        "keyword_analyzer": {
	          "tokenizer": "keyword",
	          "filter": "pinyin_filter"
	        }
	      },
	      "filter": {
	        "pinyin_filter": {
	          "type": "pinyin",
	          "keep_first_letter": false,
	          "keep_full_pinyin": false,
	          "keep_joined_full_pinyin": true,
	          "keep_original": true,
	          "limit_first_letter_length": 16,
	          "remove_duplicated_term": true,
	          "none_chinese_pinyin_tokenize": false
	        }
	      }
	    }
	  },
	  "mappings": {
	    "properties": {
	      "id": {
	        "type": "keyword"
	      },
	      "username": {
	        "type": "text",
	        "analyzer": "text_analyzer"
	      },
	      "email": {
	        "type": "keyword"
	      },
	      "password": {
	        "type": "keyword"
	      },
	      "salt": {
	        "type": "keyword"
	      }
	    }
	  }
	}
*/
func initEsExampleMapping(ctx context.Context, client *elasticsearch.TypedClient) error {
	// 判断是否存在
	exists, err := client.Indices.Exists("example").Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "check es example mapping exists error")
	}

	if exists {
		return nil
	}

	// 创建索引
	settings := types.NewIndexSettings()
	err = settings.UnmarshalJSON(bytes.NewBufferString(`{
	    "analysis": {
	      "analyzer": {
	        "text_analyzer": {
	          "tokenizer": "ik_smart",
	          "filter": "pinyin_filter"
	        },
	        "keyword_analyzer": {
	          "tokenizer": "keyword",
	          "filter": "pinyin_filter"
	        }
	      },
	      "filter": {
	        "pinyin_filter": {
	          "type": "pinyin",
	          "keep_first_letter": false,
	          "keep_full_pinyin": false,
	          "keep_joined_full_pinyin": true,
	          "keep_original": true,
	          "limit_first_letter_length": 16,
	          "remove_duplicated_term": true,
	          "none_chinese_pinyin_tokenize": false
	        }
	      }
	    }
	  }`).Bytes())
	if err != nil {
		return errors.Wrapf(err, "unmarshal es example mapping settings error")
	}

	mapping := types.NewTypeMapping()
	err = mapping.UnmarshalJSON(bytes.NewBufferString(`{
	    "properties": {
	      "id": {
	        "type": "keyword"
	      },
	      "username": {
	        "type": "text",
	        "analyzer": "text_analyzer"
	      },
	      "email": {
	        "type": "keyword"
	      },
	      "password": {
	        "type": "keyword"
	      },
	      "salt": {
	        "type": "keyword"
	      }
	    }
	  }`).Bytes())
	if err != nil {
		return errors.Wrapf(err, "unmarshal es example mapping error")
	}

	res, err := client.Indices.Create("example").Settings(settings).Mappings(mapping).Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "create es example mapping error")
	}

	if !res.Acknowledged || !res.ShardsAcknowledged {
		return errors.New("create es example mapping not acknowledged or shards not acknowledged")
	}

	return nil
}
