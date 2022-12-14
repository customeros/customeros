package config

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gocache "zgo.at/zcache"
)

var cacheHandler *gocache.Cache

const CommonSchemaPrefix = "tenant_"

type AirbyteStoreDB struct {
	CreationMutex sync.Mutex
	cache         *gocache.Cache
	cfg           *Config
}

type Context struct {
	Schema string
}

func (s *AirbyteStoreDB) CreateDBHandler(ctx *Context) *gorm.DB {
	// Maybe we could use a better mechanism to do this.
	s.CreationMutex.Lock()
	defer s.CreationMutex.Unlock()

	// Double check before moving to Creation precedures
	if gormDb, found := s.cache.Touch(ctx.Schema, gocache.DefaultExpiration); found {
		return gormDb.(*gorm.DB)
	}
	connectionString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s search_path=%s",
		s.cfg.AirbytePostgresDb.Host, s.cfg.AirbytePostgresDb.Port, s.cfg.AirbytePostgresDb.Name, s.cfg.AirbytePostgresDb.User, s.cfg.AirbytePostgresDb.Pwd, ctx.Schema)
	gormDb, err := gorm.Open(postgres.Open(connectionString), initGormConfig())
	if err != nil {
		panic(err)
	}
	sqlDb, err := gormDb.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxIdleConns(s.cfg.AirbytePostgresDb.MaxIdleConn)
	sqlDb.SetMaxOpenConns(s.cfg.AirbytePostgresDb.MaxConn)
	sqlDb.SetConnMaxLifetime(time.Duration(s.cfg.AirbytePostgresDb.ConnMaxLifetime) * time.Second)

	cacheHandler.Set(ctx.Schema, gormDb, gocache.DefaultExpiration)

	return gormDb
}

func (s *AirbyteStoreDB) GetDBHandler(ctx *Context) *gorm.DB {
	db, found := s.cache.Touch(ctx.Schema, gocache.DefaultExpiration)
	if found {
		return db.(*gorm.DB)
	}

	return s.CreateDBHandler(ctx)
}

func InitPoolManager(cfg *Config) *AirbyteStoreDB {
	cacheHandler = gocache.New(10*time.Minute, 10*time.Minute)
	cacheHandler.OnEvicted(func(s string, i interface{}) {
		// https://github.com/go-gorm/gorm/issues/3145
		sql, err := i.(*gorm.DB).DB()
		if err != nil {
			panic(err)
		}
		sql.Close()
	})

	return &AirbyteStoreDB{
		cache: cacheHandler,
		cfg:   cfg,
	}
}
