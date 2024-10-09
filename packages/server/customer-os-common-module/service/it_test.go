package service

import (
	"database/sql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/test"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/testcontainers/testcontainers-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"os"
	"testing"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.DriverWithContext

	postgresContainer testcontainers.Container
	postgresGormDB    *gorm.DB
	postgresSqlDB     *sql.DB

	CommonServices *Services
)

const tenantName = "openline"

func TestMain(m *testing.M) {
	neo4jContainer, driver = neo4jt.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.DriverWithContext, ctx context.Context) {
		neo4jt.CloseDriver(driver)
		neo4jt.TerminateNeo4j(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	postgresContainer, postgresGormDB, postgresSqlDB = neo4jt.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		neo4jt.TerminatePostgres(postgresContainer, ctx)
	}(postgresContainer, context.Background())

	prepareClient()

	os.Exit(m.Run())
}

func prepareClient() {
	appLogger := logger.NewAppLogger(&logger.Config{
		DevMode: true,
	})
	appLogger.InitLogger()

	CommonServices = InitServices(&config.GlobalConfig{}, postgresGormDB, driver, "neo4j", nil, appLogger)
}

func initContext() context.Context {
	ctx := context.Background()

	customCtx := &common.CustomContext{}
	customCtx.Tenant = tenantName

	return common.WithCustomContext(ctx, customCtx)
}

func tearDownTestCase(ctx context.Context) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jtest.CleanupAllData(ctx, driver)
	}
}
