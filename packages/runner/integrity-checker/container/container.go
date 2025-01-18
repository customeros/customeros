package container

import (
	"github.com/customeros/customeros/packages/runner/integrity-checker/caches"
	"github.com/customeros/customeros/packages/runner/integrity-checker/config"
	"github.com/customeros/customeros/packages/runner/integrity-checker/logger"
	"github.com/customeros/customeros/packages/runner/integrity-checker/repository"
)

type Container struct {
	Cfg          *config.Config
	Log          logger.Logger
	Repositories *repository.Repositories
	Cache        *caches.Cache
}
