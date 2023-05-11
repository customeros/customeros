package handler

import (
	"context"
	"fmt"
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
				log.Printf("Error checking tenant existence: %s", err.Error())
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message": fmt.Sprintf("Error checking tenant existence: %s", err.Error())}},
				})
				c.Abort()
				return
			}
			if !exists {
				log.Printf("Tenant %s does not exist", tenant)
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message": fmt.Sprintf("Tenant %s does not exist", tenant)}},
				})
				c.Abort()
				return
			}
		}
		c.Set("TenantName", tenant)
		if apiKey != aah.cfg.Admin.Key {
			log.Println("Invalid api key")
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []gin.H{{"message": "Invalid api key"}},
			})
			c.Abort()
			return
		}
		c.Set("Role", model.RoleAdmin)
		c.Next()
	}
}
