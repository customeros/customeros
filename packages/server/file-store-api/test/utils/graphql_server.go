package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/test/graph"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/test/graph/resolver"

	"log"
	"net/http/httptest"
	"testing"
)

func NewWebServer(t *testing.T) (*httptest.Server, *graphql.Client, *resolver.Resolver) {
	router := gin.Default()
	server := httptest.NewServer(router)
	handler, resolver := graph.GraphqlHandler()
	router.POST("/query", handler)
	graphqlClient := graphql.NewClient(server.URL + "/query")

	t.Cleanup(func() {
		log.Printf("shutting down webserver")
		server.Close()
	})
	return server, graphqlClient, resolver
}
