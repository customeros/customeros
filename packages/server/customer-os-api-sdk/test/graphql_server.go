package test

import (
	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/resolver"
	"log"
	"net/http/httptest"
	"testing"
)

func NewGraphQlMockedServer(t *testing.T) (*httptest.Server, *graphql.Client, *resolver.Resolver) {
	router := gin.Default()
	server := httptest.NewServer(router)
	handler, r := GraphServerMocked()
	router.POST("/query", handler)
	graphqlClient := graphql.NewClient(server.URL + "/query")

	t.Cleanup(func() {
		log.Printf("Shutting down webserver")
		server.Close()
	})
	return server, graphqlClient, r
}
