package routes

import (
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"customeros/customeros/packages/server/user-admin-api/config"
	"customeros/customeros/packages/server/user-admin-api/service"
)

// Run will start the server
func Run(config *config.Config, services *service.Services) {
	router := getRouter(config, services)
	if err := router.Run(":" + config.Service.Port); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func getRouter(config *config.Config, services *service.Services) *gin.Engine {
	router := gin.New()
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = strings.Split(config.Service.CorsURL, " ")
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS", "POST", "GET")
	corsConfig.AddAllowHeaders("WebChatApiKey")
	corsConfig.AddAllowHeaders("X-Tracker-Payload")
	corsConfig.AddAllowHeaders("TENANT_NAME")
	corsConfig.AddAllowHeaders("MASTER_USERNAME")

	router.Use(cors.New(corsConfig))
	route := router.Group("/")

	addRegistrationRoutes(route, config, services)
	addSlackRoutes(route, config, services)
	addMailRoutes(route, config, services)

	addHealthRoutes(route)

	return router
}
