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
	connectString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s ",
		cfg.PostgresDb.Host, cfg.PostgresDb.Port, cfg.PostgresDb.Name, cfg.PostgresDb.User, cfg.PostgresDb.Pwd)
	gormDb, err := gorm.Open(postgres.Open(connectString), initConfig())

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
func initConfig() *gorm.Config {
	return &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(),
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
