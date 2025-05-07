package database

import (
	"context"
	"database/sql"
	"fmt"
	commonconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func DatabaseConfigFromWarehouseDbConfig(cfg commonconf.DataWarehouseConfig) DatabaseConfig {
	return DatabaseConfig{
		Host:            cfg.Host,
		ReadPort:        cfg.ReadPort,
		WritePort:       cfg.WritePort,
		User:            cfg.User,
		DBName:          cfg.DBName,
		Password:        cfg.Password,
		MaxConn:         cfg.MaxConn,
		MaxIdleConn:     cfg.MaxIdleConn,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		LogLevel:        cfg.LogLevel,
	}
}

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
	ReadSqlDB  *sql.DB
	ReadDB     *gorm.DB
	WriteSqlDB *sql.DB
	WriteDB    *gorm.DB
}

func NewConnection(dbConfig *DatabaseConfig) (*DbConnections, error) {
	validateConfig(dbConfig)

	writeConfig := postgres.Config{
		DSN: fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbConfig.Host, dbConfig.WritePort, dbConfig.User, dbConfig.Password,
			dbConfig.DBName, dbConfig.SSLMode,
		),
		PreferSimpleProtocol: true,
	}

	// Open write connection with explicit error handling
	writeDB, err := gorm.Open(postgres.New(writeConfig), &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(dbConfig.LogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open write database connection: %w", err)
	}

	// Configure write connection pool
	writeSQLDB, err := writeDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get write database connection: %w", err)
	}

	// Configure connection pooling
	writeSQLDB.SetMaxIdleConns(dbConfig.MaxIdleConn)
	writeSQLDB.SetMaxOpenConns(dbConfig.MaxConn)
	writeSQLDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Hour)

	// Test connection with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = writeSQLDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("write database connection test failed: %w", err)
	}

	// Use the same approach for read connection
	readConfig := postgres.Config{
		DSN: fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbConfig.Host, dbConfig.ReadPort, dbConfig.User, dbConfig.Password,
			dbConfig.DBName, dbConfig.SSLMode,
		),
		PreferSimpleProtocol: true,
	}

	readDB, err := gorm.Open(postgres.New(readConfig), &gorm.Config{
		Logger: initLog(dbConfig.LogLevel),
	})
	if err != nil {
		// Clean up the write connection before returning
		if sqlDB, _ := writeDB.DB(); sqlDB != nil {
			_ = sqlDB.Close()
		}
		return nil, fmt.Errorf("failed to open read database connection: %w", err)
	}

	// Configure read connection pool
	readSQLDB, err := readDB.DB()
	if err != nil {
		// Clean up connections before returning
		if sqlDB, _ := writeDB.DB(); sqlDB != nil {
			_ = sqlDB.Close()
		}
		return nil, fmt.Errorf("failed to get read database connection: %w", err)
	}

	// Configure connection pooling for read DB
	readSQLDB.SetMaxIdleConns(dbConfig.MaxIdleConn)
	readSQLDB.SetMaxOpenConns(dbConfig.MaxConn)
	readSQLDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Hour)

	// Test read connection
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = readSQLDB.PingContext(ctx); err != nil {
		// Clean up connections before returning
		if sqlDB, _ := writeDB.DB(); sqlDB != nil {
			_ = sqlDB.Close()
		}
		return nil, fmt.Errorf("read database connection test failed: %w", err)
	}

	return &DbConnections{
		ReadDB:     readDB,
		ReadSqlDB:  readSQLDB,
		WriteDB:    writeDB,
		WriteSqlDB: writeSQLDB,
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
