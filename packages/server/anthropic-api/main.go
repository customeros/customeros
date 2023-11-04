package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/anthorpic-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/anthorpic-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/anthorpic-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/sirupsen/logrus"
)

func InitDB(cfg *config.Config) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(cfg); err != nil {
		logrus.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func main() {
	cfg := loadConfiguration()
	config.InitLogger(cfg)

	db, _ := InitDB(cfg)
	defer db.SqlDB.Close()

	services := service.InitServices(cfg, db)

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.POST("/ask",
		commonservice.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonservice.ANTHROPIC_API, commonservice.WithCache(caches.NewCommonCache())),
		func(c *gin.Context) {
			var request dto.AnthropicApiRequest

			if err := c.BindJSON(&request); err != nil {
				logrus.Printf("Fail reading request: %v", err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			if request.Temperature == nil {
				i := 1
				request.Temperature = &i
			}
			if request.MaxTokensToSample == nil {
				i := 256
				request.MaxTokensToSample = &i
			}
			if request.StopSequences == nil {
				strings := []string{"\n\nHuman:"}
				request.StopSequences = &strings
			}

			anthropicResponse := services.AnthropicService.QueryAnthropic(request)
			if anthropicResponse.Error != nil {
				c.JSON(500, gin.H{"error": anthropicResponse})
				return
			}

			c.JSON(200, anthropicResponse)
		})

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
