package service

import (
	"github.com/gin-gonic/gin"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
)

func UserToTenantEnhancer(userToTenantRepository repository.UserToTenantRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		uh := c.GetHeader("X-Openline-USERNAME")
		if uh != "" {

			tenantResult := userToTenantRepository.FindTenantByUsername(c, uh)

			if tenantResult.Error != nil {
				c.AbortWithStatus(401)
				return
			}

			tenant := tenantResult.Result.(string)

			if len(tenant) == 0 {
				c.AbortWithStatus(401)
				return
			} else {
				if c.Keys == nil {
					c.Keys = map[string]any{}
				}
				c.Keys["TenantName"] = tenant
			}

			c.Next()
			// illegal request, terminate the current process
		} else {
			c.AbortWithStatus(401)
			return
		}

	}
}
