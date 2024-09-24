package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/test"
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

func TestMain(m *testing.M) {
	neo4jContainer, driver = neo4jt.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.DriverWithContext, ctx context.Context) {
		neo4jt.CloseDriver(driver)
		neo4jt.TerminateNeo4j(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	postgresContainer, postgresGormDB, _ := neo4jt.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		neo4jt.TerminatePostgres(postgresContainer, ctx)
	}(postgresContainer, context.Background())

	repositories = InitRepos(driver, "neo4j", postgresGormDB)

	os.Exit(m.Run())
}
