package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/resolver"
	"github.com/openline-ai/openline-customer-os/customer-os-api/service"
	"log"
)

const customerOSApiPort = "1010"

func graphqlHandler(cfg *config.Config) gin.HandlerFunc {

	serviceContainer := service.InitServices(cfg)
	// instantiate graph resolver
	graphResolver := resolver.NewResolver(serviceContainer)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphResolver}))

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

func main() {
	cfg := loadConfiguration()
	// Setting up Gin
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))

	r.POST("/query", graphqlHandler(cfg))
	r.GET("/", playgroundHandler())
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
