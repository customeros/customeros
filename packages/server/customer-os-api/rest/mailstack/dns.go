package mailstack

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/gin-gonic/gin"
)

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type DNSResponse struct {
	Records []DNSRecord `json:"dnsRecords"`
}

type DNSRecordResponse struct {
	Record DNSRecord `json:"dnsRecord"`
}

func (h *MailstackHandler) DNS() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "MailstackHandler.DNS")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			spans.LogKV("result", "Missing tenant in context")
			return
		}

		domain := c.Param("domain")

		// Call service to get DNS records
		statusCode, errorMsg, records, err := h.services.CommonServices.MailstackService.GetDNSRecords(ctx, tenant, domain)
		if err != nil || statusCode != http.StatusOK {
			if errorMsg != "" {
				h.responseHandler.HandleError(c, statusCode, &errorMsg)
			} else {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			}
			return
		}

		// Convert interface records to response type
		var dnsRecords []DNSRecord
		for _, record := range records {
			dnsRecords = append(dnsRecords, DNSRecord{
				ID:      record.ID,
				Type:    record.Type,
				Name:    record.Name,
				Content: record.Content,
			})
		}

		h.responseHandler.HandleSuccess(c, DNSResponse{
			Records: dnsRecords,
		})
	}
}

func (h *MailstackHandler) AddDNSRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "MailstackHandler.AddDNSRecord")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			spans.LogKV("result", "Missing tenant in context")
			return
		}

		domain := c.Param("domain")

		// get dns record payload
		record, err := h.getDNSRequestPayload(c)
		if err != nil {
			message := "Unable to parse request"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Convert to interfaces.DNSRecord
		interfaceRecord := interfaces.DNSRecord{
			ID:      record.ID,
			Type:    record.Type,
			Name:    record.Name,
			Content: record.Content,
		}

		// Call service to add DNS record
		statusCode, errorMsg, dnsRecord, err := h.services.CommonServices.MailstackService.AddDNSRecord(ctx, tenant, domain, interfaceRecord)
		if err != nil || statusCode != http.StatusOK {
			if errorMsg != "" {
				h.responseHandler.HandleError(c, statusCode, &errorMsg)
			} else {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			}
			return
		}

		h.responseHandler.HandleSuccess(c, DNSRecordResponse{
			Record: DNSRecord{
				ID:      dnsRecord.ID,
				Type:    dnsRecord.Type,
				Name:    dnsRecord.Name,
				Content: dnsRecord.Content,
			},
		})
	}
}

func (h *MailstackHandler) DeleteDNSRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "MailstackHandler.DeleteDNSRecord")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			spans.LogKV("result", "Missing tenant in context")
			return
		}

		domain := c.Param("domain")
		dnsId := c.Param("dnsId")

		// Call service to delete DNS record
		statusCode, errorMsg, err := h.services.CommonServices.MailstackService.DeleteDNSRecord(ctx, tenant, domain, dnsId)
		if err != nil || statusCode != http.StatusOK {
			if errorMsg != "" {
				h.responseHandler.HandleError(c, statusCode, &errorMsg)
			} else {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			}
			return
		}

		h.responseHandler.HandleSuccess(c, nil)
	}
}

func (h *MailstackHandler) getDNSRequestPayload(c *gin.Context) (DNSRecord, error) {
	spans, _ := telemetry.StartRestSpan(c.Request.Context(), "MailstackHandler.getDNSRequestPayload")
	defer spans.Finish()

	var req DNSRecord
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return req, err
	}

	return req, nil
}
