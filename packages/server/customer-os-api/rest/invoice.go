package rest

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

func RedirectToPayInvoice(services *service.Services) gin.HandlerFunc {
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
			return
		}
		tracing.TagTenant(span, tenant)
		span.LogKV(log.String("invoiceStatus", invoice.Status.String()))

		// get organization linked to invoice
		organizationDbNode, err := services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		organizationEntity := neo4jentity.OrganizationEntity{}
		if organizationDbNode != nil {
			organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
		}

		// Check invoice status
		switch invoice.Status {
		case neo4jenum.InvoiceStatusPaid:
			// Handle scenario: Invoice already paid
			c.Redirect(http.StatusSeeOther, services.Cfg.AppConfig.InvoicePaidRedirectUrl)
			return
		case neo4jenum.InvoiceStatusVoid:
			// Handle scenario: Invoice voided
			c.JSON(http.StatusGone, gin.H{"error": "Invoice is voided"})
			return
		case neo4jenum.InvoiceStatusOnHold:
			// Handle scenario: Invoice voided
			c.JSON(http.StatusGone, gin.H{"error": "Invoice is on holde"})
			return
		}

		paymentLink := invoice.PaymentDetails.PaymentLink
		validUntil := invoice.PaymentDetails.PaymentLinkValidUntil
		span.LogFields(log.String("initial.paymentLink", paymentLink), log.Object("initial.validUntil", validUntil), log.Object("now", utils.Now()))
		generateNewLink := false
		if paymentLink == "" {
			generateNewLink = true
		} else if validUntil != nil && validUntil.Before(utils.Now()) {
			generateNewLink = true
		}
		span.LogFields(log.Bool("generateNewLink", generateNewLink))

		if generateNewLink {
			paymentLink = ""

			primaryStripeCustomerId, err := services.CommonServices.ExternalSystemService.GetPrimaryExternalId(ctx, neo4jenum.Stripe.String(), organizationEntity.ID, model.ORGANIZATION)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error fetching primary stripe customer ID"))
			}

			err = callIntegrationAppWithApiRequestForNewPaymentLink(ctx, services.Cfg.ExternalServices.IntegrationApp.WorkspaceKey,
				services.Cfg.ExternalServices.IntegrationApp.WorkspaceSecret, tenant,
				services.Cfg.ExternalServices.IntegrationApp.ApiTriggerUrlCreatePaymentLinks, primaryStripeCustomerId, invoice)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error calling integration app"))
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

func GetInvoicePaymentLink(services *service.Services) gin.HandlerFunc {
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
		span.LogKV(log.String("invoiceStatus", invoice.Status.String()))

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

type ApiRequestCreatePaymentLinks struct {
	Input ApiRequestCreatePaymentLinksInput `json:"input"`
}

type ApiRequestCreatePaymentLinksInput struct {
	InvoiceId                    string `json:"invoiceId"`
	AmountInSmallestCurrencyUnit int64  `json:"amountInSmallestCurrencyUnit"`
	Currency                     string `json:"currency"`
	InvoiceDescription           string `json:"invoiceDescription"`
	CustomerEmail                string `json:"customerEmail"`
	PrimaryStripeCustomerId      string `json:"stripeCustomerId"`
}

func callIntegrationAppWithApiRequestForNewPaymentLink(ctx context.Context, key, secret, tenant, url, primaryStripeCustomerId string, invoice *neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "callIntegrationAppWithApiRequestForNewPaymentLink")
	defer span.Finish()
	span.LogKV("url", url)
	span.LogKV("primaryStripeCustomerId", primaryStripeCustomerId)
	span.SetTag(tracing.SpanTagTenant, tenant)

	var SigningKey = []byte(secret)

	claims := jwt.MapClaims{
		"id":   tenant,
		"name": tenant,
		// To prevent token from being used for too long
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iss": key,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SigningKey)
	if err != nil {
		return errors.Wrap(err, "Error signing JWT token")
	}

	amountInSmallestCurrencyUnit, err := data.InSmallestCurrencyUnit(invoice.Currency.String(), invoice.TotalAmount)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error converting amount to smallest currency unit"))
		return err
	}

	input := ApiRequestCreatePaymentLinks{
		Input: ApiRequestCreatePaymentLinksInput{
			InvoiceId:                    invoice.Id,
			AmountInSmallestCurrencyUnit: amountInSmallestCurrencyUnit,
			Currency:                     invoice.Currency.String(),
			InvoiceDescription:           fmt.Sprintf("Invoice %s", invoice.Number),
			CustomerEmail:                invoice.Customer.Email,
			PrimaryStripeCustomerId:      primaryStripeCustomerId,
		},
	}
	payload, err := json.Marshal(input)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error marshalling input"))
		return err
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating HTTP request"))
		return err
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)
	req.Header.Add("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error calling integration app"))
		return err
	}
	return nil
}

func notifyOnSlackPaymentFailed(ctx context.Context, services *service.Services, tenant, invoiceId, invoiceNumber string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "notifyOnSlackPaymentFailed")
	defer span.Finish()

	tenantSettings, err := services.CommonServices.TenantService.GetTenantSettingsForTenant(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error fetching tenant settings"))
	}
	if tenantSettings == nil || tenantSettings.SharedSlackChannelUrl == "" {
		span.LogKV("msg", "Notification slack channel URL not set")
		return
	}

	organizationDbNode, err := services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error fetching organization"))
		return
	}
	if organizationDbNode == nil {
		span.LogKV("msg", "Organization not found")
		return
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	slackMessageText := fmt.Sprintf("Customer %s encountered an error trying to pay invoice %s", organizationEntity.Name, invoiceNumber)

	err = utils.SendSlackMessage(ctx, tenantSettings.SharedSlackChannelUrl, slackMessageText)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error sending slack message"))
	}
}
