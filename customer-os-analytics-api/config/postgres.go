package config

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type StorageDB struct {
	SqlDB  *sql.DB
	GormDB *gorm.DB
}

func NewDBConn(host, port, name, user, pass string, maxConn, maxIdleConn, connMaxLifetime int) (*StorageDB, error) {
	connectString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s ", host, port, name, user, pass)
	storageDB := new(StorageDB)
	gormdb, err := gorm.Open(postgres.Open(connectString), &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	if storageDB.SqlDB, err = gormdb.DB(); err != nil {
		return nil, err
	}
	if err = storageDB.SqlDB.Ping(); err != nil {
		return nil, err
	}
	storageDB.SqlDB.SetMaxIdleConns(maxIdleConn)
	storageDB.SqlDB.SetMaxOpenConns(maxConn)
	storageDB.SqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	storageDB.GormDB = gormdb
	return &StorageDB{
		SqlDB:  storageDB.SqlDB,
		GormDB: storageDB.GormDB,
	}, nil
}
