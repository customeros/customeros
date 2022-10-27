package service

type ServiceContainer struct {
	ContactService ContactService
}

func InitServices() *ServiceContainer {
	return &ServiceContainer{
		ContactService: NewContactService(),
	}
}
