package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/config/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/repository"
	"log"
)

const apiPort = "10000"

func InitDB(cfg *config.Config) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Db,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.MaxConn,
		cfg.Postgres.MaxIdleConn,
		cfg.Postgres.ConnMaxLifetime); err != nil {
		log.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func init() {
	logger.Logger = logger.New(log.New(log.Default().Writer(), "", log.Ldate|log.Ltime|log.Lmicroseconds), logger.Config{
		Colorful: true,
		LogLevel: logger.Info,
	})
}

//// Declare a simple handler for pingpong as a request accepting behavior
//func ApiKeyChecker(appKeyRepo repository.AppKeyRepository) func(c *gin.Context) {
//	return func(c *gin.Context) {
//		kh := c.GetHeader("X-Openline-API-KEY")
//		if kh != "" {
//
//			keyResult := appKeyRepo.FindByKey(c, kh)
//
//			if keyResult.Error != nil {
//				c.AbortWithStatus(401)
//				return
//			}
//
//			appKey := keyResult.Result.(*entity.AppKeyEntity)
//
//			if appKey == nil {
//				c.AbortWithStatus(401)
//				return
//			} else {
//				// todo set tenant in context
//			}
//
//			c.Next()
//			// illegal request, terminate the current process
//		} else {
//			c.AbortWithStatus(401)
//			return
//		}
//
//	}
//}

func main() {
	cfg := loadConfiguration()

	db, _ := InitDB(cfg)
	defer db.SqlDB.Close()

	repositoryContainer := repository.InitRepositories(db.GormDB)

	if repositoryContainer == nil {
		panic("a")
	}

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.GET("/", (func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	}))

	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)

	port := cfg.ApiPort
	if port == "" {
		port = apiPort
	}

	r.Run(":" + port)
}

func loadConfiguration() *config.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] Error loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}
