package resolver

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service/container"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
)

//go:generate go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Services                    *container.Services
	PostgresRepositoryContainer *commonRepository.PostgresCommonRepositoryContainer
}

func NewResolver(serviceContainer *container.Services, postgresRepositoryContainer *commonRepository.PostgresCommonRepositoryContainer) *Resolver {
	return &Resolver{Services: serviceContainer, PostgresRepositoryContainer: postgresRepositoryContainer}
}
