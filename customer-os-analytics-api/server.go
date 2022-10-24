package main

import (
	common "github.com.openline-ai.customer-os-analytics-api/common"
	"github.com.openline-ai.customer-os-analytics-api/config"
	"github.com.openline-ai.customer-os-analytics-api/dataloader"
	"github.com.openline-ai.customer-os-analytics-api/graph/resolver"
	"github.com.openline-ai.customer-os-analytics-api/repository"
	"github.com/gin-gonic/gin"
	"log"
	"os"

	"github.com.openline-ai.customer-os-analytics-api/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func InitDB() (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PWD"),
		100,
		10,
		0); err != nil {
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
	db, _ := InitDB()
	defer db.SqlDB.Close()

	r := gin.Default()
	r.POST("/query", graphqlHandler(db))
	r.GET("/", playgroundHandler())
	r.Run()
}
