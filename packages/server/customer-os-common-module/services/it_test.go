package service

import (
	"database/sql"
	"os"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	test "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/test"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.DriverWithContext

	postgresContainer testcontainers.Container
	postgresGormDB    *gorm.DB
	postgresSqlDB     *sql.DB

	rabbitMqContainer testcontainers.Container
	rabbitMqConn      *amqp091.Connection
)

const tenantName = "openline"

func TestMain(m *testing.M) {
	neo4jContainer, driver = test.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.DriverWithContext, ctx context.Context) {
		test.CloseDriver(driver)
		test.TerminateNeo4j(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	postgresContainer, postgresGormDB, postgresSqlDB = test.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		test.TerminatePostgres(postgresContainer, ctx)
	}(postgresContainer, context.Background())

	//rabbitMqContainer, rabbitMqConn = test.InitTestRabbitMQ()
	//defer func(rabbitMqContainer testcontainers.Container, ctx context.Context) {
	//	test.TerminateRabbitMq(rabbitMqContainer, ctx)
	//}(rabbitMqContainer, context.Background())

	prepareClient()

	os.Exit(m.Run())
}

func prepareClient() {
	appLogger := logger.NewAppLogger(&logger.Config{
		DevMode: true,
	})
	appLogger.InitLogger()
}

func InitContext() context.Context {
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
