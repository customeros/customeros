package integration_tests

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/repository/entity"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB = "testdb"

type TestContainerDbConfig struct {
	Host string
	Port string
	Name string
	User string
	Pwd  string
}

func createPostgresTestDBContainer() (testcontainers.Container, *pgxpool.Pool, *TestContainerDbConfig, error) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, nil, nil, err
	}
	port, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		return nil, nil, nil, err
	}
	host, err := dbContainer.Host(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	dbURI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/%v", host, port.Port(), testDB)
	connPool, err := pgxpool.Connect(context.Background(), dbURI)
	if err != nil {
		return nil, nil, nil, err
	}

	testContainerDbConfig := TestContainerDbConfig{
		Host: host,
		Port: port.Port(),
		User: "postgres",
		Pwd:  "postgres",
		Name: testDB,
	}

	return dbContainer, connPool, &testContainerDbConfig, err
}

func createAllTables(db *gorm.DB) {
	db.Exec("create schema if not exists derived")

	var err error
	err = db.AutoMigrate(&entity.ApplicationEntity{})
	if err != nil {
		log.Panicf("Error creating %v table", entity.ApplicationEntity{}.TableName())
	}

	err = db.AutoMigrate(&entity.PageViewEntity{})
	if err != nil {
		log.Panicf("Error creating %v table", entity.PageViewEntity{}.TableName())
	}

	err = db.AutoMigrate(&entity.SessionEntity{})
	if err != nil {
		log.Panicf("Error creating %v table", entity.SessionEntity{}.TableName())
	}
}

func InitTestDB() (testcontainers.Container, *config.StorageDB) {
	var err error
	dbContainer, _, dbConfig, err := createPostgresTestDBContainer()
	if err != nil {
		log.Fatal("Error setup PostgresDB container")
	}

	var db *config.StorageDB
	if db, err = config.NewDBConn(
		dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.User, dbConfig.Pwd, 10, 10, 0); err != nil {
		log.Panicf("Coud not open db connection: %s", err.Error())
	}

	createAllTables(db.GormDB)

	return dbContainer, db
}
