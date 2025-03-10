package database

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	DBName          string
	Password        string
	MaxConn         int
	MaxIdleConn     int
	ConnMaxLifetime int
	LogLevel        string
	SSLMode         string
}

func NewConnection(dbConfig *DatabaseConfig) (*gorm.DB, error) {
	validateConfig(dbConfig)

	portInt, err := strconv.Atoi(dbConfig.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %w", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, portInt, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func validateConfig(config *DatabaseConfig) {
	switch {
	case config == nil:
		log.Fatalf("Database config is nil")
	case config.Host == "":
		log.Fatalf("Database host config is empty")
	case config.Port == "":
		log.Fatalf("Database port config is empty")
	case config.User == "":
		log.Fatalf("Database user config is empty")
	case config.Password == "":
		log.Fatalf("Database password config is empty")
	case config.DBName == "":
		log.Fatalf("Database name config is empty")
	case config.SSLMode == "":
		log.Fatalf("Database SSLMode config is empty")
	}
	return
}
