package repository

import testRepo "github.com/ihezebin/go-template-ddd/domain/repository/test"

func Init(db string) {
	testRepo.SetRepository(testRepo.NewRepoRedis(db))
	testRepo.SetRepository(testRepo.NewRepoMemory(db))
	testRepo.SetRepository(testRepo.NewRepoMongo(db))
}
