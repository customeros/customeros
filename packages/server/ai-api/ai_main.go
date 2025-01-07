package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/service"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/sirupsen/logrus"
)

const defaultAnthropicModel = "claude-3-haiku-20240307"

func main() {
	cfg := loadConfiguration()
	config.InitLogger(cfg)

	postgresDb, err := commonConfig.InitPostgres(&commonConfig.GlobalConfig{
		PostgresConfig:      &cfg.PostgresConfig,
		PostgresAsyncConfig: &cfg.PostgresAsyncConfig,
	})
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	services := service.InitServices(cfg, postgresDb, logger.NewAppLogger(nil))

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.POST("/askAI",
		security.ApiKeyCheckerHTTP(
			services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository,
			services.CommonServices.PostgresRepositories.AppKeyRepository,
			security.AI_API,
			security.WithCache(services.CommonServices.Cache),
		),
		routes.AskAI(services),
	)

	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)

	r.Run(":" + cfg.ApiPort)
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func loadConfiguration() *config.Config {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Failed loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("%+v", err)
	}

	return &cfg
}
