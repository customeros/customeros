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

func NewPostgresClient(cfg *Config) (*sql.DB, *gorm.DB, error) {
	connectString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s ",
		cfg.PostgresDb.Host, cfg.PostgresDb.Port, cfg.PostgresDb.Db, cfg.PostgresDb.User, cfg.PostgresDb.Password)
	gormDb, err := gorm.Open(postgres.Open(connectString), initGormConfig(cfg))

	var sqlDb *sql.DB
	if err != nil {
		return nil, nil, err
	}
	if sqlDb, err = gormDb.DB(); err != nil {
		return nil, nil, err
	}
	if err = sqlDb.Ping(); err != nil {
		return nil, nil, err
	}

	sqlDb.SetMaxIdleConns(cfg.PostgresDb.MaxIdleConn)
	sqlDb.SetMaxOpenConns(cfg.PostgresDb.MaxConn)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.PostgresDb.ConnMaxLifetime) * time.Second)

	return sqlDb, gormDb, nil
}

// initConfig Initialize Config
func initGormConfig(cfg *Config) *gorm.Config {
	return &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(cfg),
	}
}

// initLog Connection Log Configuration
func initLog(cfg *Config) logger.Interface {
	var logLevel = logger.Silent
	switch cfg.PostgresDb.LogLevel {
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
