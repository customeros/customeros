package webhook

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type InvoiceLineItem struct {
	Description string
	MetadataID  string
}

// source: https://docs.customeros.ai/en/api/invoice-events#example-event
type InvoicePayload struct {
	Data struct {
		AmountDue          float64 `json:"amountDue"`
		AmountPaid         float64 `json:"amountPaid"`
		AmountRemaining    float64 `json:"amountRemaining"`
		Currency           string  `json:"currency"`
		Due                string  `json:"due"`
		InvoiceNumber      string  `json:"invoiceNumber"`
		InvoicePeriodEnd   string  `json:"invoicePeriodEnd"`
		InvoicePeriodStart string  `json:"invoicePeriodStart"`
		InvoiceUrl         string  `json:"invoiceUrl"`
		Note               string  `json:"note"`
		Paid               bool    `json:"paid"`
		Status             string  `json:"status"`
		Subtotal           float64 `json:"subtotal"`
		TaxDue             float64 `json:"taxDue"`
		Contract           struct {
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
			Created string `json:"created"`
			ID      string `json:"id"`
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

func PopulateInvoiceFinalizedPayload(invoice *neo4jentity.InvoiceEntity, org *neo4jentity.OrganizationEntity, contract *neo4jentity.ContractEntity, slis []*neo4jentity.ServiceLineItemEntity) *InvoicePayload {
	if invoice == nil || org == nil || contract == nil {
		return nil
	}

	invoiceLineItems := make([]InvoiceLineItem, 0)

	for _, sli := range slis {
		invoiceLineItems = append(invoiceLineItems, InvoiceLineItem{
			Description: sli.Name,
			MetadataID:  sli.ID,
		})
	}

	payload := createInvoicePayload(
		invoice.Amount,
		invoice.Amount,
		0,
		invoice.Currency.String(),
		invoice.DueDate.String(),
		invoice.Number,
		invoice.PeriodEndDate.String(),
		invoice.PeriodStartDate.String(),
		"", // FIXME: Where does invoiceUrl come from?
		invoice.Note,
		true, // FIXME: Where does paid come from?
		invoice.Status.String(),
		invoice.Amount,
		invoice.Vat,
		contract.Name,
		contract.ContractStatus.String(),
		contract.Id,
		"", // FIXME: Where does metadataCreated come from?
		"", // FIXME: Where does metadataID come from?
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
	currency,
	due,
	invoiceNumber,
	invoicePeriodEnd,
	invoicePeriodStart,
	invoiceUrl,
	note string,
	paid bool,
	status string,
	subtotal,
	taxDue float64,
	contractName,
	contractStatus,
	contractMetadataID,
	metadataCreated,
	metadataID,
	organizationCustomerOsID,
	organizationMetadataID,
	organizationName string,
	eventType WebhookEvent,
	invoiceLineItems []InvoiceLineItem,
) *InvoicePayload {
	payload := &InvoicePayload{}

	payload.Data.AmountDue = amountDue
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
