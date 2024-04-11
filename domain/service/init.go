package service

func Init() {
	SetExampleDomainService(NewExampleServiceImpl())
}
