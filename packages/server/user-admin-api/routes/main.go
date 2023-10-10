package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/routes/generate"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"log"
	"strings"
)

// Run will start the server
func Run(config *config.Config, services *service.Services) {
	router := getRouter(config, services)
	if err := router.Run(config.Service.ServerAddress); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func getRouter(config *config.Config, services *service.Services) *gin.Engine {
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

	addRegistrationRoutes(route, config, services)
	addSlackRoutes(route, config, services)
	generate.AddDemoTenantRoutes(route, config, services)

	route2 := router.Group("/")

	addHealthRoutes(route2)

	return router
}
