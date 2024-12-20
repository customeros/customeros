package restmailstack

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type DNSResponse struct {
	enum.BaseResponse
	Records []DNSRecord `json:"dnsRecords"`
}

func DNS(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mailstack.DNS", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		var records []DNSRecord

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			span.LogFields(log.String("result", "Missing tenant in context"))
			return
		}

		// validate domain belongs to tenant

		domain := c.Param("domain")
		dns, err := services.CommonServices.CloudflareService.GetDNSRecords(ctx, domain)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to get DNS records"))
			return
		}
		if dns == nil {
			rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("unable to locate DNS records for domain"))
			return
		}

		for _, record := range *dns {
			records = append(records, DNSRecord{
				ID:      record.ID,
				Type:    record.Type,
				Content: record.Content,
			})
		}

		c.JSON(http.StatusOK, DNSResponse{
			enum.BuildBaseResponse(enum.StatusSuccess),
			records,
		})
	}
}
