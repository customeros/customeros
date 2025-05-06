package database

import (
	"database/sql"
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
	ReadPort        string
	WritePort       string
	User            string
	DBName          string
	Password        string
	MaxConn         int
	MaxIdleConn     int
	ConnMaxLifetime int
	LogLevel        string
	SSLMode         string
}

type DatabaseConnection struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

func NewConnection(dbConfig *DatabaseConfig) (*DatabaseConnection, error) {
	validateConfig(dbConfig)

	// Create write connection
	writeConnectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.WritePort, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	writeDB, err := gorm.Open(postgres.Open(writeConnectString), &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(dbConfig.LogLevel),
	})
	if err != nil {
		log.Printf("Error opening write DB: %v", err)
		return nil, err
	}

	// Configure write connection pool
	writeSqlDB, err := writeDB.DB()
	if err != nil {
		log.Printf("Error getting write DB: %v", err)
		return nil, err
	}

	// Test the write connection
	if err = writeSqlDB.Ping(); err != nil {
		log.Printf("Error pinging write DB: %v", err)
		return nil, err
	}

	// Create read connection
	readConnectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.ReadPort, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	readDB, err := gorm.Open(postgres.Open(readConnectString), &gorm.Config{
		Logger: initLog(dbConfig.LogLevel),
	})
	if err != nil {
		log.Printf("Error opening read DB: %v", err)
		return nil, err
	}

	// Configure read connection pool
	readSqlDB, err := readDB.DB()
	if err != nil {
		log.Printf("Error getting read DB: %v", err)
		return nil, err
	}

	// Test the read connection
	if err = readSqlDB.Ping(); err != nil {
		log.Printf("Error pinging read DB: %v", err)
		return nil, err
	}

	// Configure both connection pools
	configureConnectionPool(writeSqlDB, dbConfig)
	configureConnectionPool(readSqlDB, dbConfig)

	return &DatabaseConnection{
		ReadDB:  readDB,
		WriteDB: writeDB,
	}, nil
}

func configureConnectionPool(db *sql.DB, config *DatabaseConfig) {
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	db.SetMaxIdleConns(config.MaxIdleConn)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	db.SetMaxOpenConns(config.MaxConn)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Hour)
}

func validateConfig(config *DatabaseConfig) {
	switch {
	case config == nil:
		log.Fatalf("Database config is nil")
	case config.Host == "":
		log.Fatalf("Database host config is empty")
	case config.ReadPort == "":
		log.Fatalf("Database read port config is empty")
	case config.WritePort == "":
		log.Fatalf("Database write port config is empty")
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
