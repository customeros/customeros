package public

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

func RedirectToPayInvoice(services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "RedirectToPayInvoice")
		defer spans.Finish()

		// validate integration app is configured
		if services.Cfg.Common.External.IntegrationAppConfig.WorkspaceKey == "" || services.Cfg.Common.External.IntegrationAppConfig.WorkspaceSecret == "" {
			err := errors.New("Integration app not configured")
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to obtain payment link, please try again later"})
			return
		}

		clientIP := getClientIP(c)

		// Get invoice ID from path parameter
		invoiceID := c.Param("invoiceId")
		spans.LogKV("invoiceId", invoiceID)

		// Fetch invoice by ID
		invoice, tenant, err := services.CommonServices.InvoiceService.GetByIdAcrossAllTenants(ctx, invoiceID)
		if err != nil {
			spans.TraceError(err)
		}
		if invoice == nil || invoice.DryRun {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
			return
		}
		spans.TagTenant(tenant)
		spans.LogKV("invoiceStatus", invoice.Status.String())

		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceCustomerOsApiRest,
		})

		// Save Client IP
		saveErr := saveClientIP(innerCtx, services, clientIP, invoiceID, tenant)
		if saveErr != nil {
			spans.TraceError(errors.Wrap(err, "Error saving clientIP"))
		}

		// Check invoice status
		switch invoice.Status {
		case neo4jenum.InvoiceStatusPaid:
			// Handle scenario: Invoice already paid
			c.Redirect(http.StatusSeeOther, services.Cfg.App.InvoicePaidRedirectUrl)
			return
		case neo4jenum.InvoiceStatusVoid:
			// Handle scenario: Invoice voided
			c.JSON(http.StatusGone, gin.H{"error": "Invoice is voided"})
			return
		case neo4jenum.InvoiceStatusPaymentProcessing:
			// Handle scenario: Invoice in payment processing
			c.JSON(http.StatusSeeOther, gin.H{"error": "Processing bank payment. Cannot pay now."})
			return
		case neo4jenum.InvoiceStatusOnHold:
			// Handle scenario: Invoice voided
			c.JSON(http.StatusGone, gin.H{"error": "Invoice is on hold"})
			return
		}

		paymentLink := invoice.PaymentDetails.PaymentLink
		validUntil := invoice.PaymentDetails.PaymentLinkValidUntil
		spans.LogFields(log.String("initial.paymentLink", paymentLink), log.Object("initial.validUntil", validUntil), log.Object("now", utils.Now()))
		generateNewLink := false
		if paymentLink == "" {
			generateNewLink = true
		} else if validUntil != nil && validUntil.Before(utils.Now()) {
			generateNewLink = true
		}
		spans.LogKV("generateNewLink", generateNewLink)

		if generateNewLink {
			err = services.CommonServices.InvoiceService.GenerateNewPaymentLink(innerCtx, invoice.Id)
			if err != nil {
				spans.TraceError(errors.Wrap(err, "error generating payment link"))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to obtain payment link, please try again later"})
				return
			}
			c.File("../static/pay-invoice.html")
			return
		} else {
			// If all good, redirect to payment link
			c.Redirect(http.StatusFound, paymentLink)
			return
		}
	}
}

func GetInvoicePaymentLink(services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RedirectToPayInvoice", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		// Get invoice ID from path parameter
		invoiceID := c.Param("invoiceId")
		span.LogKV("invoiceId", invoiceID)

		// Fetch invoice by ID
		invoice, tenant, err := services.CommonServices.InvoiceService.GetByIdAcrossAllTenants(ctx, invoiceID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error fetching invoice"))
		}
		if invoice == nil || invoice.DryRun {
			c.String(http.StatusNotFound, "")
			return
		}
		tracing.TagTenant(span, tenant)
		span.LogKV("invoiceStatus", invoice.Status.String())

		// Check invoice status
		switch invoice.Status {
		case neo4jenum.InvoiceStatusPaid:
			// Handle scenario: Invoice already paid
			c.String(http.StatusConflict, "")
			return
		case neo4jenum.InvoiceStatusVoid:
			// Handle scenario: Invoice voided
			c.String(http.StatusConflict, "")
			return
		case neo4jenum.InvoiceStatusOnHold:
			// Handle scenario: Invoice on hold
			c.String(http.StatusConflict, "")
		}

		paymentLink := invoice.PaymentDetails.PaymentLink
		validUntil := invoice.PaymentDetails.PaymentLinkValidUntil
		span.LogFields(log.String("paymentLink", paymentLink), log.Object("validUntil", validUntil))

		if paymentLink != "" && validUntil != nil && validUntil.After(utils.Now()) {
			c.String(http.StatusOK, paymentLink)
			return
		} else {
			c.String(http.StatusOK, "")
			return
		}
	}
}

func getClientIP(c *gin.Context) string {
	originalIP := c.Request.Header["X-Original-Forwarded-For"]
	cloudflareIP := c.Request.Header["Cf-Connecting-Ip"]

	if cloudflareIP[0] != "" {
		return cloudflareIP[0]
	}
	if originalIP[0] != "" {
		return originalIP[0]
	}
	return ""
}

func saveClientIP(ctx context.Context, s *cosapi_services.Services, clientIP, invoiceID, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Rest.SaveIP")
	defer span.Finish()
	tracing.TagComponentRest(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("clientIP", clientIP, "invoiceID", invoiceID)

	if clientIP == "" {
		return nil
	}

	contracts, err := s.ContractService.GetContractsForInvoices(ctx, []string{invoiceID})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if len(*contracts) == 0 {
		err := errors.New("Could not find contract for invoice")
		tracing.TraceErr(span, err)
		return err
	}

	invoiceEmail := (*contracts)[0].InvoiceEmail
	verifyEmail := mailvalidate.ValidateEmailSyntax(invoiceEmail)
	if !verifyEmail.IsValid {
		err := errors.New("Invalid invoice email address")
		span.LogKV("invoiceEmail", invoiceEmail)
		tracing.TraceErr(span, err)
		return err
	}

	details := postgres_entity.EnrichDetailsTracking{
		IP:             clientIP,
		CompanyDomain:  &verifyEmail.Domain,
		CompanyWebsite: &verifyEmail.Domain,
		SourceEmail:    &verifyEmail.CleanEmail,
	}

	createErr := s.Repositories.PostgresRepositories.EnrichDetailsTrackingRepository.Save(ctx, details)
	if createErr != nil {
		tracing.TraceErr(span, createErr)
		return createErr
	}
	return nil
}

func notifyOnSlackPaymentFailed(ctx context.Context, services *cosapi_services.Services, tenant, invoiceId, invoiceNumber string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "notifyOnSlackPaymentFailed")
	defer span.Finish()

	tenantSettings, err := services.CommonServices.TenantSettingsService.GetTenantSettingsForTenant(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error fetching tenant settings"))
	}
	if tenantSettings == nil || tenantSettings.SharedSlackChannelUrl == "" {
		span.LogKV("msg", "Notification slack channel URL not set")
		return
	}

	organizationDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error fetching organization"))
		return
	}
	if organizationDbNode == nil {
		span.LogKV("msg", "Organization not found")
		return
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	slackMessageText := fmt.Sprintf("OrganizationRelationshipCustomer %s encountered an error trying to pay invoice %s", organizationEntity.Name, invoiceNumber)

	err = utils.SendSlackMessage(ctx, tenantSettings.SharedSlackChannelUrl, slackMessageText)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error sending slack message"))
	}
}
