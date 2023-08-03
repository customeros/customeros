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

func NewPostgresDBConn(cfg PostgresConfig) (*StorageDB, error) {
	connectString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s ",
		cfg.Host, cfg.Port, cfg.Db, cfg.User, cfg.Password)
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

	sqlDb.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDb.SetMaxOpenConns(cfg.MaxConn)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return &StorageDB{
		SqlDB:  sqlDb,
		GormDB: gormDb,
	}, nil
}

// initConfig Initialize Config
func initConfig(cfg PostgresConfig) *gorm.Config {
	return &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(cfg),
	}
}

// initLog Connection Log Configuration
func initLog(cfg PostgresConfig) logger.Interface {
	var logLevel = logger.Silent
	switch cfg.LogLevel {
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
