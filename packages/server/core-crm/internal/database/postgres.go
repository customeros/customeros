package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
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

	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)

	gormDb, err := gorm.Open(postgres.Open(connectString), &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(dbConfig.LogLevel),
	})
	if err != nil {
		log.Printf("Error opening DB: %v", err)
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := gormDb.DB()
	if err != nil {
		log.Printf("Error getting DB: %v", err)
		return nil, err
	}

	// Test the connection
	if err = sqlDB.Ping(); err != nil {
		log.Printf("Error pinging DB: %v", err)
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(dbConfig.MaxConn)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Hour)

	return gormDb, nil
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
	}
}

func initLog(logLevel string) gormlogger.Interface {
	postgresLogLevel := gormlogger.Silent
	switch logLevel {
	case "ERROR":
		postgresLogLevel = gormlogger.Error
	case "WARN":
		postgresLogLevel = gormlogger.Warn
	case "INFO":
		postgresLogLevel = gormlogger.Info
	}
	newLogger := gormlogger.New(log.New(io.MultiWriter(os.Stdout), "\r\n", log.LstdFlags), gormlogger.Config{
		Colorful:      true,
		LogLevel:      postgresLogLevel,
		SlowThreshold: time.Second,
	})
	return newLogger
}
