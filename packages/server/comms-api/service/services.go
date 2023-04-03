package service

import (
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
)

type Services struct {
	MailService       MailService
	CustomerOsService CustomerOSService
}

func InitServices(graphqlClient *graphql.Client, config *c.Config) *Services {
	cosService := NewCustomerOSService(graphqlClient, config)
	services := Services{
		CustomerOsService: cosService,
		MailService:       NewMailService(config, cosService),
	}

	return &services
}
