package service

import (
	"context"
	"database/sql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/test/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var (
	postgresGormDB    *gorm.DB
	postgresSqlDB     *sql.DB
	serviceContainer  *Services
	postgresContainer testcontainers.Container
)

func TestMain(m *testing.M) {
	postgresContainer, postgresGormDB, postgresSqlDB = postgrest.InitTestDB()
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
	serviceContainer = InitServices(&config.Config{}, postgresGormDB)
	log.Printf("%v", serviceContainer)
}
