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
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service/container"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"log"
)

const customerOSApiPort = "10000"

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

func graphqlHandler(driver neo4j.Driver, repositoryContainer *repository.PostgresRepositoryContainer) gin.HandlerFunc {
	serviceContainer := container.InitServices(&driver)
	// instantiate graph resolver
	graphResolver := resolver.NewResolver(serviceContainer, repositoryContainer)

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

// Declare a simple handler for pingpong as a request accepting behavior
func ApiKeyChecker(appKeyRepo repository.AppKeyRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		kh := c.GetHeader("X-Openline-API-KEY")
		if kh != "" {

			keyResult := appKeyRepo.FindByKey(c, kh)

			if keyResult.Error != nil {
				c.AbortWithStatus(401)
				return
			}

			appKey := keyResult.Result.(*entity.AppKeyEntity)

			if appKey == nil {
				c.AbortWithStatus(401)
				return
			} else {
				// todo set tenant in context
			}

			c.Next()
			// illegal request, terminate the current process
		} else {
			c.AbortWithStatus(401)
			return
		}

	}
}

func main() {
	cfg := loadConfiguration()

	db, _ := InitDB(cfg)
	defer db.SqlDB.Close()

	repositoryContainer := repository.InitRepositories(db.GormDB)

	neo4jDriver, err := config.NewDriver(cfg)
	if err != nil {
		log.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close()

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.POST("/query", ApiKeyChecker(repositoryContainer.AppKeyRepo), graphqlHandler(neo4jDriver, repositoryContainer))
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
