package container

import (
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/repository"
)

type Container struct {
	Cfg                     *config.Config
	Log                     logger.Logger
	Repositories            *repository.Repositories
	CustomerOsGraphQLClient *graphql.Client
	RawDataStoreDB          *config.RawDataStoreDB
}
