package rest

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

const retryCountFetchFreshData = 15

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
		}

		paymentLink := invoice.PaymentDetails.PaymentLink
		validUntil := invoice.PaymentDetails.PaymentLinkValidUntil
		generateNewLink := false
		if paymentLink == "" {
			generateNewLink = true
		} else if validUntil != nil && validUntil.Before(utils.Now()) {
			generateNewLink = true
		}

		if generateNewLink {
			paymentLink, err = generateAndGetNewStripeCheckoutSession(ctx, services, invoice, tenant)
		}

		if paymentLink == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Please try again later"})
			return
		}

		// If all good, redirect to payment link
		c.Redirect(http.StatusFound, paymentLink)
	}
}

func generateAndGetNewStripeCheckoutSession(ctx context.Context, services *service.Services, invoice *neo4jentity.InvoiceEntity, tenant string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "generateAndGetNewStripeCheckoutSession")
	defer span.Finish()

	// Call integration app to create new payment link
	err := callIntegrationAppWithApiRequestForNewPaymentLink(ctx, services.Cfg.ExternalServices.IntegrationApp.WorkspaceKey, services.Cfg.ExternalServices.IntegrationApp.WorkspaceSecret, tenant, services.Cfg.ExternalServices.IntegrationApp.ApiTriggerUrlCreatePaymentLinks, invoice)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error calling integration app"))
		return "", err
	}

	// Wait for payment link to be generated
	for i := 0; i < retryCountFetchFreshData; i++ {
		// Fetch invoice again to get updated payment link
		latestInvoice, _, err := services.CommonServices.InvoiceService.GetByIdAcrossAllTenants(ctx, invoice.Id)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error fetching invoice"))
			return "", err
		}
		linkValid := latestInvoice.PaymentDetails.PaymentLinkValidUntil == nil ||
			(invoice.PaymentDetails.PaymentLinkValidUntil != nil && latestInvoice.PaymentDetails.PaymentLinkValidUntil.After(utils.Now()))
		if latestInvoice.PaymentDetails.PaymentLink != "" &&
			latestInvoice.PaymentDetails.PaymentLink != invoice.PaymentDetails.PaymentLink &&
			linkValid {
			span.LogKV("result.paymentLink", latestInvoice.PaymentDetails.PaymentLink)
			return latestInvoice.PaymentDetails.PaymentLink, nil
		}
		// sleep for 1 second
		time.Sleep(time.Second)
	}

	return "", nil
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
}

func callIntegrationAppWithApiRequestForNewPaymentLink(ctx context.Context, key, secret, tenant, url string, invoice *neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "callIntegrationAppWithApiRequestForNewPaymentLink")
	defer span.Finish()
	span.LogKV("url", url)
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
