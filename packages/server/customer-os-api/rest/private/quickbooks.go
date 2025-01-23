package private

import (
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
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

		quickbooksRequestAccessUrl := "https://appcenter.intuit.com/connect/oauth2?client_id=" + s.Cfg.Common.External.QuickbooksConfig.ClientId + "&redirect_uri=" + s.Cfg.Common.External.QuickbooksConfig.RedirectUrl + "&response_type=code&scope=com.intuit.quickbooks.accounting&state=12345"

		span.LogFields(log.Object("quickbooksRequestAccessUrl", quickbooksRequestAccessUrl))

		c.JSON(http.StatusOK, gin.H{"url": quickbooksRequestAccessUrl})
	}
}

func CallbackQuickbooks(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/slack/oauth/callback", c.Request.Header)
		defer span.Finish()

		code := c.Request.URL.Query().Get("code")
		realmId := c.Request.URL.Query().Get("realmId")

		requestData := url.Values{}
		requestData.Set("grant_type", "authorization_code")
		requestData.Set("code", code)
		requestData.Set("redirect_uri", s.Cfg.Common.External.QuickbooksConfig.RedirectUrl)

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

		c.JSON(http.StatusOK, gin.H{})
	}
}
