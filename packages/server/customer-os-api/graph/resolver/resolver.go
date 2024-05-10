package resolver

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

//go:generate go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	log      logger.Logger
	cfg      *config.Config
	Services *service.Services
	Clients  *grpc_client.Clients
}

func NewResolver(log logger.Logger, serviceContainer *service.Services, grpcContainer *grpc_client.Clients, cfg *config.Config) *Resolver {
	return &Resolver{
		log:      log,
		cfg:      cfg,
		Services: serviceContainer,
		Clients:  grpcContainer,
	}
}
