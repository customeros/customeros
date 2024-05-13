package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/ContactHub"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"log"
)

// Run will start the server
func Run(config *c.Config, hub *ContactHub.ContactHub, services *service.Services) {
	router := getRouter(config, hub, services)
	if err := router.Run(":" + config.Service.Port); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func getRouter(config *c.Config, hub *ContactHub.ContactHub, services *service.Services) *gin.Engine {
	router := gin.New()
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"*"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS", "POST", "GET")
	corsConfig.AddAllowHeaders("WebChatApiKey")

	router.Use(cors.New(corsConfig))
	route := router.Group("/")

	addMailRoutes(config, route, services, hub)

	AddCalComRoutes(route, services.CustomerOsService, services.CommonServices.PostgresRepositories.PersonalIntegrationRepository)
	AddQueryRoutes(route, services.CustomerOsService, services.RedisService)

	addHealthRoutes(route)
	return router
}
