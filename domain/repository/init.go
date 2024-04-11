package repository

func Init() {
	// 二级缓存实现
	SetExampleRepository(NewExampleMemoryRepository(NewExampleRedisRepository(NewExampleMongoRepository())))
	//SetExampleRepository(NewExampleMysqlRepository())
}
