package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/service"
	"github.com/sirupsen/logrus"
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
	graphqlClient := graphql.NewClient(cfg.Service.CustomerOsAPI)
	services := service.InitServices(cfg, graphqlClient)

	jwtTennantUserService := service.NewJWTTenantUserService(commonRepositoryContainer, cfg)

	// Setting up Gin
	r := gin.Default()
	r.MaxMultipartMemory = cfg.MaxFileSizeMB << 20

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS", "POST", "GET")

	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	corsConfig.AddAllowHeaders("X-Openline-JWT")

	r.Use(cors.New(corsConfig))

	r.POST("/file",
		jwtTennantUserService.GetJWTTenantUserEnhancer(),
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME_OR_TENANT, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.FILE_STORE_API),
		func(ctx *gin.Context) {
			tenantName, _ := ctx.Keys["TenantName"].(string)
			userEmail, _ := ctx.Keys["UserEmail"].(string)

			multipartFileHeader, err := ctx.FormFile("file")
			if err != nil {
				ctx.AbortWithStatusJSON(500, map[string]string{"error": "missing field file"}) //todo
				return
			}

			fileEntity, err := services.FileService.UploadSingleFile(userEmail, tenantName, multipartFileHeader)
			if err != nil {
				ctx.AbortWithStatusJSON(500, map[string]string{"error": fmt.Sprintf("Error Uploading File %v", err)}) //todo
				return
			}

			ctx.JSON(200, MapFileEntityToDTO(cfg, fileEntity))
		})
	r.GET("/file/:id",
		jwtTennantUserService.GetJWTTenantUserEnhancer(),
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.FILE_STORE_API),
		func(ctx *gin.Context) {
			tenantName, _ := ctx.Keys["TenantName"].(string)
			userEmail, _ := ctx.Keys["UserEmail"].(string)

			byId, err := services.FileService.GetById(userEmail, tenantName, ctx.Param("id"))
			if err != nil && err.Error() != "record not found" {
				ctx.AbortWithStatus(500) //todo
				return
			}
			if err != nil && err.Error() == "record not found" {
				ctx.AbortWithStatus(404)
				return
			}

			ctx.JSON(200, MapFileEntityToDTO(cfg, byId))
		})
	r.GET("/file/:id/download",
		jwtTennantUserService.GetJWTTenantUserEnhancer(),
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.FILE_STORE_API),
		func(ctx *gin.Context) {
			tenantName, _ := ctx.Keys["TenantName"].(string)
			userEmail, _ := ctx.Keys["UserEmail"].(string)

			_, err := services.FileService.DownloadSingleFile(userEmail, tenantName, ctx.Param("id"), ctx, ctx.Query("inline") == "true")
			if err != nil && err.Error() != "record not found" {
				ctx.AbortWithStatus(500) //todo
				return
			}
			if err != nil && err.Error() == "record not found" {
				ctx.AbortWithStatus(404)
				return
			}

			//ctx.Header("Accept-Length", fmt.Sprintf("%d", len(bytes)))
			//ctx.Writer.Write(bytes)
		})
	r.GET("/file/:id/base64",
		jwtTennantUserService.GetJWTTenantUserEnhancer(),
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.FILE_STORE_API),
		func(ctx *gin.Context) {
			tenantName, _ := ctx.Keys["TenantName"].(string)
			userEmail, _ := ctx.Keys["UserEmail"].(string)

			base64Encoded, err := services.FileService.Base64Image(userEmail, tenantName, ctx.Param("id"))
			if err != nil && err.Error() != "record not found" {
				ctx.AbortWithStatus(500) //todo
				return
			}
			if err != nil && err.Error() == "record not found" {
				ctx.AbortWithStatus(404)
				return
			}

			bytes := []byte(*base64Encoded)
			ctx.Writer.Write(bytes)
		})

	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)

	r.GET("/jwt",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.FILE_STORE_API),
		func(ctx *gin.Context) {
			jwtTennantUserService.MakeJWT(ctx)
		})

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

func MapFileEntityToDTO(cfg *config.Config, fileEntity *model.File) *dto.File {
	return mapper.MapFileEntityToDTO(fileEntity, cfg.ApiServiceUrl)
}
