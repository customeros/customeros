package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/resolver"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service/container"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"log"
	"strings"
)

const customerOSApiPort = "10000"

func graphqlHandler(driver neo4j.Driver) gin.HandlerFunc {
	serviceContainer := container.InitServices(&driver)
	// instantiate graph resolver
	graphResolver := resolver.NewResolver(serviceContainer)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphResolver}))

	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return gqlerror.Errorf("Internal server error!")
	})
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		// Error hook place, Returned error can be customized. Check https://gqlgen.com/reference/errors/
		return err
	})

	customCtx := &common.CustomContext{
		Tenant: "openline", // TODO replace with tenant from authentication
	}
	h := common.CreateContext(customCtx, srv)

	return func(c *gin.Context) {
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
	logger.Logger = logger.New(log.New(log.Default().Writer(), "", log.Ldate|log.Ltime|log.Lmicroseconds), logger.Config{
		Colorful: true,
		LogLevel: logger.Info,
	})
}

func main() {
	cfg := loadConfiguration()
	neo4jDriver, err := config.NewDriver(cfg)
	if err != nil {
		log.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close()

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = strings.Split(cfg.CorsUrl, " ")
	r.Use(cors.New(corsConfig))

	r.POST("/query", graphqlHandler(neo4jDriver))
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
