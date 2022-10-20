package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
)

const customerOSApiPort = "1010"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = customerOSApiPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
