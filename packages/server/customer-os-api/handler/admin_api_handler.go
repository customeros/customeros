package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"log"
	"net/http"
)

type AdminApiHandler struct {
	cfg            *config.Config
	commonServices *commonService.Services
}

func NewAdminApiHandler(config *config.Config, commonServices *commonService.Services) *AdminApiHandler {
	return &AdminApiHandler{
		cfg:            config,
		commonServices: commonServices,
	}
}

func (aah *AdminApiHandler) GetAdminApiHandlerEnhancer() func(c *gin.Context) {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(commonService.ApiKeyHeader)
		ctx := context.Background()

		tenant := c.GetHeader(commonService.TenantHeader)
		if tenant != "" {
			exists, err := aah.commonServices.CommonRepositories.TenantRepository.TenantExists(ctx, tenant)
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if !exists {
				log.Printf("Tenant %s does not exist", tenant)
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		c.Set("TenantName", tenant)
		if apiKey != aah.cfg.Admin.Key {
			log.Println("Invalid api key")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("Role", model.RoleAdmin)
		c.Next()
	}
}
