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
		if apiKey != aah.cfg.Admin.Key {
			log.Println("Invalid api key")
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []gin.H{{"message": "Invalid api key"}},
			})
			c.Abort()
			return
		}

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
		c.Set(commonService.KEY_TENANT_NAME, tenant)

		//TODO DROP THIS. WE NEED TO USE THE TenantUserContextEnhancer + a check on the ADMIN ROLE
		usernameHeader := c.GetHeader(commonService.UsernameHeader)
		if usernameHeader != "" {
			userId, tenantName, roles, err := aah.commonServices.CommonRepositories.UserRepository.FindUserByEmail(ctx, usernameHeader)
			if err != nil {
				log.Printf("Error checking user existence: %s", err.Error())
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message": fmt.Sprintf("Error checking user existence: %s", err.Error())}},
				})
				c.Abort()
				return
			}
			if tenant != tenantName {
				log.Printf("User %s does not belong to tenant %s", usernameHeader, tenant)
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message": fmt.Sprintf("User %s does not belong to tenant %s", usernameHeader, tenant)}},
				})
				c.Abort()
				return
			}
			c.Set(commonService.KEY_USER_ID, userId)
			c.Set(commonService.KEY_USER_ROLES, roles)
		} else {
			c.Set(commonService.KEY_USER_ROLES, []string{model.RoleAdmin.String()})
		}

		c.Next()
	}
}
