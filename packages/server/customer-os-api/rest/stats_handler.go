package rest_handlers

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/gin-gonic/gin"

	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

func StatsSuccessHandler(api string, services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call your actual handler function
		c.Next()

		tenant := common.GetTenantFromContext(c.Request.Context())

		// Check the response status code
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			// Increment the API call stat
			_, err := services.Repositories.PostgresRepositories.StatsApiCallsRepository.Increment(c.Request.Context(), tenant, api)
			if err != nil {
				services.Log.Errorf("Error incrementing API calls stats %s for tenant: %s", api, tenant, err.Error())
			}
		}
	}
}
