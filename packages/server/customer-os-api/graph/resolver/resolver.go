package resolver

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service/container"
)

//go:generate go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Services                    *container.Services
	PostgresRepositoryContainer *repository.PostgresRepositoryContainer
}

func NewResolver(serviceContainer *container.Services, postgresRepositoryContainer *repository.PostgresRepositoryContainer) *Resolver {
	return &Resolver{Services: serviceContainer, PostgresRepositoryContainer: postgresRepositoryContainer}
}
