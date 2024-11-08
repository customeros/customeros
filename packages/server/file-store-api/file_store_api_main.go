package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/service"
	"github.com/opentracing/opentracing-go"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

func main() {
	parentCtx := context.Background()
	ctx, cancel := signal.NotifyContext(parentCtx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg := loadConfiguration()

	// Initialize Logging
	appLogger := initLogger(cfg)

	// Initialize Tracing
	tracingCloser := initTracing(cfg, appLogger)
	if tracingCloser != nil {
		defer tracingCloser.Close()
	}

	// Setting up Neo4j
	neo4jDriver, err := commonconf.NewNeo4jDriver(cfg.Neo4j)
	if err != nil {
		appLogger.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	// initialize db
	db, _ := InitDB(cfg)
	defer db.SqlDB.Close()

	commonServices := commonservice.InitServices(&commonconf.GlobalConfig{}, db.GormDB, &neo4jDriver, cfg.Neo4j.Database, nil, appLogger)

	graphqlClient := graphql.NewClient(cfg.Service.CustomerOsAPI)
	services := service.InitServices(cfg, commonServices, graphqlClient, appLogger)

	jwtTennantUserService := service.NewJWTTenantUserService(cfg)

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
		tracing.TracingEnhancer(ctx, "POST /file"),
		jwtTennantUserService.GetJWTTenantUserEnhancer(),
		security.TenantUserContextEnhancer(security.USERNAME_OR_TENANT, commonServices.Neo4jRepositories, security.WithCache(commonServices.Cache)),
		security.ApiKeyCheckerHTTP(commonServices.PostgresRepositories.TenantWebhookApiKeyRepository, commonServices.PostgresRepositories.AppKeyRepository, security.FILE_STORE_API, security.WithCache(commonServices.Cache)),
		func(c *gin.Context) {
			ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreApi)

			cdnUpload := c.Request.FormValue("cdnUpload") == "true"
			basePath := c.Request.FormValue("basePath")
			fileId := c.Request.FormValue("fileId")

			multipartFileHeader, err := c.FormFile("file")
			if err != nil {
				c.AbortWithStatusJSON(500, map[string]string{"error": "missing field file"}) //todo
				return
			}

			fileEntity, err := services.FileService.UploadSingleFile(ctx, basePath, fileId, multipartFileHeader, cdnUpload)
			if err != nil {
				c.AbortWithStatusJSON(500, map[string]string{"error": fmt.Sprintf("Error Uploading File %v", err)}) //todo
				return
			}

			c.JSON(http.StatusOK, MapFileEntityToDTO(cfg, fileEntity))
		})
	r.GET("/file/:id",
		tracing.TracingEnhancer(ctx, "GET /file/:id"),
		jwtTennantUserService.GetJWTTenantUserEnhancer(),
		security.TenantUserContextEnhancer(security.USERNAME_OR_TENANT, commonServices.Neo4jRepositories, security.WithCache(commonServices.Cache)),
		security.ApiKeyCheckerHTTP(commonServices.PostgresRepositories.TenantWebhookApiKeyRepository, commonServices.PostgresRepositories.AppKeyRepository, security.FILE_STORE_API, security.WithCache(commonServices.Cache)),
		func(c *gin.Context) {
			ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreApi)

			byId, err := services.FileService.GetById(ctx, c.Param("id"))
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
		tracing.TracingEnhancer(ctx, "GET /file/:id/download"),
		jwtTennantUserService.GetJWTTenantUserEnhancer(),
		security.ApiKeyCheckerHTTP(commonServices.PostgresRepositories.TenantWebhookApiKeyRepository, commonServices.PostgresRepositories.AppKeyRepository, security.FILE_STORE_API, security.WithCache(commonServices.Cache)),
		security.TenantUserContextEnhancer(security.USERNAME_OR_TENANT, commonServices.Neo4jRepositories, security.WithCache(commonServices.Cache)),
		func(c *gin.Context) {
			ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreApi)

			_, err := services.FileService.DownloadSingleFile(ctx, c.Param("id"), c, c.Query("inline") == "true")
			if err != nil && err.Error() != "record not found" {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if err != nil && err.Error() == "record not found" {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		})
	r.GET("/file/:id/base64",
		tracing.TracingEnhancer(ctx, "GET /file/:id/base64"),
		jwtTennantUserService.GetJWTTenantUserEnhancer(),
		security.TenantUserContextEnhancer(security.USERNAME_OR_TENANT, commonServices.Neo4jRepositories, security.WithCache(commonServices.Cache)),
		security.ApiKeyCheckerHTTP(commonServices.PostgresRepositories.TenantWebhookApiKeyRepository, commonServices.PostgresRepositories.AppKeyRepository, security.FILE_STORE_API, security.WithCache(commonServices.Cache)),
		func(c *gin.Context) {
			ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreApi)

			base64Encoded, err := services.FileService.Base64Image(ctx, c.Param("id"))
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
	r.GET("/file/:id/public-url",
		tracing.TracingEnhancer(ctx, "GET /file/:id/public-url"),
		jwtTennantUserService.GetJWTTenantUserEnhancer(),
		security.TenantUserContextEnhancer(security.USERNAME_OR_TENANT, commonServices.Neo4jRepositories, security.WithCache(commonServices.Cache)),
		security.ApiKeyCheckerHTTP(commonServices.PostgresRepositories.TenantWebhookApiKeyRepository, commonServices.PostgresRepositories.AppKeyRepository, security.FILE_STORE_API, security.WithCache(commonServices.Cache)),
		func(c *gin.Context) {
			ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreApi)

			publicUrl, err := services.FileService.GetFilePublicUrl(ctx, c.Param("id"))
			if err != nil && err.Error() != "record not found" {
				c.JSON(500, gin.H{"error": "Internal Server Error"})
				return
			}
			if err != nil && err.Error() == "record not found" {
				c.JSON(404, gin.H{"error": "File not found"})
				return
			}

			// return public url
			c.JSON(200, gin.H{"publicUrl": publicUrl})
		})

	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)

	r.GET("/jwt",
		tracing.TracingEnhancer(ctx, "GET /jwt"),
		security.TenantUserContextEnhancer(security.USERNAME, commonServices.Neo4jRepositories, security.WithCache(commonServices.Cache)),
		security.ApiKeyCheckerHTTP(commonServices.PostgresRepositories.TenantWebhookApiKeyRepository, commonServices.PostgresRepositories.AppKeyRepository, security.FILE_STORE_API, security.WithCache(commonServices.Cache)),
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

func MapFileEntityToDTO(cfg *config.Config, fileEntity *model.File) *fsc.FileDTO {
	return mapper.MapFileEntityToDTO(fileEntity, cfg.ApiServiceUrl)
}

func initLogger(cfg *config.Config) logger.Logger {
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName(constants.AppSourceFileStoreApi)
	return appLogger
}

func initTracing(cfg *config.Config, appLogger logger.Logger) io.Closer {
	if cfg.Jaeger.Enabled {
		tracer, closer, err := tracing.NewJaegerTracer(&cfg.Jaeger, appLogger)
		if err != nil {
			appLogger.Fatalf("Could not initialize jaeger tracer: %v", err.Error())
		}
		opentracing.SetGlobalTracer(tracer)
		return closer
	}
	return nil
}
