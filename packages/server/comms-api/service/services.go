package service

import (
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/repository"
)

type Services struct {
	MailService       MailService
	CustomerOsService CustomerOSService
}

func InitServices(graphqlClient *graphql.Client, config *c.Config, db *c.StorageDB) *Services {
	cosService := NewCustomerOSService(graphqlClient, config)
	apiKeyRepository := repository.NewApiKeyRepository(db)
	services := Services{
		CustomerOsService: cosService,
		MailService:       NewMailService(config, cosService, apiKeyRepository),
	}

	return &services
}
