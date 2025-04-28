package private

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func RequestAccessQuickbooks(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, _ := telemetry.StartRestSpan(c.Request.Context(), "RequestAccessQuickbooks")
		defer spans.Finish()

		quickbooksRequestAccessUrl := "https://appcenter.intuit.com/connect/oauth2?client_id=" + s.Cfg.Common.External.QuickbooksConfig.ClientId

		redirectUrl := c.Query("redirect_url")
		if redirectUrl != "" {
			quickbooksRequestAccessUrl += "&redirect_uri=" + url.QueryEscape(redirectUrl)
		}

		quickbooksRequestAccessUrl += "&response_type=code&scope=com.intuit.quickbooks.accounting"

		state := c.Query("state")
		if state != "" {
			quickbooksRequestAccessUrl += "&state=" + state
		}

		spans.LogKV("quickbooksRequestAccessUrl", quickbooksRequestAccessUrl)

		c.JSON(http.StatusOK, gin.H{"url": quickbooksRequestAccessUrl})
	}
}

func CallbackQuickbooks(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "CallbackQuickbooks")
		defer spans.Finish()

		code := c.Request.URL.Query().Get("code")
		realmId := c.Request.URL.Query().Get("realmId")
		redirectUrl := c.Request.URL.Query().Get("redirect_url")
		spans.LogKV("code", code)
		spans.LogKV("redirectUrl", url.QueryEscape(redirectUrl))
		spans.LogKV("realmId", realmId)

		requestData := url.Values{}
		requestData.Set("grant_type", "authorization_code")
		requestData.Set("code", code)
		requestData.Set("redirect_uri", redirectUrl)

		tenant, _ := c.Get(security.KEY_TENANT_NAME)

		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant: tenant.(string),
		})

		_, err := s.CommonServices.QuickbooksService.GetAndStoreAccessToken(ctx, realmId, requestData)
		if err != nil {
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		//publish SKUs not pushed already to QB on initial connection
		skuList, err := s.CommonServices.PostgresRepositories.SkuRepository.GetAll(ctx, tenant.(string), nil)
		if err != nil {
			spans.TraceError(err)
		}

		if skuList != nil && len(skuList) > 0 {
			for _, sku := range skuList {
				if sku.QuickbooksId == "" {
					err = s.CommonServices.Events.Publisher.PublishFanoutEvent(ctx, sku.ID, model.SKU, dto.SkuUpdate{})
					if err != nil {
						spans.TraceError(err)
					}
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func RevokeQuickbooks(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "RevokeQuickbooks")
		defer spans.Finish()

		// Retrieve the tenant from context.
		tenant, exists := c.Get(security.KEY_TENANT_NAME)
		if !exists || tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found in context"})
			return
		}

		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant.(string),
			AppSource: constants.AppSourceCustomerOsApi,
		})
		spans.LogKV("tenant", tenant.(string))

		// Call the QuickBooks service to revoke the access.
		err := s.CommonServices.QuickbooksService.RevokeAccess(ctx)
		if err != nil {
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke QuickBooks access"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func GetQuickbooksSettings(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "GetQuickbooksSettings")
		defer spans.Finish()

		tenant, _ := c.Get(security.KEY_TENANT_NAME)

		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant: tenant.(string),
		})

		quickbooksConnected, err := s.CommonServices.QuickbooksService.QuickbooksConnected(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"quickbooksConnected": quickbooksConnected})
	}
}
