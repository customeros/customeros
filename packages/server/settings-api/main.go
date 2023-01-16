package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/config/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
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

	commonRepositoryContainer := commonRepository.InitCommonRepositories(db.GormDB)
	repositories := repository.InitRepositories(db.GormDB)

	// Setting up Gin
	r := gin.Default()
	r.MaxMultipartMemory = cfg.MaxFileSizeMB << 20

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.GET("/settings",
		commonService.UserToTenantEnhancer(commonRepositoryContainer.UserToTenantRepo),
		commonService.ApiKeyChecker(commonRepositoryContainer.AppKeyRepo, commonService.SETTINGS_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)
			qr := repositories.TenantSettingsRepository.FindForTenantName(tenantName)
			if qr.Error != nil {
				c.AbortWithStatus(500) //todo
				return
			}

			settings := qr.Result.(entity.TenantSettings)
			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(&settings))
		})

	r.POST("/settings/hubspot",
		commonService.UserToTenantEnhancer(commonRepositoryContainer.UserToTenantRepo),
		commonService.ApiKeyChecker(commonRepositoryContainer.AppKeyRepo, commonService.SETTINGS_API),
		func(c *gin.Context) {
			var request dto.TenantSettingsDTO

			if err := c.BindJSON(&request); err != nil {
				println(err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			id := c.Keys["TenantName"].(string)

			settingsForTenant := repositories.TenantSettingsRepository.FindForTenantName(id)
			if settingsForTenant.Error != nil {
				println(settingsForTenant.Error.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			var savedSettings entity.TenantSettings
			if settingsForTenant.Result == nil {
				e := new(entity.TenantSettings)
				e.TenantId = id
				e.HubspotPrivateAppKey = request.HubspotPrivateAppKey

				qr := repositories.TenantSettingsRepository.Save(*e)

				if qr.Error != nil {
					println(qr.Error.Error())
					c.AbortWithStatus(500) //todo
					return
				}
				savedSettings = qr.Result.(entity.TenantSettings)
			} else {
				existingSettings := settingsForTenant.Result.(entity.TenantSettings)
				existingSettings.HubspotPrivateAppKey = request.HubspotPrivateAppKey
				qr := repositories.TenantSettingsRepository.Save(existingSettings)
				if qr.Error != nil {
					println(qr.Error.Error())
					c.AbortWithStatus(500) //todo
					return
				}
				savedSettings = qr.Result.(entity.TenantSettings)
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(&savedSettings))
		})

	r.POST("/settings/zendesk",
		commonService.UserToTenantEnhancer(commonRepositoryContainer.UserToTenantRepo),
		commonService.ApiKeyChecker(commonRepositoryContainer.AppKeyRepo, commonService.SETTINGS_API),
		func(c *gin.Context) {
			var request dto.TenantSettingsDTO

			if err := c.BindJSON(&request); err != nil {
				println(err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			id := c.Keys["TenantName"].(string)

			//find
			settingsForTenant := repositories.TenantSettingsRepository.FindForTenantName(id)
			if settingsForTenant.Error != nil {
				println(settingsForTenant.Error.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			//if new, create
			var savedSettings entity.TenantSettings
			if settingsForTenant.Result == nil {
				e := new(entity.TenantSettings)
				e.TenantId = id
				e.ZendeskAPIKey = request.ZendeskAPIKey
				e.ZendeskAdminEmail = request.ZendeskAdminEmail
				e.ZendeskSubdomain = request.ZendeskSubdomain

				qr := repositories.TenantSettingsRepository.Save(*e)
				if qr.Error != nil {
					println(qr.Error.Error())
					c.AbortWithStatus(500) //todo
					return
				}
				savedSettings = qr.Result.(entity.TenantSettings)
			} else {
				existingSettings := settingsForTenant.Result.(entity.TenantSettings)

				existingSettings.ZendeskAPIKey = request.ZendeskAPIKey
				existingSettings.ZendeskAdminEmail = request.ZendeskAdminEmail
				existingSettings.ZendeskSubdomain = request.ZendeskSubdomain
				qr := repositories.TenantSettingsRepository.Save(existingSettings)
				if qr.Error != nil {
					println(qr.Error.Error())
					c.AbortWithStatus(500) //todo
					return
				}
				savedSettings = qr.Result.(entity.TenantSettings)
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(&savedSettings))
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
