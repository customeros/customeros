package main

import (
	common "github.com.openline-ai.customer-os-analytics-api/common"
	"github.com.openline-ai.customer-os-analytics-api/config"
	"github.com.openline-ai.customer-os-analytics-api/repository"
	"github.com/gin-gonic/gin"
	"log"
	"os"

	"github.com.openline-ai.customer-os-analytics-api/graph"
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
	repositoryHandler := repository.InitRepositories(db.GormDB)

	customCtx := &common.CustomContext{}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		RepositoryHandler: repositoryHandler,
	}}))

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

func main() {
	db, _ := InitDB()
	defer db.SqlDB.Close()

	r := gin.Default()
	r.POST("/query", graphqlHandler(db))
	r.GET("/", playgroundHandler())
	r.Run()
}
