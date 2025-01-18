package mailstack

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-api/enum"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type DNSResponse struct {
	enum.BaseResponse
	Records []DNSRecord `json:"dnsRecords"`
}

type DNSRecordResponse struct {
	enum.BaseResponse
	Record DNSRecord `json:"dnsRecord"`
}

func DNS(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mailstack.DNS", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		var records []DNSRecord

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			span.LogFields(log.String("result", "Missing tenant in context"))
			return
		}

		// validate domain belongs to tenant
		domain := c.Param("domain")
		mailboxTenant, err := s.CommonServices.MailstackService.GetTenantForMailstackDomain(ctx, domain)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to validate domain ownership"))
			return
		}
		if !strings.EqualFold(tenant, mailboxTenant) {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound)
			return
		}

		// get dns records
		dns, err := s.CommonServices.CloudflareService.GetDNSRecords(ctx, domain)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to get DNS records"))
			return
		}
		if dns == nil {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("unable to locate DNS records for domain"))
			return
		}

		for _, record := range *dns {
			records = append(records, DNSRecord{
				ID:      fmt.Sprintf("dns_%s", record.ID),
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

func AddDNSRecord(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mailstack.AddDNSRecord", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			span.LogFields(log.String("result", "Missing tenant in context"))
			return
		}

		// validate domain belongs to tenant
		domain := c.Param("domain")
		mailboxTenant, err := s.CommonServices.MailstackService.GetTenantForMailstackDomain(ctx, domain)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to validate domain ownership"))
			return
		}
		if !strings.EqualFold(tenant, mailboxTenant) {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("domain not found"))
			return
		}

		domainExists, zoneId, err := s.CommonServices.CloudflareService.CheckDomainExists(ctx, domain)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to verify domain"))
			return
		}
		if !domainExists {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("domain not found"))
			return
		}

		// get dns record payload
		record, err := getDNSRequestPayload(c)
		if err != nil {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("unable to parse request"))
			return
		}
		err = s.CommonServices.CloudflareService.AddDNSRecord(ctx, zoneId, record.Type, record.Name, record.Content, 1, false, nil)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to create DNS record"))
			return
		}

		c.JSON(http.StatusCreated, DNSRecordResponse{
			enum.BuildBaseResponse(enum.StatusSuccess),
			record,
		})
	}
}

func DeleteDNSRecord(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mailstack.DeleteDNSRecord", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			span.LogFields(log.String("result", "Missing tenant in context"))
			return
		}

		// validate domain belongs to tenant
		domain := c.Param("domain")
		mailboxTenant, err := s.CommonServices.MailstackService.GetTenantForMailstackDomain(ctx, domain)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to validate domain ownership"))
			return
		}
		if !strings.EqualFold(tenant, mailboxTenant) {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("domain not found"))
			return
		}

		domainExists, zoneId, err := s.CommonServices.CloudflareService.CheckDomainExists(ctx, domain)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to verify domain"))
			return
		}
		if !domainExists {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("domain not found"))
			return
		}

		// delete dns record
		dnsRecordId := c.Param("dnsId")
		dnsRecordId = strings.TrimPrefix(dnsRecordId, "dns_")
		err = s.CommonServices.CloudflareService.DeleteDNSRecord(ctx, zoneId, dnsRecordId)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to delete dns record"))
			return
		}

		c.JSON(http.StatusNoContent,
			enum.BuildBaseResponse(enum.StatusSuccess),
		)
	}
}

func getDNSRequestPayload(c *gin.Context) (DNSRecord, error) {
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
