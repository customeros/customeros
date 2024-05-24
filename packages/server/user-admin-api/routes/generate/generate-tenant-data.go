package generate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"io/ioutil"
	"net/http"
)

type CleanupRequest struct {
	Tenant        string `json:"tenant"`
	ConfirmTenant string `json:"confirmTenant"`
}

func AddDemoTenantRoutes(rg *gin.RouterGroup, config *config.Config, services *service.Services) {
	rg.POST("/demo-tenant-data", func(context *gin.Context) {

		apiKey := context.GetHeader("X-Openline-Api-Key")
		if apiKey != config.Service.ApiKey {
			context.JSON(http.StatusUnauthorized, gin.H{
				"result": fmt.Sprintf("invalid api key"),
			})
			return
		}

		sourceData, err := validateRequestAndGetFileBytes(context)
		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		tenant := context.GetHeader("TENANT_NAME")
		username := context.GetHeader("MASTER_USERNAME")

		err = services.TenantDataInjector.InjectTenantData(context, tenant, username, sourceData)
		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.JSON(200, gin.H{
			"tenant": "tenant initiated",
		})
	})

	rg.POST("/demo-tenant-delete", func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(context.Background(), "GET /demo-tenant-delete", c.Request.Header)
		defer span.Finish()

		apiKey := c.GetHeader("X-Openline-Api-Key")
		if apiKey != config.Service.ApiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"result": fmt.Sprintf("invalid api key"),
			})
			return
		}

		tenant := c.GetHeader("TENANT_NAME")
		username := c.GetHeader("MASTER_USERNAME")

		var req CleanupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"result": "invalid tenant delete payload",
			})
			return
		}
		err := services.TenantDataInjector.CleanupTenantData(ctx, tenant, username, req.Tenant, req.ConfirmTenant)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"tenant": "tenant " + tenant + " successfully cleaned up",
		})
	})
}

func validateRequestAndGetFileBytes(context *gin.Context) (*service.SourceData, error) {
	tenant := context.GetHeader("TENANT_NAME")
	if tenant == "" {
		return nil, errors.New("tenant is required")
	}

	username := context.GetHeader("MASTER_USERNAME")
	if username == "" {
		return nil, errors.New("username is required")
	}

	multipartFileHeader, err := context.FormFile("file")
	if err != nil {
		return nil, err
	}

	multipartFile, err := multipartFileHeader.Open()
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(multipartFile)
	if err != nil {
		return nil, err
	}

	var sourceData service.SourceData
	if err := json.Unmarshal(bytes, &sourceData); err != nil {
		return nil, err
	}

	return &sourceData, nil
}
