package main

import (
	"context"
	commonconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/service"
	"github.com/opentracing/opentracing-go"
	"io"

	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"log"
)

func main() {
	cfg := loadConfiguration()

	// Initialize Logging
	appLogger := initLogger(cfg)

	// Initialize Tracing
	tracingCloser := initTracing(cfg, appLogger)
	if tracingCloser != nil {
		defer tracingCloser.Close()
	}

	ctx := context.Background()

	// Initialize postgres db
	postgresDb, _ := InitDB(cfg, appLogger)
	defer postgresDb.SqlDB.Close()

	// Setting up Neo4j
	neo4jDriver, err := commonconf.NewNeo4jDriver(cfg.Neo4j)
	if err != nil {
		appLogger.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	services := service.InitServices(cfg, postgresDb.GormDB, &neo4jDriver, appLogger)

	r.POST("/validateAddress",
		handler.TracingEnhancer(ctx, "POST /validateAddress"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API),
		func(c *gin.Context) {
			var request model.ValidationAddressRequest
			if err := c.BindJSON(&request); err != nil {
				errorMessage := "Invalid request body"
				c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
				return
			}
			if request.International {
				internationalAddressLookup, err := services.AddressValidationService.ValidateInternationalAddress(request.Address, request.Country)
				if err != nil {
					errorMessage := err.Error()
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				if internationalAddressLookup == nil {
					errorMessage := "Invalid address"
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				addressVerified := false
				for _, v := range internationalAddressLookup.Results {
					if v.Analysis.VerificationStatus == "Verified" {
						addressVerified = true
						break
					}
				}

				if !addressVerified {
					errorMessage := "Address could not be verified"
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				c.JSON(200, model.MapValidationInternationalAddressResponse(internationalAddressLookup, nil, true))
			} else {
				validatedAddressResponse, err := services.AddressValidationService.ValidateUsAddress(request.Address)
				if err != nil {
					errorMessage := err.Error()
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				if validatedAddressResponse == nil {
					errorMessage := "Invalid address"
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				addressVerified := false
				for _, v := range validatedAddressResponse.Result.Addresses {
					if v.Verified {
						addressVerified = true
						break
					}
				}

				if !addressVerified {
					errorMessage := "Address could not be verified"
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				c.JSON(200, model.MapValidationUsAddressResponse(validatedAddressResponse, nil, true))
			}
		})

	r.POST("/validatePhoneNumber",
		handler.TracingEnhancer(ctx, "POST /validatePhoneNumber"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API),
		func(c *gin.Context) {
			var request model.ValidationPhoneNumberRequest

			if err := c.BindJSON(&request); err != nil {
				errorMessage := "Invalid request body"
				c.JSON(400, model.MapValidationPhoneNumberResponse(nil, nil, &errorMessage, false))
				return
			}

			e164, country, err := services.PhoneNumberValidationService.ValidatePhoneNumber(ctx, request.Country, request.PhoneNumber)
			if err != nil {
				errorMessage := err.Error()
				c.JSON(500, model.MapValidationPhoneNumberResponse(nil, nil, &errorMessage, false))
				return
			}

			if e164 == nil {
				errorMessage := "Invalid phone number"
				c.JSON(400, model.MapValidationPhoneNumberResponse(nil, nil, &errorMessage, false))
				return
			}

			c.JSON(200, model.MapValidationPhoneNumberResponse(e164, country, nil, true))
		})

	r.POST("/validateEmail",
		handler.TracingEnhancer(ctx, "POST /validateEmail"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API),
		func(c *gin.Context) {
			var request model.ValidationEmailRequest

			if err := c.BindJSON(&request); err != nil {
				appLogger.Errorf("Fail reading request: %v", err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			response, err := services.EmailValidationService.ValidateEmail(ctx, request.Email)
			if err != nil {
				appLogger.Errorf("Error validating email: %v", err.Error())
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, model.MapValidationEmailResponse(response, nil))
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

func initLogger(cfg *config.Config) logger.Logger {
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("VALIDATION-API")
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

func InitDB(cfg *config.Config, log logger.Logger) (db *commonconf.StorageDB, err error) {
	if db, err = commonconf.NewPostgresDBConn(cfg.Postgres); err != nil {
		log.Fatalf("Could not open db connection: %s", err.Error())
	}
	return
}
