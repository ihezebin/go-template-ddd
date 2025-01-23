package storage

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Suggest interface{}

func TestInsertElasticsearchData(t *testing.T) {
	err := InitElasticsearchClient(context.Background(), "http://localhost:9200")
	if err != nil {
		t.Fatal(err)
	}

	client := ElasticsearchClient()
	ctx := context.Background()

	// 创建索引
	indexSettings := strings.NewReader(`{
		"settings": {
			"analysis": {
				"analyzer": {
					"text_analyzer": {
						"tokenizer": "ik_max_word",
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
				"user_id": {"type": "keyword"},
				"name": {
					"type": "text",
					"analyzer": "text_analyzer",
					"search_analyzer": "ik_smart",
					"copy_to": "search_text"
				},
				"age": {"type": "integer"},
				"interests": {
					"type": "keyword",
					"copy_to": "search_text"
				},
				"city": {
					"type": "keyword",
					"copy_to": "search_text"
				},
				"avatar": {
					"type": "keyword",
					"index": false
				},
				"search_text": {
					"type": "text",
					"analyzer": "text_analyzer",
					"search_analyzer": "ik_smart"
				},
				"suggestion": {
					"type": "completion",
					"analyzer": "keyword_analyzer"
				}
			}
		}
	}`)

	resp, err := client.Indices.Create("user").Raw(indexSettings).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Acknowledged {
		t.Fatal("index creation not acknowledged")
	}

	// 插入文档
	doc1 := map[string]interface{}{
		"user_id":    "U010",
		"name":       "王十二",
		"age":        31,
		"interests":  []string{"电影", "阅读"},
		"city":       "长沙",
		"avatar":     "https://example.com/avatar10.jpg",
		"suggestion": []string{"电影", "阅读", "长沙"},
	}

	_, err = client.Index("user").Id("1").Request(doc1).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	doc2 := map[string]interface{}{
		"user_id":    "U009",
		"name":       "郑十一",
		"age":        24,
		"interests":  []string{"摄影", "游泳"},
		"city":       "武汉",
		"avatar":     "https://example.com/avatar9.jpg",
		"suggestion": []string{"摄影", "游泳", "武汉"},
	}

	_, err = client.Index("user").Id("2").Request(doc2).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchElasticsearchClient(t *testing.T) {
	err := InitElasticsearchClient(context.Background(), "http://localhost:9200")
	if err != nil {
		t.Fatal(err)
	}

	client := ElasticsearchClient()
	ctx := context.Background()

	// 搜索
	resp, err := client.Search().Index("user").Request(&search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"search_text": {Query: "biancheng"},
			},
		},
	}).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range resp.Hits.Hits {
		var source map[string]interface{}
		err = json.Unmarshal(v.Source_, &source)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Hit: %+v", source)
	}

	// 自动补全
	resp, err = client.Search().Index("user").Raw(strings.NewReader(`{
		"suggest": {
			"text": "b",
			"user_suggest": {
				"completion": {
					"field": "suggestion",
					"skip_duplicates": true,
					"size": 5
				}
			}
		}
	}`)).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Suggest != nil {
		for _, suggestions := range resp.Suggest {
			for _, suggestion := range suggestions {
				switch s := suggestion.(type) {
				case *types.CompletionSuggest:
					for _, option := range s.Options {
						t.Logf("Completion Suggestion: %+v", option.Text)
					}
				case *types.PhraseSuggest:
					for _, option := range s.Options {
						t.Logf("Phrase Suggestion: %+v", option.Text)
					}
				case *types.TermSuggest:
					for _, option := range s.Options {
						t.Logf("Term Suggestion: %+v", option.Text)
					}
				default:
					// 获取实际类型
					t.Logf("Unknown suggestion type: %T", s)
				}
			}
		}
	}
}

func TestDeleteIndex(t *testing.T) {
	err := InitElasticsearchClient(context.Background(), "http://localhost:9200")
	if err != nil {
		t.Fatal(err)
	}

	client := ElasticsearchClient()
	ctx := context.Background()

	resp, err := client.Indices.Delete("user").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Acknowledged {
		t.Fatal("index deletion not acknowledged")
	}
}
