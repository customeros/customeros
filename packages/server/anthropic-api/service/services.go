package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/anthorpic-api/config"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
)

type Services struct {
	PostgresRepositories *postgresRepository.Repositories

	AnthropicService AnthropicService
}

func InitServices(cfg *config.Config, db *config.StorageDB) *Services {
	services := &Services{
		PostgresRepositories: postgresRepository.InitRepositories(db.GormDB),
	}

	services.AnthropicService = NewAnthropicService(cfg)

	return services
}
