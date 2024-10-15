package config

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

type StorageDB struct {
	SqlDB  *sql.DB
	GormDB *gorm.DB
}

func NewDBConn(cfg *Config) (*StorageDB, error) {
	connectString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s ", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Db, cfg.Postgres.User, cfg.Postgres.Password)
	gormDb, err := gorm.Open(postgres.Open(connectString), initConfig(cfg))

	var sqlDb *sql.DB
	if err != nil {
		return nil, err
	}
	if sqlDb, err = gormDb.DB(); err != nil {
		return nil, err
	}
	if err = sqlDb.Ping(); err != nil {
		return nil, err
	}

	sqlDb.SetMaxIdleConns(cfg.Postgres.MaxIdleConn)
	sqlDb.SetMaxOpenConns(cfg.Postgres.MaxConn)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.Postgres.ConnMaxLifetime) * time.Second)

	return &StorageDB{
		SqlDB:  sqlDb,
		GormDB: gormDb,
	}, nil
}

// initConfig Initialize Config
func initConfig(cfg *Config) *gorm.Config {
	return &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(cfg),
	}
}

// initLog Connection Log Configuration
func initLog(cfg *Config) logger.Interface {
	var logLevel = logger.Silent
	switch cfg.Postgres.LogLevel {
	case "ERROR":
		logLevel = logger.Error
	case "WARN":
		logLevel = logger.Warn
	case "INFO":
		logLevel = logger.Info
	}
	newLogger := logger.New(log.New(io.MultiWriter(os.Stdout), "\r\n", log.LstdFlags), logger.Config{
		Colorful:      true,
		LogLevel:      logLevel,
		SlowThreshold: time.Second,
	})
	return newLogger
}
