package service

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/mailsherpa-api/config"
)

type Services struct {
	Cache                *caches.Cache
	PostgresRepositories *postgres_repository.Repositories
	MailsherpaService    *MailSherpaService
}

func InitServices(
	log logger.Logger,
	postgres *postgres_repository.Repositories,
	config *config.Config,
) *Services {
	services := Services{
		PostgresRepositories: postgres,
		MailsherpaService: NewMailSherpaService(
			log,
			config,
			postgres,
		),
	}

	return &services
}
