package service

import (
	"context"
	"database/sql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/config"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/test/neo4j"
	postgres "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/test/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.DriverWithContext

	postgresGormDB    *gorm.DB
	postgresSqlDB     *sql.DB
	serviceContainer  *Services
	postgresContainer testcontainers.Container
)

func TestMain(m *testing.M) {
	neo4jContainer, driver = neo4jt.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.DriverWithContext, ctx context.Context) {
		neo4jt.CloseDriver(driver)
		neo4jt.Terminate(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	postgresContainer, postgresGormDB, postgresSqlDB = postgres.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		err := postgresContainer.Terminate(ctx)
		if err != nil {
			log.Fatal("Error during container termination")
		}
	}(postgresContainer, context.Background())

	prepareClient()

	os.Exit(m.Run())
}

func prepareClient() {
	log.Printf("******************PREPARING CLIENT******************")
	appLogger := logger.NewAppLogger(&logger.Config{
		DevMode: true,
	})
	appLogger.InitLogger()

	commonServices := commonService.InitServices(postgresGormDB, driver)
	serviceContainer = InitServices(&config.Config{}, commonServices, postgresGormDB)
	log.Printf("%v", serviceContainer)
}
