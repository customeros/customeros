package rest

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/test"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.DriverWithContext

	postgresContainer testcontainers.Container
	postgresGormDB    *gorm.DB
	serviceContainer  *service.Services
)

const tenantName = "openline"
const testUserId = "test-user-id"
const testContactId = "test-contact-id"
const testPlayerId = "test-player-id"

func TestMain(m *testing.M) {
	neo4jContainer, driver = neo4jt.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.DriverWithContext, ctx context.Context) {
		neo4jt.CloseDriver(driver)
		neo4jt.TerminateNeo4j(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	postgresContainer, postgresGormDB, _ = neo4jt.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		neo4jt.TerminatePostgres(postgresContainer, ctx)
	}(postgresContainer, context.Background())

	prepareClient()

	os.Exit(m.Run())
}

func tearDownTestCase(ctx context.Context) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jtest.CleanupAllData(ctx, driver)
	}
}

func prepareClient() {
	log.Printf("******************PREPARING CLIENT******************")
	appLogger := logger.NewAppLogger(&logger.Config{
		DevMode: true,
	})
	appLogger.InitLogger()
	commonServices := commonService.InitServices(&commonConfig.GlobalConfig{}, postgresGormDB, driver, "neo4j", nil, appLogger)
	serviceContainer = service.InitServices(appLogger, driver, &config.Config{}, commonServices, nil, postgresGormDB)
	dataloader.NewDataLoader(serviceContainer)
	log.Printf("%v", serviceContainer)
}
