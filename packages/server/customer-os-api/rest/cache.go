package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
)

func CacheHandler(serviceContainer *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {

		tenant := c.Keys[security.KEY_TENANT_NAME].(string)

		organizations := serviceContainer.Caches.GetOrganizations(tenant)

		c.JSON(200, organizations)

	}
}
