package repository

func Init() {
	// 二级缓存实现
	SetExampleRepository(NewExampleMemoryRepository(NewExampleRedisRepository(NewExampleMongoRepository())))
	// es 单独实例
	SetExampleEsRepository(NewExampleEsRepository())
	//SetExampleRepository(NewExampleMysqlRepository())
	//SetExampleRepository(NewExampleClickhouseRepository())
}
