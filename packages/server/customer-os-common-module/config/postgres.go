package config

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type PostgresDB struct {
	SqlDB  *sql.DB
	GormDB *gorm.DB

	AsyncSqlDB  *sql.DB
	AsyncGormDB *gorm.DB
}

func InitPostgres(cfg *CommonConfig) (*PostgresDB, error) {
	var err error
	db := &PostgresDB{}

	db.SqlDB, db.GormDB, err = NewPostgresDBConn(cfg.PostgresConfig.Host, cfg.PostgresConfig.Port, cfg.PostgresConfig.User, cfg.PostgresConfig.Password, cfg.PostgresConfig.Db, cfg.PostgresConfig.LogLevel, cfg.PostgresConfig.MaxConn, cfg.PostgresConfig.MaxIdleConn, cfg.PostgresConfig.ConnMaxLifetime)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
		return nil, err
	}
	db.AsyncSqlDB, db.AsyncGormDB, err = NewPostgresDBConn(cfg.PostgresAsyncConfig.Host, cfg.PostgresAsyncConfig.Port, cfg.PostgresAsyncConfig.User, cfg.PostgresAsyncConfig.Password, cfg.PostgresAsyncConfig.Db, cfg.PostgresAsyncConfig.LogLevel, cfg.PostgresAsyncConfig.MaxConn, cfg.PostgresAsyncConfig.MaxIdleConn, cfg.PostgresAsyncConfig.ConnMaxLifetime)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
		return nil, err
	}

	return db, nil
}

func (db *PostgresDB) Close() {
	db.SqlDB.Close()
	db.AsyncSqlDB.Close()
}

func NewPostgresDBConn(host, port, user, password, db, logLevel string, maxConn, maxIdleConn, connMaxLifetime int) (*sql.DB, *gorm.DB, error) {
	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, db)
	gormDb, err := gorm.Open(postgres.Open(connectString), initConfig(logLevel))

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

	sqlDb.SetMaxIdleConns(maxConn)
	sqlDb.SetMaxOpenConns(maxIdleConn)
	sqlDb.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	return sqlDb, gormDb, nil
}

// initConfig Initialize Config
func initConfig(logLevel string) *gorm.Config {
	return &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(logLevel),
	}
}

// initLog Connection Log Configuration
func initLog(logLevel string) gormLogger.Interface {
	postgresLogLevel := gormLogger.Silent
	switch logLevel {
	case "ERROR":
		postgresLogLevel = gormLogger.Error
	case "WARN":
		postgresLogLevel = gormLogger.Warn
	case "INFO":
		postgresLogLevel = gormLogger.Info
	}
	newLogger := gormLogger.New(log.New(io.MultiWriter(os.Stdout), "\r\n", log.LstdFlags), gormLogger.Config{
		Colorful:      true,
		LogLevel:      postgresLogLevel,
		SlowThreshold: time.Second,
	})
	return newLogger
}
