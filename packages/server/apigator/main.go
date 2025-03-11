package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/customeros/customeros/packages/server/apigator/config"
	"github.com/customeros/customeros/packages/server/apigator/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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

	cache := caches.NewCommonCache()
	postgresRepositories := postgresRepository.InitRepositories(postgresDb)
	neo4jRepositories := neo4jRepository.InitNeo4jRepositories(&neo4jDriver, cfg.Neo4jConfig.Database)

	r := gin.Default()

	r.Use(tracing.RecoveryWithJaeger(opentracing.GlobalTracer()))

	r.GET("/validate", validateToken(
		postgresRepositories.TenantWebhookApiKeyRepository,
		cfg.AppKey,
		neo4jRepositories,
		"apigator",
		security.WithCache(cache)))

	log.Printf("Auth service running on port %s", cfg.ApiPort)
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

func validateToken(
	tenantApiKeyRepo postgresRepository.TenantWebhookApiKeyRepository,
	appKey string,
	repos *neo4jRepository.Repositories,
	app security.App,
	opts ...security.CommonServiceOption,
) gin.HandlerFunc {

	return func(c *gin.Context) {
		if isIntrospectionQuery(c.Request) {
			c.Next()
			return
		}

		security.ApiKeyCheckerHTTP(tenantApiKeyRepo, appKey, opts...)(c)
		security.TenantUserContextEnhancer(repos, opts...)(c)
	}
}
