package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
	"log"
)

func InitDB(cfg *config.Config, appLogger logger.Logger) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(cfg); err != nil {
		appLogger.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func main() {
	cfg := loadConfiguration()

	// Initialize logger
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName(config.AppName)

	db, _ := InitDB(cfg, appLogger)
	defer db.SqlDB.Close()

	neo4jDriver, err := config.NewDriver(cfg)
	if err != nil {
		appLogger.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	ctx := context.Background()
	defer neo4jDriver.Close(ctx)

	services := service.InitServices(db.GormDB, &neo4jDriver, appLogger)

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	routes.InitIntegrationRoutes(r, ctx, services)
	routes.InitUserSettingsRoutes(r, ctx, services)
	routes.InitPersonalIntegrationRoutes(r, ctx, services)

	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)

	r.Run(":" + cfg.ApiPort)

	// Flush logs and exit
	appLogger.Sync()
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
