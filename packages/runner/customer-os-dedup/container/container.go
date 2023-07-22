package container

import (
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/repository"
)

type Container struct {
	Cfg                     *config.Config
	Log                     logger.Logger
	Repositories            *repository.Repositories
	CustomerOsGraphQLClient *graphql.Client
}
