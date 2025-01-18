package container

import (
	"github.com/machinebox/graphql"
	"github.com/customeros/customeros/packages/runner/customer-os-dedup/config"
	"github.com/customeros/customeros/packages/runner/customer-os-dedup/logger"
	"github.com/customeros/customeros/packages/runner/customer-os-dedup/repository"
)

type Container struct {
	Cfg                     *config.Config
	Log                     logger.Logger
	Repositories            *repository.Repositories
	CustomerOsGraphQLClient *graphql.Client
}
