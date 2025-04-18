package container

import (
	"github.com/customeros/customeros/packages/runner/integrity-checker/caches"
	"github.com/customeros/customeros/packages/runner/integrity-checker/config"
	"github.com/customeros/customeros/packages/runner/integrity-checker/logger"
	neo4jRepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type Container struct {
	Cfg      *config.Config
	Log      logger.Logger
	Neo4j    *neo4jRepository.Repositories
	Postgres *postgresRepository.Repositories
	Cache    *caches.Cache
}
