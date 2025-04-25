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

type DbConnections struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

func NewConnection(dbConfig *DatabaseConfig) (*DbConnections, error) {
	validateConfig(dbConfig)

	// Create write connection
	writeConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.WritePort, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	writeDB, err := gorm.Open(postgres.Open(writeConnString), &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(dbConfig.LogLevel),
	})
	if err != nil {
		log.Printf("Error opening write DB: %v", err)
		return nil, err
	}

	// Configure connection pool for write DB
	writeSQL, err := writeDB.DB()
	if err != nil {
		log.Printf("Error getting write DB: %v", err)
		return nil, err
	}

	// Test the write connection
	if err = writeSQL.Ping(); err != nil {
		log.Printf("Error pinging write DB: %v", err)
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	writeSQL.SetMaxIdleConns(dbConfig.MaxIdleConn)
	// SetMaxOpenConns sets the maximum number of open connections to the database
	writeSQL.SetMaxOpenConns(dbConfig.MaxConn)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	writeSQL.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Hour)

	// Create read connection
	readConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.ReadPort, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	readDB, err := gorm.Open(postgres.Open(readConnString), &gorm.Config{
		Logger: initLog(dbConfig.LogLevel),
	})
	if err != nil {
		log.Printf("Error opening read DB: %v", err)
		return nil, err
	}

	// Configure connection pool for read DB
	readSQL, err := readDB.DB()
	if err != nil {
		log.Printf("Error getting read DB: %v", err)
		return nil, err
	}

	// Test the read connection
	if err = readSQL.Ping(); err != nil {
		log.Printf("Error pinging read DB: %v", err)
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	readSQL.SetMaxIdleConns(dbConfig.MaxIdleConn)
	// SetMaxOpenConns sets the maximum number of open connections to the database
	readSQL.SetMaxOpenConns(dbConfig.MaxConn)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	readSQL.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Hour)

	return &DbConnections{
		ReadDB:  readDB,
		WriteDB: writeDB,
	}, nil
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
