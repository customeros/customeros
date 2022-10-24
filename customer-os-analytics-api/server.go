package main

import (
	common "github.com.openline-ai.customer-os-analytics-api/common"
	"github.com.openline-ai.customer-os-analytics-api/config"
	"github.com.openline-ai.customer-os-analytics-api/dataloader"
	"github.com.openline-ai.customer-os-analytics-api/graph/resolver"
	"github.com.openline-ai.customer-os-analytics-api/repository"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com.openline-ai.customer-os-analytics-api/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultApiPort = "8080"

type Config struct {
	Db struct {
		Host            string `env:"DB_HOST,required"`
		Port            string `env:"DB_PORT" envDefault:"5432"`
		Pwd             string `env:"DB_PWD,unset"`
		Name            string `env:"DB_NAME,required"`
		User            string `env:"DB_USER,required"`
		MaxConn         int    `env:"DB_MAX_CONN"`
		MaxIdleConn     int    `env:"DB_MAX_IDLE_CONN"`
		ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME"`
	}
}

func InitDB(cfg *Config) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.Name,
		cfg.Db.User,
		cfg.Db.Pwd,
		cfg.Db.MaxConn,
		cfg.Db.MaxIdleConn,
		cfg.Db.ConnMaxLifetime); err != nil {
		log.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

// Defining the Graphql handler
func graphqlHandler(db *config.StorageDB) gin.HandlerFunc {
	// instantiate repository handler
	repositoryContainer := repository.InitRepositories(db.GormDB)
	// instantiate graph resolver
	graphResolver := resolver.NewResolver(repositoryContainer)
	// make a data loader
	loader := dataloader.NewDataLoader(repositoryContainer)
	// make a custom context
	customCtx := &common.CustomContext{
		Tenant: "openline", // FIXME alexb replace with authentication
	}
	// create query handler
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphResolver}))
	// wrap the query handler with middleware to inject dataloader
	dataloaderSrv := dataloader.Middleware(loader, srv)

	h := common.CreateContext(customCtx, dataloaderSrv)

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

	db, _ := InitDB(cfg)
	defer db.SqlDB.Close()

	r := gin.Default()
	r.POST("/query", graphqlHandler(db))
	r.GET("/", playgroundHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultApiPort
	}
	r.Run(":" + port)
}

func loadConfiguration() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] Error loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}
