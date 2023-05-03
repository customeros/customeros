package config

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type StorageDB struct {
	SqlDB  *sql.DB
	GormDB *gorm.DB
}

func NewDBConn(host, port, name, user, pass string, maxConn, maxIdleConn, connMaxLifetime int) (*StorageDB, error) {
	connectString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s ", host, port, name, user, pass)
	gormDb, err := gorm.Open(postgres.Open(connectString), &gorm.Config{})

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
