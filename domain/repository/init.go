package repository

import (
	"github.com/ihezebin/go-template-ddd/domain/entity"
)

var (
	ExampleGenericRepository *GenericRepository[entity.Example]
)

/*
通用 CRUD 实现

若要实现单独方法，通用 CRUD 继续用 GenericRepository[T]，专用逻辑用组合嵌套即可

	type ExampleSpecialRepository struct {
		*GenericRepository[entity.Example]
	}

	func NewExampleSpecialRepository() *ExampleSpecialRepository {
		return &ExampleSpecialRepository{
			GenericRepository: NewGenericRepository[entity.Example](),
		}
	}

专用方法

	func (r *ExampleSpecialRepository) FindRecentFailed(ctx context.Context, limit int64) ([]*entity.Example, error) {
		opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit)
		cur, err := r.coll.Find(ctx, bson.M{"status": "failed"}, opts)
		...
	}
*/
func InitGenericRepository() {
	ExampleGenericRepository = NewGenericRepository[entity.Example](exampleCollectionName)
}

// 完全自定义接口 + Mongo 实现，适合和通用 CRUD 差很多、或要 mock 的实体
func InitInterfaceRepository() {
	// 二级缓存实现
	SetExampleRepository(NewExampleMemoryRepository(NewExampleRedisRepository(NewExampleMongoRepository())))
	// es 单独实例
	SetExampleEsRepository(NewExampleEsRepository())
	//SetExampleRepository(NewExampleMysqlRepository())
	//SetExampleRepository(NewExampleClickhouseRepository())
}
