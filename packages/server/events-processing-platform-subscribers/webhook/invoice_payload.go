package webhook

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	"time"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
)

type InvoiceLineItem struct {
	Description string
	MetadataID  string
}

// source: https://docs.customeros.ai/en/api/invoice-events#example-event
type InvoicePayload struct {
	Data struct {
		AmountDue               float64   `json:"amountDue"`
		AmountDueInSmallestUnit int64     `json:"amountDueInSmallestUnit"`
		AmountPaid              float64   `json:"amountPaid"`
		AmountRemaining         float64   `json:"amountRemaining"`
		Currency                string    `json:"currency"`
		Due                     time.Time `json:"due"`
		InvoiceNumber           string    `json:"invoiceNumber"`
		InvoicePeriodEnd        time.Time `json:"invoicePeriodEnd"`
		InvoicePeriodStart      time.Time `json:"invoicePeriodStart"`
		InvoiceUrl              string    `json:"invoiceUrl"`
		Note                    string    `json:"note"`
		Paid                    bool      `json:"paid"`
		Status                  string    `json:"status"`
		Subtotal                float64   `json:"subtotal"`
		TaxDue                  float64   `json:"taxDue"`
		Contract                struct {
			ContractName   string `json:"contractName"`
			ContractStatus string `json:"contractStatus"`
			Metadata       struct {
				ID string `json:"id"`
			} `json:"metadata"`
		} `json:"contract"`
		InvoiceLineItems []struct {
			Description string `json:"description"`
			Metadata    struct {
				ID string `json:"id"`
			} `json:"metadata"`
		} `json:"invoiceLineItems"`
		Metadata struct {
			Created time.Time `json:"created"`
			ID      string    `json:"id"`
		} `json:"metadata"`
		Organization struct {
			CustomerOsID string `json:"customerOsId"`
			Metadata     struct {
				ID string `json:"id"`
			} `json:"metadata"`
			Name string `json:"name"`
		} `json:"organization"`
	} `json:"data"`
	Event string `json:"event"`
}

func PopulateInvoicePayload(invoice *neo4jentity.InvoiceEntity, org *neo4jentity.OrganizationEntity, contract *neo4jentity.ContractEntity, ils []*neo4jentity.InvoiceLineEntity) *InvoicePayload {
	if invoice == nil || org == nil || contract == nil || ils == nil {
		return nil
	}

	invoiceLineItems := make([]InvoiceLineItem, 0)

	for _, il := range ils {
		invoiceLineItems = append(invoiceLineItems, InvoiceLineItem{
			Description: il.Name,
			MetadataID:  il.Id,
		})
	}

	paid := false
	amountRemaining := invoice.TotalAmount
	amountPaid := float64(0)

	if invoice.Status == neo4jenum.InvoiceStatusPaid {
		paid = true
		amountRemaining = float64(0)
		amountPaid = invoice.TotalAmount
	}

	payload := createInvoicePayload(
		invoice.TotalAmount,
		amountPaid,
		amountRemaining,
		invoice.Currency.String(),
		invoice.DueDate,
		invoice.Number,
		invoice.PeriodEndDate,
		invoice.PeriodStartDate,
		fmt.Sprintf(constants.UrlFileStoreFileDownloadUrlTemplate, invoice.RepositoryFileId),
		invoice.Note,
		paid,
		invoice.Status.String(),
		invoice.Amount,
		invoice.Vat,
		contract.Name,
		contract.ContractStatus.String(),
		contract.Id,
		invoice.CreatedAt,
		invoice.Id,
		org.CustomerOsId,
		org.ID,
		org.Name,
		WebhookEventInvoiceFinalized,
		invoiceLineItems,
	)

	return payload
}

func createInvoicePayload(
	amountDue,
	amountPaid,
	amountRemaining float64,
	currency string,
	due time.Time,
	invoiceNumber string,
	invoicePeriodEnd,
	invoicePeriodStart time.Time,
	invoiceUrl,
	note string,
	paid bool,
	status string,
	subtotal,
	taxDue float64,
	contractName,
	contractStatus,
	contractMetadataID string,
	metadataCreated time.Time,
	metadataID,
	organizationCustomerOsID,
	organizationMetadataID,
	organizationName string,
	eventType WebhookEvent,
	invoiceLineItems []InvoiceLineItem,
) *InvoicePayload {
	payload := &InvoicePayload{}

	// convert amount to the smallest currency unit
	amountDueInSmallestCurrencyUnit, _ := data.InSmallestCurrencyUnit(currency, amountDue)

	payload.Data.AmountDue = amountDue
	payload.Data.AmountDueInSmallestUnit = amountDueInSmallestCurrencyUnit
	payload.Data.AmountPaid = amountPaid
	payload.Data.AmountRemaining = amountRemaining
	payload.Data.Currency = currency
	payload.Data.Due = due
	payload.Data.InvoiceNumber = invoiceNumber
	payload.Data.InvoicePeriodEnd = invoicePeriodEnd
	payload.Data.InvoicePeriodStart = invoicePeriodStart
	payload.Data.InvoiceUrl = invoiceUrl
	payload.Data.Note = note
	payload.Data.Paid = paid
	payload.Data.Status = status
	payload.Data.Subtotal = subtotal
	payload.Data.TaxDue = taxDue

	payload.Data.Contract.ContractName = contractName
	payload.Data.Contract.ContractStatus = contractStatus
	payload.Data.Contract.Metadata.ID = contractMetadataID

	for _, item := range invoiceLineItems {
		lineItem := struct {
			Description string `json:"description"`
			Metadata    struct {
				ID string `json:"id"`
			} `json:"metadata"`
		}{
			Description: item.Description,
			Metadata: struct {
				ID string `json:"id"`
			}{
				ID: item.MetadataID,
			},
		}
		payload.Data.InvoiceLineItems = append(payload.Data.InvoiceLineItems, lineItem)
	}

	payload.Data.Metadata.Created = metadataCreated
	payload.Data.Metadata.ID = metadataID

	payload.Data.Organization.CustomerOsID = organizationCustomerOsID
	payload.Data.Organization.Metadata.ID = organizationMetadataID
	payload.Data.Organization.Name = organizationName

	payload.Event = eventType.String()

	return payload
}
