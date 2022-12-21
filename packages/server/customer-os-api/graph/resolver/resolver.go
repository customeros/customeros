package resolver

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service/container"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
)

type Resolver struct {
	Services                    *container.Services
	PostgresRepositoryContainer *commonRepository.PostgresCommonRepositoryContainer
}

func NewResolver(serviceContainer *container.Services, postgresRepositoryContainer *commonRepository.PostgresCommonRepositoryContainer) *Resolver {
	return &Resolver{Services: serviceContainer, PostgresRepositoryContainer: postgresRepositoryContainer}
}
