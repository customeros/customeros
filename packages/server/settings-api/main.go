package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/mapper"
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

	commonRepositoryContainer := commonRepository.InitRepositories(db.GormDB, &neo4jDriver)
	services := service.InitServices(db.GormDB, appLogger)

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.GET("/integrations",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)

			tenantIntegrationSettings, activeServices, err := services.TenantSettingsService.GetForTenant(tenantName)

			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(tenantIntegrationSettings, activeServices))
		})

	r.POST("/integration",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			var request map[string]interface{}

			if err := c.BindJSON(&request); err != nil {
				println(err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			tenantName := c.Keys["TenantName"].(string)

			tenantIntegrationSettings, activeServices, err := services.TenantSettingsService.SaveIntegrationData(tenantName, request)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(tenantIntegrationSettings, activeServices))
		})

	r.DELETE("/integration/:identifier",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			identifier := c.Param("identifier")
			if identifier == "" {
				c.JSON(500, gin.H{"error": "integration identifier is empty"})
				return
			}
			tenantName := c.Keys["TenantName"].(string)

			data, activeServices, err := services.TenantSettingsService.ClearIntegrationData(tenantName, identifier)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(data, activeServices))
		})

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
