package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository/postgres/entity"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

func startPostgresContainer(ctx context.Context) (testcontainers.Container, error) {
	request := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
}

func InitTestDB() (testcontainers.Container, *gorm.DB, *sql.DB) {
	var ctx = context.Background()
	var err error
	postgresContainer, err := startPostgresContainer(ctx)
	if err != nil {
		log.Panic("Container should start")
	}

	port, err := postgresContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		panic(err)
	}
	host, err := postgresContainer.Host(context.Background())
	if err != nil {
		panic(err)
	}

	connectString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s ", host, port.Port(), "testdb", "postgres", "postgres")
	gormDb, err := gorm.Open(postgres.Open(connectString), &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(),
	})
	if err != nil {
		panic(err)
	}

	sqlDb, err := gormDb.DB()
	if err != nil {
		panic(err)
	}
	if err = sqlDb.Ping(); err != nil {
		panic(err)
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(2)
	sqlDb.SetConnMaxLifetime(time.Duration(1) * time.Second)

	createAllTables(gormDb)

	return postgresContainer, gormDb, sqlDb
}

func createAllTables(db *gorm.DB) {
	db.Exec("create schema if not exists derived")

	var err error
	err = db.AutoMigrate(&entity.AppKey{})
	if err != nil {
		log.Panicf("Error creating %v table", entity.AppKeyEntity{}.TableName())
	}
}

// initLog Connection Log Configuration
func initLog() logger.Interface {
	//f, _ := os.Create("gorm.log")
	newLogger := logger.New(log.New(io.MultiWriter(os.Stdout), "\r\n", log.LstdFlags), logger.Config{
		Colorful:      true,
		LogLevel:      logger.Info,
		SlowThreshold: time.Second,
	})
	return newLogger
}

func Close(closer io.Closer, resourceName string) {
	err := closer.Close()
	if err != nil {
		log.Panicf("%s should close", resourceName)
	}
}

func Terminate(container testcontainers.Container, ctx context.Context) {
	err := container.Terminate(ctx)
	if err != nil {
		log.Fatal("Container should stop")
	}
}
