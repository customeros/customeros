package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"log"
	"strings"
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

	corsConfig.AllowOrigins = strings.Split(config.Service.CorsUrl, " ")
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

	//tracking configuration
	//all all sources and filter down by tenant information in database
	trackingRoute := router.Group("/tracking")

	trackingCorsConfig := cors.DefaultConfig()
	trackingCorsConfig.AllowOrigins = []string{"*"}
	trackingCorsConfig.AllowHeaders = []string{"Host", "Content-Type", "Content-Length", "Accept", "Origin", "Referer", "User-Agent"}
	trackingCorsConfig.AllowCredentials = false
	trackingCorsConfig.AddAllowMethods("OPTIONS", "POST")

	trackingRoute.Use(cors.New(trackingCorsConfig))

	addTrackingRoutes(trackingRoute, services)

	return router
}
