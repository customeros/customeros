package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"

	"log"
	"strings"
)

// Run will start the server
func Run(config *config.Config, graphqlClient *graphql.Client) {
	router := getRouter(config)
	if err := router.Run(config.Service.ServerAddress); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func getRouter(config *config.Config) *gin.Engine {
	router := gin.New()
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = strings.Split(config.Service.CorsUrl, " ")
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS", "POST", "GET")
	corsConfig.AddAllowHeaders("WebChatApiKey")

	router.Use(cors.New(corsConfig))
	route := router.Group("/")

	addRegistrationRoutes(route, config)

	route2 := router.Group("/")

	addHealthRoutes(route2)
	return router
}
