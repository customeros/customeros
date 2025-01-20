package mailstack

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
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
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mailstack.DNS", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		var records []DNSRecord

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(log.String("result", "Missing tenant in context"))
			return
		}

		// validate domain belongs to tenant
		domain := c.Param("domain")
		mailboxTenant, err := h.services.CommonServices.MailstackService.GetTenantForMailstackDomain(ctx, domain)
		if err != nil {
			message := "Unable to validate domain ownership"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)

			return
		}
		if !strings.EqualFold(tenant, mailboxTenant) {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		// get dns records
		dns, err := h.services.CommonServices.CloudflareService.GetDNSRecords(ctx, domain)
		if err != nil {
			message := "Unable to get DNS records"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		if dns == nil {
			message := "Unable to loacate DNS records for domain"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}

		for _, record := range *dns {
			records = append(records, DNSRecord{
				ID:      fmt.Sprintf("dns_%s", record.ID),
				Type:    record.Type,
				Content: record.Content,
			})
		}

		h.responseHandler.HandleSuccess(c, DNSResponse{
			records,
		})
	}
}

func (h *MailstackHandler) AddDNSRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mailstack.AddDNSRecord", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(log.String("result", "Missing tenant in context"))
			return
		}

		// validate domain belongs to tenant
		domain := c.Param("domain")
		mailboxTenant, err := h.services.CommonServices.MailstackService.GetTenantForMailstackDomain(ctx, domain)
		if err != nil {
			message := "Unable to validate domain ownership"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		if !strings.EqualFold(tenant, mailboxTenant) {
			message := "domain not found"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}

		domainExists, zoneId, err := h.services.CommonServices.CloudflareService.CheckDomainExists(ctx, domain)
		if err != nil {
			message := "Unable to verify domain"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		if !domainExists {
			message := "domain not found"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}

		// get dns record payload
		record, err := h.getDNSRequestPayload(c)
		if err != nil {
			message := "Unable to parse request"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		err = h.services.CommonServices.CloudflareService.AddDNSRecord(ctx, zoneId, record.Type, record.Name, record.Content, 1, false, nil)
		if err != nil {
			message := "Unable to create DNS record"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, DNSRecordResponse{
			record,
		})
	}
}

func (h *MailstackHandler) DeleteDNSRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mailstack.DeleteDNSRecord", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(log.String("result", "Missing tenant in context"))
			return
		}

		// validate domain belongs to tenant
		domain := c.Param("domain")
		mailboxTenant, err := h.services.CommonServices.MailstackService.GetTenantForMailstackDomain(ctx, domain)
		if err != nil {
			message := "Unable to validate domain ownership"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		if !strings.EqualFold(tenant, mailboxTenant) {
			message := "domain not found"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}

		domainExists, zoneId, err := h.services.CommonServices.CloudflareService.CheckDomainExists(ctx, domain)
		if err != nil {
			message := "Unable to verify domain"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		if !domainExists {
			message := "domain not found"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}

		// delete dns record
		dnsRecordId := c.Param("dnsId")
		dnsRecordId = strings.TrimPrefix(dnsRecordId, "dns_")
		err = h.services.CommonServices.CloudflareService.DeleteDNSRecord(ctx, zoneId, dnsRecordId)
		if err != nil {
			message := "Unable to delete DNS record"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, nil)
	}
}

func (h *MailstackHandler) getDNSRequestPayload(c *gin.Context) (DNSRecord, error) {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Flows.getDNSRequestPayload")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var req DNSRecord
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return req, err
	}

	return req, nil
}
