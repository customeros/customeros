package private

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
	"net/url"
)

func RequestAccessQuickbooks(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/quickbooks/requestAccess", c.Request.Header)
		defer span.Finish()

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

		span.LogFields(log.Object("quickbooksRequestAccessUrl", quickbooksRequestAccessUrl))

		c.JSON(http.StatusOK, gin.H{"url": quickbooksRequestAccessUrl})
	}
}

func CallbackQuickbooks(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/quickbooks/oauth/callback", c.Request.Header)
		defer span.Finish()

		code := c.Request.URL.Query().Get("code")
		realmId := c.Request.URL.Query().Get("realmId")
		redirectUrl := c.Request.URL.Query().Get("redirect_url")
		span.LogKV("code", code)
		span.LogKV("redirectUrl", url.QueryEscape(redirectUrl))
		span.LogKV("realmId", realmId)

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
			span.LogFields(log.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		//publish SKUs not pushed already to QB on initial connection
		skuList, err := s.CommonServices.PostgresRepositories.SkuRepository.GetAll(ctx, tenant.(string), nil)
		if err != nil {
			tracing.TraceErr(span, err)
		}

		if skuList != nil && len(skuList) > 0 {
			for _, sku := range skuList {
				if sku.QuickbooksId == "" {
					err = s.CommonServices.Events.Publisher.PublishFanoutEvent(ctx, sku.ID, model.SKU, dto.SkuUpdate{})
					if err != nil {
						tracing.TraceErr(span, err)
					}
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func RevokeQuickbooks(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start a tracer span for the revoke callback endpoint.
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/quickbooks/revoke", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

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
		tracing.TagTenant(span, tenant.(string))

		// Call the QuickBooks service to revoke the access.
		err := s.CommonServices.QuickbooksService.RevokeAccess(ctx)
		if err != nil {
			span.LogFields(log.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke QuickBooks access"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func GetQuickbooksSettings(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/tenant/settings", c.Request.Header)
		defer span.Finish()

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
