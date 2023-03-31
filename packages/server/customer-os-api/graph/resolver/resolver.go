package resolver

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
)

//go:generate go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Services                    *service.Services
	PostgresRepositoryContainer *commonRepository.Repositories
	Clients                     *grpc_client.Clients
}

func NewResolver(
	serviceContainer *service.Services,
	postgresRepositoryContainer *commonRepository.Repositories,
	grpcContainer *grpc_client.Clients) *Resolver {
	return &Resolver{
		Services:                    serviceContainer,
		PostgresRepositoryContainer: postgresRepositoryContainer,
		Clients:                     grpcContainer,
	}
}
