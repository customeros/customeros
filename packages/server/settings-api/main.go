package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/config/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/mapper"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	//"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
	"log"
)

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

func main() {
	cfg := loadConfiguration()

	db, _ := InitDB(cfg)
	defer db.SqlDB.Close()

	neo4jDriver, err := config.NewDriver(cfg)
	if err != nil {
		logrus.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	ctx := context.Background()
	defer neo4jDriver.Close(ctx)

	commonRepositoryContainer := commonRepository.InitRepositories(db.GormDB, &neo4jDriver)
	services := service.InitServices(db.GormDB)

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.GET("/settings",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.SETTINGS_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)

			tenant, err := services.TenantSettingsService.GetForTenant(tenantName)

			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(tenant))
		})

	r.POST("/settings/hubspot",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.SETTINGS_API),
		func(c *gin.Context) {
			var request dto.TenantSettingsHubspotDTO

			if err := c.BindJSON(&request); err != nil {
				println(err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			tenantName := c.Keys["TenantName"].(string)

			data, err := services.TenantSettingsService.SaveHubspotData(tenantName, request)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(data))
		})

	r.DELETE("/settings/hubspot",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.SETTINGS_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)

			data, err := services.TenantSettingsService.ClearHubspotData(tenantName)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(data))
		})

	r.POST("/settings/zendesk",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.SETTINGS_API),
		func(c *gin.Context) {
			var request dto.TenantSettingsZendeskDTO

			if err := c.BindJSON(&request); err != nil {
				println(err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			tenantName := c.Keys["TenantName"].(string)

			data, err := services.TenantSettingsService.SaveZendeskData(tenantName, request)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(data))
		})

	r.DELETE("/settings/zendesk",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.SETTINGS_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)

			data, err := services.TenantSettingsService.ClearZendeskData(tenantName)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(data))
		})

	r.POST("/settings/smartSheet",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.SETTINGS_API),
		func(c *gin.Context) {
			var request dto.TenantSettingsSmartSheetDTO

			if err := c.BindJSON(&request); err != nil {
				println(err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			tenantName := c.Keys["TenantName"].(string)

			data, err := services.TenantSettingsService.SaveSmartSheetData(tenantName, request)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(data))
		})

	r.DELETE("/settings/smartSheet",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.SETTINGS_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)

			data, err := services.TenantSettingsService.ClearSmartSheetData(tenantName)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(data))
		})

	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)

	r.Run(":" + cfg.ApiPort)
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
