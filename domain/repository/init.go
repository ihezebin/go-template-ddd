package repository

import "github.com/ihezebin/go-template-ddd/domain/repository/test"

func Init(db string) {
	test.SetRepository(test.NewRepoRedis(db))
	test.SetRepository(test.NewRepoMemory(db))
	test.SetRepository(test.NewRepoMongo(db))
}
