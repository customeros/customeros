package container

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/repository"
)

type Container struct {
	Cfg          *config.Config
	Log          logger.Logger
	Repositories *repository.Repositories
}
