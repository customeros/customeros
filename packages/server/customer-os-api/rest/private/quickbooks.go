package private

import (
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

func RequestAccessQuickbooks(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/quickbooks/requestAccess", c.Request.Header)
		defer span.Finish()

		quickbooksRequestAccessUrl := "https://appcenter.intuit.com/connect/oauth2?client_id=" + s.Cfg.Common.External.QuickbooksCofig.ClientId + "&redirect_uri=" + s.Cfg.Common.External.QuickbooksCofig.RedirectUrl + "&response_type=code&scope=com.intuit.quickbooks.accounting&state=12345"

		span.LogFields(log.Object("quickbooksRequestAccessUrl", quickbooksRequestAccessUrl))

		c.JSON(http.StatusOK, gin.H{"url": quickbooksRequestAccessUrl})
	}
}
