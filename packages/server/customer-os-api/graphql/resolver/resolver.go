package resolver

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/config"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
)

//go:generate go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	log      logger.Logger
	cfg      *config.Config
	Services *cosapi_services.Services
}

func NewResolver(log logger.Logger, serviceContainer *cosapi_services.Services, cfg *config.Config) *Resolver {
	return &Resolver{
		log:      log,
		cfg:      cfg,
		Services: serviceContainer,
	}
}
