package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	comlog "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/test/neo4j"
	postgrest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/test/postgres"
	"github.com/testcontainers/testcontainers-go"
	"os"
	"testing"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.DriverWithContext

	repositories *Repositories
)

const tenantName = "openline"

func SetupTestLogger() logger.Logger {
	testLogger := logger.NewExtendedAppLogger(&comlog.Config{
		DevMode: true,
	})
	testLogger.InitLogger()
	return testLogger
}

func TestMain(m *testing.M) {
	log := SetupTestLogger()

	neo4jContainer, driver = neo4jt.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.DriverWithContext, ctx context.Context) {
		neo4jt.CloseDriver(driver)
		neo4jt.Terminate(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	postgresContainer, postgresGormDB, _ := postgrest.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		err := postgresContainer.Terminate(ctx)
		if err != nil {
			log.Fatal("Error during container termination")
		}
	}(postgresContainer, context.Background())

	repositories = InitRepos(driver, postgresGormDB, "neo4j")

	os.Exit(m.Run())
}
