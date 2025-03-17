package private

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	tracingLog "github.com/opentracing/opentracing-go/log"

	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

func GetMailboxes(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetMailboxes", c.Request.Header)
		defer span.Finish()
		tracing.SetDefaultRestSpanTags(ctx, span)

		// get tenant from context
		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			c.Status(http.StatusNotFound)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Get mailboxes using service
		statusCode, _, mailboxes, err := s.CommonServices.MailstackService.GetMailboxes(ctx, tenant, "", "")
		if err != nil {
			tracing.TraceErr(span, err)
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
