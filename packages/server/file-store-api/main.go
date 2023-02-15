package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/service"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
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
	services := service.InitServices(cfg, db.GormDB)

	// Setting up Gin
	r := gin.Default()
	r.MaxMultipartMemory = cfg.MaxFileSizeMB << 20

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.POST("/file",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.FILE_STORE_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)

			multipartFileHeader, err := c.FormFile("file")
			if err != nil {
				c.AbortWithStatus(500) //todo
				return
			}

			fileEntity, err := services.FileService.UploadSingleFile(tenantName, multipartFileHeader)
			if err != nil {
				c.AbortWithStatus(500) //todo
				return
			}

			c.JSON(200, MapFileEntityToDTO(cfg, fileEntity))
		})
	r.GET("/file/:id",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.FILE_STORE_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)

			byId, err := services.FileService.GetById(tenantName, c.Param("id"))
			if err != nil && err.Error() != "record not found" {
				c.AbortWithStatus(500) //todo
				return
			}
			if err != nil && err.Error() == "record not found" {
				c.AbortWithStatus(404)
				return
			}

			c.JSON(200, MapFileEntityToDTO(cfg, byId))
		})
	r.GET("/file/:id/download",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.FILE_STORE_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)

			byId, bytes, err := services.FileService.DownloadSingleFile(tenantName, c.Param("id"))
			if err != nil && err.Error() != "record not found" {
				c.AbortWithStatus(500) //todo
				return
			}
			if err != nil && err.Error() == "record not found" {
				c.AbortWithStatus(404)
				return
			}

			c.Header("Content-Disposition", "attachment; filename="+byId.Name)
			c.Header("Content-Type", fmt.Sprintf("%s", byId.MIME))
			c.Header("Accept-Length", fmt.Sprintf("%d", len(bytes)))
			c.Writer.Write(bytes)
		})
	r.GET("/file/:id/base64",
		commonService.UserToTenantEnhancer(ctx, commonRepositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepo, commonService.FILE_STORE_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)

			base64Encoded, err := services.FileService.Base64Image(tenantName, c.Param("id"))
			if err != nil && err.Error() != "record not found" {
				c.AbortWithStatus(500) //todo
				return
			}
			if err != nil && err.Error() == "record not found" {
				c.AbortWithStatus(404)
				return
			}

			bytes := []byte(*base64Encoded)
			c.Writer.Write(bytes)
		})

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

func MapFileEntityToDTO(cfg *config.Config, fileEntity *entity.File) *dto.File {
	return mapper.MapFileEntityToDTO(fileEntity, cfg.ApiServiceUrl)
}
