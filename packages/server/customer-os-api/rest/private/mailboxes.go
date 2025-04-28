package private

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/gin-gonic/gin"

	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

func GetMailboxes(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "GetMailboxes")
		defer spans.Finish()

		// get tenant from context
		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			c.Status(http.StatusNotFound)
			spans.LogKV("result", "Missing tenant in context")
			return
		}

		// Get mailboxes using service
		statusCode, _, mailboxes, err := s.CommonServices.MailstackService.GetMailboxes(ctx, tenant, "", "")
		if err != nil {
			spans.TraceError(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		// Handle non-200 responses
		if statusCode != http.StatusOK {
			c.Status(statusCode)
			return
		}

		c.JSON(http.StatusOK, mailboxes)
	}
}
