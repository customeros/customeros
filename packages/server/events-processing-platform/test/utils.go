package test

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	common_logger "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	logger "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	postgrest "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/postgres"
	"github.com/testcontainers/testcontainers-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TestDatabase struct {
	Neo4jContainer testcontainers.Container
	Driver         *neo4j.DriverWithContext
	Repositories   *repository.Repositories
	GormDB         *gorm.DB
}

func SetupTestDatabase() (TestDatabase, func()) {
	testDBs := TestDatabase{}

	appLogger := logger.NewExtendedAppLogger(&common_logger.Config{
		DevMode: true,
	})
	appLogger.InitLogger()

	testDBs.Neo4jContainer, testDBs.Driver = neo4jt.InitTestNeo4jDB()

	postgresContainer, postgresGormDB, _ := postgrest.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		err := postgresContainer.Terminate(ctx)
		if err != nil {
			appLogger.Fatal("Error during container termination")
		}
	}(postgresContainer, context.Background())

	testDBs.GormDB = postgresGormDB
	testDBs.Repositories = repository.InitRepos(testDBs.Driver, postgresGormDB)

	shutdown := func() {
		neo4jt.CloseDriver(*testDBs.Driver)
		neo4jt.Terminate(testDBs.Neo4jContainer, context.Background())
	}
	return testDBs, shutdown
}
