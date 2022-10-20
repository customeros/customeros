package main

import (
	common "github.com.openline-ai.customer-os-analytics-api/common"
	"github.com.openline-ai.customer-os-analytics-api/config"
	"github.com.openline-ai.customer-os-analytics-api/repository"
	"log"
	"net/http"
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, _ := InitDB()
	defer db.SqlDB.Close()

	repositoryHandler := repository.InitRepositories(db.GormDB)

	customCtx := &common.CustomContext{}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		RepositoryHandler: repositoryHandler,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", common.CreateContext(customCtx, srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
