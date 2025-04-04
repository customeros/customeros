package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/customeros/customeros/packages/server/apigator/config"
	entities "github.com/customeros/customeros/packages/server/apigator/entity"
	"github.com/customeros/customeros/packages/server/apigator/logger"
	apigator_service "github.com/customeros/customeros/packages/server/apigator/service"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	utils "github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jRepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := loadConfiguration()
	ctx := context.Background()
	appLogger := initLogger(cfg)

	tracingCloser := initTracing(cfg, appLogger)
	if tracingCloser != nil {
		defer tracingCloser.Close()
	}
	defer tracing.RecoverAndLogToJaeger(appLogger)

	neo4jDriver, err := commonConfig.NewNeo4jDriver(cfg.Neo4jConfig)
	if err != nil {
		logrus.Fatalf(
			"Could not establish connection with neo4j at: %v, error: %v",
			cfg.Neo4jConfig.Target,
			err.Error(),
		)
	}
	defer neo4jDriver.Close(ctx)

	postgresDb, err := commonConfig.InitPostgres(&commonConfig.CommonConfig{
		Infrastructure: commonConfig.InfrastructureConfig{
			PostgresConfig:      cfg.PostgresConfig,
			PostgresAsyncConfig: cfg.PostgresAsyncConfig,
		},
	})
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	postgresRepos := postgresRepository.InitRepositories(postgresDb)
	neo4jRepos := neo4jRepository.InitNeo4jRepositories(&neo4jDriver, cfg.Neo4jConfig.Database)

	service := apigator_service.Service{}
	service.Init(
		postgresRepos.TenantWebhookApiKeyRepository,
		neo4jRepos.UserReadRepository,
	)

	r := gin.Default()

	r.Use(tracing.RecoveryWithJaeger(opentracing.GlobalTracer(), appLogger))
	r.GET("/validate", validate(&service, cfg.AppKey))

	log.Printf("Apigator running on port %s", cfg.ApiPort)
	r.Run(cfg.ApiPort)
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

func initLogger(cfg *config.Config) logger.Logger {
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("APIGATOR")
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

func isIntrospectionQuery(req *http.Request) bool {
	var requestMap map[string]interface{}
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		return false
	}

	req.Body = io.NopCloser(bytes.NewReader(requestBody))

	if err = json.Unmarshal(requestBody, &requestMap); err != nil {
		return false
	}

	if opName, ok := requestMap["operationName"].(string); ok && opName == "IntrospectionQuery" {
		// Check if "__schema" is present in the request
		if selectionSet, ok := requestMap["query"].(string); ok {
			if strings.Contains(selectionSet, "__schema {") && strings.Contains(selectionSet, "query IntrospectionQuery") {
				return true
			}
		}
	}
	return false
}

const (
	INTERNAL_API_KEY_HEADER = "X-OPENLINE-API-KEY"
	TENANT_API_KEY_HEADER   = "X-CUSTOMER-OS-API-KEY"
	USERNAME_HEADER         = "X-OPENLINE-USERNAME"
	TENANT_HEADER           = "X-OPENLINE-TENANT"
)

func validate(
	service *apigator_service.Service,
	appKey string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "validate")
		service.SetContext(ctx)

		traceErr := ""
		defer func() {
			if traceErr != "" {
				tracing.TraceErr(span, errors.New(traceErr))
			}
			span.Finish()
		}()

		// ✅ Skip auth entirely for WebSocket upgrade requests
		if strings.EqualFold(c.GetHeader("Upgrade"), "websocket") {
			c.Status(http.StatusOK)
			return
		}

		if isIntrospectionQuery(c.Request) {
			c.Next()
			return
		}

		var response ApigatorResponse
		response.RequestId = utils.GenerateNanoIdWithPrefix("api", 16)

		internalApiKey := c.GetHeader(INTERNAL_API_KEY_HEADER)
		tenantApiKey := c.GetHeader(TENANT_API_KEY_HEADER)
		username := c.GetHeader(USERNAME_HEADER)
		tenant := c.GetHeader(TENANT_HEADER)

		var userDetails *entities.UserDetails

		span.LogKV(
			INTERNAL_API_KEY_HEADER, internalApiKey,
			TENANT_API_KEY_HEADER, tenantApiKey,
			USERNAME_HEADER, username,
			TENANT_HEADER, tenant,
		)

		if internalApiKey == "" && tenantApiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Payload("error", "Missing API key"))
			return
		}

		/**
		 * Validate internal service API key.
		 * Username header is mandatory
		 * Tenant header is optional
		 * If Service API Key is present, Tenant API Key will be ignored
		 */
		if internalApiKey != "" {
			if internalApiKey != appKey {
				traceErr = "Invalid API Key"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Payload("error", traceErr))
				return
			}
			if username == "" {
				traceErr = "Missing Username header"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Payload("error", traceErr))
				return
			}

			foundTenant, err := service.GetTenantByUser(username)
			if err != nil || foundTenant == "" {
				traceErr = "Failed to authenticate user"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Payload("error", traceErr))
				c.Abort()
				return
			}

			if tenant != "" && tenant != foundTenant {
				traceErr = "Invalid tenant"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Payload("error", traceErr))
				return
			} else {
				tenant = foundTenant
			}

			foundUser, err := service.GetUserDetails(tenant, username)
			if err != nil {
				traceErr = "User not found"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Payload("error", traceErr))
				return
			}

			userDetails = foundUser
		} else if tenantApiKey != "" {
			/**
			* Validate tenant API key.
			* Both username header and tenant header are optional.
			* If username header is present, validate it against the tenant & decorate the user details.
			* If tenant header is present, validate it against the tenant associated with the API key.
			 */
			foundTenant, err := service.GetTenantByApiKey(tenantApiKey)
			if err != nil || foundTenant == "" {
				traceErr = "Invalid API key"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Payload("error", traceErr))
				c.Abort()
				return
			}

			if tenant != "" && tenant != foundTenant {
				traceErr = "Invalid tenant"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Payload("error", traceErr))
				c.Abort()
				return
			} else {
				tenant = foundTenant
			}

			if username != "" {
				foundUser, err := service.GetUserDetails(foundTenant, username)

				if err != nil || foundUser == nil {
					traceErr = "User not found"
					c.AbortWithStatusJSON(http.StatusUnauthorized, response.Payload("error", traceErr))
					c.Abort()
					return
				}

				userDetails = foundUser
			}
		}

		c.Header("X-Tenant", tenant)
		c.Header("X-Username", username)
		c.Header("X-Request-Id", response.RequestId)
		if userDetails != nil {
			userDetails.ToHeaders(c)
		}

		c.Next()
	}
}

type ApigatorResponse struct {
	RequestId string `json:"requestId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

func (instance *ApigatorResponse) Payload(status, message string) *ApigatorResponse {
	instance.Status = status
	instance.Message = message
	return instance
}
