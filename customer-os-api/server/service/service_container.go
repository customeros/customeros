package service

import "github.com/openline-ai/openline-customer-os/customer-os-api/config"

type ServiceContainer struct {
	ContactService ContactService
}

func InitServices(cfg *config.Config) *ServiceContainer {
	return &ServiceContainer{
		ContactService: NewContactService(cfg),
	}
}
