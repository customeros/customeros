package config

import (
	"fmt"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gocache "zgo.at/zcache"
)

var cacheHandler *gocache.Cache

type RawDataStoreDB struct {
	CreationMutex sync.Mutex
	cache         *gocache.Cache
	cfg           *Config
}

type Context struct {
	Schema string
}

func (s *RawDataStoreDB) CreateDBHandler(ctx *Context) *gorm.DB {
	// Maybe we could use a better mechanism to do this.
	s.CreationMutex.Lock()
	defer s.CreationMutex.Unlock()

	// Double check before moving to Creation precedures
	if gormDb, found := s.cache.Touch(ctx.Schema, gocache.DefaultExpiration); found {
		return gormDb.(*gorm.DB)
	}
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s search_path=%s",
		s.cfg.Postgres.Host, s.cfg.Postgres.Port, s.cfg.RawDataStoreDBName, s.cfg.Postgres.User, s.cfg.Postgres.Password, ctx.Schema)
	gormDb, err := gorm.Open(postgres.Open(connectionString), initGormConfig(s.cfg))
	if err != nil {
		panic(err)
	}
	sqlDb, err := gormDb.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxIdleConns(s.cfg.Postgres.MaxIdleConn)
	sqlDb.SetMaxOpenConns(s.cfg.Postgres.MaxConn)
	sqlDb.SetConnMaxLifetime(time.Duration(s.cfg.Postgres.ConnMaxLifetime) * time.Second)

	cacheHandler.Set(ctx.Schema, gormDb, gocache.DefaultExpiration)

	return gormDb
}

func (s *RawDataStoreDB) GetDBHandler(ctx *Context) *gorm.DB {
	db, found := s.cache.Touch(ctx.Schema, gocache.DefaultExpiration)
	if found {
		return db.(*gorm.DB)
	}

	return s.CreateDBHandler(ctx)
}

func InitPoolManager(cfg *Config) *RawDataStoreDB {
	cacheHandler = gocache.New(10*time.Minute, 10*time.Minute)
	cacheHandler.OnEvicted(func(s string, i interface{}) {
		// https://github.com/go-gorm/gorm/issues/3145
		sql, err := i.(*gorm.DB).DB()
		if err != nil {
			panic(err)
		}
		sql.Close()
	})

	return &RawDataStoreDB{
		cache: cacheHandler,
		cfg:   cfg,
	}
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
