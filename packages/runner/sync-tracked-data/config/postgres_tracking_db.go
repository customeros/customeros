package config

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func NewPostgresTrackingClient(cfg *Config) (*sql.DB, *gorm.DB, error) {
	connectString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s ",
		cfg.TrackingPostgresDb.Host, cfg.TrackingPostgresDb.Port, cfg.TrackingPostgresDb.Name, cfg.TrackingPostgresDb.User, cfg.TrackingPostgresDb.Pwd)
	gormDb, err := gorm.Open(postgres.Open(connectString), initConfig(cfg))

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

	sqlDb.SetMaxIdleConns(cfg.TrackingPostgresDb.MaxIdleConn)
	sqlDb.SetMaxOpenConns(cfg.TrackingPostgresDb.MaxConn)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.TrackingPostgresDb.ConnMaxLifetime) * time.Second)

	return sqlDb, gormDb, nil
}
