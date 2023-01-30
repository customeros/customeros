package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/resolver"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service/container"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const customerOSApiPort = "10000"

func InitDB(cfg *config.Config) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(cfg); err != nil {
		logrus.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func graphqlHandler(cfg *config.Config, driver neo4j.Driver, repositoryContainer *commonRepository.Repositories) gin.HandlerFunc {
	serviceContainer := container.InitServices(&driver)
	// instantiate graph resolver
	graphResolver := resolver.NewResolver(serviceContainer, repositoryContainer)

	schemaConfig := generated.Config{Resolvers: graphResolver}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(schemaConfig))
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return gqlerror.Errorf("Internal server error!")
	})
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		// Error hook place, Returned error can be customized. Check https://gqlgen.com/reference/errors/
		return err
	})
	srv.Use(extension.FixedComplexityLimit(cfg.GraphQL.FixedComplexityLimit))

	return func(c *gin.Context) {
		customCtx := &common.CustomContext{
			Tenant: c.Keys["TenantName"].(string),
			UserId: c.Keys["UserId"].(string),
		}
		h := common.WithContext(customCtx, srv)

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func init() {
}

func main() {
	cfg := loadConfiguration()
	config.InitLogger(cfg)

	db, _ := InitDB(cfg)
	defer db.SqlDB.Close()

	neo4jDriver, err := config.NewDriver(cfg)
	if err != nil {
		logrus.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close()

	repositoryContainer := commonRepository.InitRepositories(db.GormDB, &neo4jDriver)

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.POST("/query",
		commonService.UserToTenantEnhancer(repositoryContainer.UserRepo),
		commonService.ApiKeyCheckerHTTP(repositoryContainer.AppKeyRepo, commonService.CUSTOMER_OS_API),
		graphqlHandler(cfg, neo4jDriver, repositoryContainer))
	if cfg.GraphQL.PlaygroundEnabled {
		r.GET("/", playgroundHandler())
	}
	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)

	port := cfg.ApiPort
	if port == "" {
		port = customerOSApiPort
	}

	r.Run(":" + port)
}

func loadConfiguration() *config.Config {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Error loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("%+v\n", err)
	}

	return &cfg
}

func healthCheckHandler(c *gin.Context) {

	c.JSON(200, gin.H{"status": "OK"})
}
