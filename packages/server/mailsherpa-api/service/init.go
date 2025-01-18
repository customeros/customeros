package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/mailsherpa-api/config"
)

type Services struct {
	Cache                *caches.Cache
	PostgresRepositories *repository.Repositories
	MailsherpaService    *MailSherpaService
}

func InitServices(
	log logger.Logger,
	postgres *repository.Repositories,
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
