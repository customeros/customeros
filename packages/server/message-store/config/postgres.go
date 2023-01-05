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

func NewDBConn(host, port, name, user, pass string, maxConn, maxIdleConn, connMaxLifetime int) (*StorageDB, error) {
	connectString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s ", host, port, name, user, pass)
	gormDb, err := gorm.Open(postgres.Open(connectString), initConfig())

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

	sqlDb.SetMaxIdleConns(maxIdleConn)
	sqlDb.SetMaxOpenConns(maxConn)
	sqlDb.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	return &StorageDB{
		SqlDB:  sqlDb,
		GormDB: gormDb,
	}, nil
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
