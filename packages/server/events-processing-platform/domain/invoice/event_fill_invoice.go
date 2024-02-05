package invoice

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type InvoiceFillEvent struct {
	Tenant                        string                   `json:"tenant" validate:"required"`
	UpdatedAt                     time.Time                `json:"updatedAt"`
	Amount                        float64                  `json:"amount"`
	VAT                           float64                  `json:"vat"`
	TotalAmount                   float64                  `json:"totalAmount" `
	InvoiceLines                  []InvoiceLineEvent       `json:"invoiceLines" validate:"required"`
	ContractId                    string                   `json:"contractId"`
	DryRun                        bool                     `json:"dryRun"`
	InvoiceNumber                 string                   `json:"invoiceNumber"`
	Currency                      string                   `json:"currency"`
	PeriodStartDate               time.Time                `json:"periodStartDate"`
	PeriodEndDate                 time.Time                `json:"periodEndDate"`
	BillingCycle                  string                   `json:"billingCycle"`
	Status                        string                   `json:"status"`
	Note                          string                   `json:"note"`
	DomesticPaymentsBankInfo      string                   `json:"domesticPaymentsBankInfo"`
	InternationalPaymentsBankInfo string                   `json:"internationalPaymentsBankInfo"`
	Customer                      InvoiceFillCustomerEvent `json:"customer"`
	Provider                      InvoiceFillProviderEvent `json:"provider"`
}
type InvoiceFillCustomerEvent struct {
	Name         string `json:"name"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	Zip          string `json:"zip"`
	Locality     string `json:"locality"`
	Country      string `json:"country"`
	Email        string `json:"email"`
}
type InvoiceFillProviderEvent struct {
	LogoUrl      string `json:"logoUrl"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	Zip          string `json:"zip"`
	Locality     string `json:"locality"`
	Country      string `json:"country"`
}

func NewInvoiceFillEvent(aggregate eventstore.Aggregate, updatedAt time.Time, invoice Invoice,
	domesticPaymentsBankInfo, internationalPaymentsBankInfo,
	customerName, customerAddressLine1, customerAddressLine2, customerAddressZip, customerAddressLocality, customerAddressCountry, customerEmail,
	providerLogoUrl, providerName, providerEmail, providerAddressLine1, providerAddressLine2, providerAddressZip, providerAddressLocality, providerAddressCountry,
	note, status, invoiceNumber string, amount, vat, totalAmount float64, invoiceLines []InvoiceLineEvent) (eventstore.Event, error) {
	eventData := InvoiceFillEvent{
		Tenant:                        aggregate.GetTenant(),
		UpdatedAt:                     updatedAt,
		Amount:                        amount,
		VAT:                           vat,
		TotalAmount:                   totalAmount,
		Currency:                      invoice.Currency,
		ContractId:                    invoice.ContractId,
		InvoiceLines:                  invoiceLines,
		BillingCycle:                  invoice.BillingCycle,
		PeriodStartDate:               invoice.PeriodStartDate,
		PeriodEndDate:                 invoice.PeriodEndDate,
		InvoiceNumber:                 invoiceNumber,
		DryRun:                        invoice.DryRun,
		Status:                        status,
		Note:                          note,
		DomesticPaymentsBankInfo:      domesticPaymentsBankInfo,
		InternationalPaymentsBankInfo: internationalPaymentsBankInfo,
		Customer: InvoiceFillCustomerEvent{
			Name:         customerName,
			Email:        customerEmail,
			AddressLine1: customerAddressLine1,
			AddressLine2: customerAddressLine2,
			Zip:          customerAddressZip,
			Locality:     customerAddressLocality,
			Country:      customerAddressCountry,
		},
		Provider: InvoiceFillProviderEvent{
			LogoUrl:      providerLogoUrl,
			Name:         providerName,
			Email:        providerEmail,
			AddressLine1: providerAddressLine1,
			AddressLine2: providerAddressLine2,
			Zip:          providerAddressZip,
			Locality:     providerAddressLocality,
			Country:      providerAddressCountry,
		},
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceFillEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceFillV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceFillEvent")
	}

	return event, nil
}

type InvoiceLineEvent struct {
	Id                      string             `json:"id" validate:"required"`
	CreatedAt               time.Time          `json:"createdAt" validate:"required"`
	SourceFields            commonmodel.Source `json:"sourceFields"`
	Name                    string             `json:"name" validate:"required"`
	Price                   float64            `json:"price" validate:"required"`
	Quantity                int64              `json:"quantity" validate:"required"`
	Amount                  float64            `json:"amount" validate:"required"`
	VAT                     float64            `json:"vat" validate:"required"`
	TotalAmount             float64            `json:"totalAmount" validate:"required"`
	ServiceLineItemId       string             `json:"serviceLineItemId"`
	ServiceLineItemParentId string             `json:"serviceLineItemParentId"`
	BilledType              string             `json:"billedType" validate:"required"`
}
