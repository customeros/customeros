package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/dto"
	"github.com/sirupsen/logrus"

	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/service"
	"log"
)

func InitDB(cfg *config.Config) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(cfg); err != nil {
		logrus.Fatalf("Coud not open db connection: %s", err.Error())
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

	services := service.InitServices(cfg, db, &neo4jDriver)

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.POST("/validateAddress",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME_OR_TENANT, services.CommonServices.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.VALIDATION_API),
		func(c *gin.Context) {
			var request dto.ValidationAddressRequest

			if err := c.BindJSON(&request); err != nil {
				errorMessage := "Invalid request body"
				c.JSON(400, dto.MapValidationAddressResponse(nil, &errorMessage, false))
				return
			}

			validatedAddressResponse, err := services.AddressValidationService.ValidateAddress(request.Address)
			if err != nil {
				errorMessage := err.Error()
				c.JSON(400, dto.MapValidationAddressResponse(nil, &errorMessage, false))
				return
			}

			if validatedAddressResponse == nil {
				errorMessage := "Invalid address"
				c.JSON(400, dto.MapValidationAddressResponse(nil, &errorMessage, false))
				return
			}

			c.JSON(200, dto.MapValidationAddressResponse(validatedAddressResponse, nil, true))
		})

	r.POST("/validatePhoneNumber",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME_OR_TENANT, services.CommonServices.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.VALIDATION_API),
		func(c *gin.Context) {
			var request dto.ValidationPhoneNumberRequest

			if err := c.BindJSON(&request); err != nil {
				errorMessage := "Invalid request body"
				c.JSON(400, dto.MapValidationPhoneNumberResponse(nil, nil, &errorMessage, false))
				return
			}

			e164, country, err := services.PhoneNumberValidationService.ValidatePhoneNumber(ctx, request.Country, request.PhoneNumber)
			if err != nil {
				errorMessage := err.Error()
				c.JSON(500, dto.MapValidationPhoneNumberResponse(nil, nil, &errorMessage, false))
				return
			}

			if e164 == nil {
				errorMessage := "Invalid phone number"
				c.JSON(400, dto.MapValidationPhoneNumberResponse(nil, nil, &errorMessage, false))
				return
			}

			c.JSON(200, dto.MapValidationPhoneNumberResponse(e164, country, nil, true))
		})

	r.POST("/validateEmail",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME_OR_TENANT, services.CommonServices.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.VALIDATION_API),
		func(c *gin.Context) {
			var request dto.ValidationEmailRequest

			if err := c.BindJSON(&request); err != nil {
				println(err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			response, err := services.EmailValidationService.ValidateEmail(ctx, request.Email)
			if err != nil {
				errorMessage := err.Error()
				c.JSON(500, dto.MapValidationEmailResponse(nil, &errorMessage, false))
				return
			}

			c.JSON(200, dto.MapValidationEmailResponse(response, nil, true))
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
