package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ContractUpdateEvent struct {
	Tenant                 string                     `json:"tenant" validate:"required"`
	Name                   string                     `json:"name"`
	ContractUrl            string                     `json:"contractUrl"`
	UpdatedAt              time.Time                  `json:"updatedAt"`
	ServiceStartedAt       *time.Time                 `json:"serviceStartedAt,omitempty"`
	SignedAt               *time.Time                 `json:"signedAt,omitempty"`
	EndedAt                *time.Time                 `json:"endedAt,omitempty"`
	ExternalSystem         commonmodel.ExternalSystem `json:"externalSystem,omitempty"`
	Source                 string                     `json:"source"`
	InvoicingStartDate     *time.Time                 `json:"invoicingStartDate,omitempty"`
	Currency               string                     `json:"currency,omitempty"`
	BillingCycle           string                     `json:"billingCycle,omitempty"` //Deprecated: Use BillingCycleInMonths instead
	BillingCycleInMonths   int64                      `json:"billingCycleInMonths,omitempty"`
	AddressLine1           string                     `json:"addressLine1,omitempty"`
	AddressLine2           string                     `json:"addressLine2,omitempty"`
	Locality               string                     `json:"locality,omitempty"`
	Country                string                     `json:"country,omitempty"`
	Region                 string                     `json:"region,omitempty"`
	Zip                    string                     `json:"zip,omitempty"`
	OrganizationLegalName  string                     `json:"organizationLegalName,omitempty"`
	InvoiceEmail           string                     `json:"invoiceEmail,omitempty" validate:"omitempty,email"`
	InvoiceEmailCC         []string                   `json:"invoiceEmailCC,omitempty" validate:"omitempty,dive,email"`
	InvoiceEmailBCC        []string                   `json:"invoiceEmailBCC,omitempty" validate:"omitempty,dive,email"`
	InvoiceNote            string                     `json:"invoiceNote,omitempty"`
	FieldsMask             []string                   `json:"fieldsMask,omitempty"`
	NextInvoiceDate        *time.Time                 `json:"nextInvoiceDate,omitempty"`
	CanPayWithCard         bool                       `json:"canPayWithCard,omitempty"`
	CanPayWithDirectDebit  bool                       `json:"canPayWithDirectDebit,omitempty"`
	CanPayWithBankTransfer bool                       `json:"canPayWithBankTransfer,omitempty"`
	PayOnline              bool                       `json:"payOnline,omitempty"`
	PayAutomatically       bool                       `json:"payAutomatically,omitempty"`
	InvoicingEnabled       bool                       `json:"invoicingEnabled,omitempty"`
	AutoRenew              bool                       `json:"autoRenew,omitempty"`
	Check                  bool                       `json:"check,omitempty"`
	DueDays                int64                      `json:"dueDays,omitempty"`
	LengthInMonths         int64                      `json:"lengthInMonths,omitempty"`
	Approved               bool                       `json:"approved,omitempty"`
}

func NewContractUpdateEvent(a eventstore.Aggregate, dataFields model.ContractDataFields, externalSystem commonmodel.ExternalSystem, source string, updatedAt time.Time, fieldsMask []string) (eventstore.Event, error) {
	eventData := ContractUpdateEvent{
		Tenant:                 a.GetTenant(),
		Name:                   dataFields.Name,
		ContractUrl:            dataFields.ContractUrl,
		ServiceStartedAt:       utils.ToDatePtr(dataFields.ServiceStartedAt),
		SignedAt:               utils.ToDatePtr(dataFields.SignedAt),
		EndedAt:                utils.ToDatePtr(dataFields.EndedAt),
		Currency:               dataFields.Currency,
		BillingCycleInMonths:   dataFields.BillingCycleInMonths,
		AddressLine1:           dataFields.AddressLine1,
		AddressLine2:           dataFields.AddressLine2,
		Locality:               dataFields.Locality,
		Country:                dataFields.Country,
		Region:                 dataFields.Region,
		Zip:                    dataFields.Zip,
		OrganizationLegalName:  dataFields.OrganizationLegalName,
		CanPayWithCard:         dataFields.CanPayWithCard,
		CanPayWithDirectDebit:  dataFields.CanPayWithDirectDebit,
		CanPayWithBankTransfer: dataFields.CanPayWithBankTransfer,
		PayOnline:              dataFields.PayOnline,
		PayAutomatically:       dataFields.PayAutomatically,
		InvoicingEnabled:       dataFields.InvoicingEnabled,
		AutoRenew:              dataFields.AutoRenew,
		Check:                  dataFields.Check,
		DueDays:                dataFields.DueDays,
		UpdatedAt:              updatedAt,
		Source:                 source,
		FieldsMask:             fieldsMask,
		LengthInMonths:         dataFields.LengthInMonths,
		Approved:               dataFields.Approved,
	}

	if eventData.LengthInMonths < 0 {
		eventData.LengthInMonths = 0
	} else if eventData.LengthInMonths > 1200 {
		eventData.LengthInMonths = 1200
	}

	if eventData.DueDays < 0 {
		eventData.DueDays = 0
	} else if eventData.DueDays > 365 {
		eventData.DueDays = 365
	}
	if eventData.UpdateNextInvoiceDate() {
		eventData.NextInvoiceDate = utils.ToDatePtr(dataFields.NextInvoiceDate)
	}
	if eventData.UpdateInvoicingStartDate() {
		eventData.InvoicingStartDate = utils.ToDatePtr(dataFields.InvoicingStartDate)
	}
	if eventData.UpdateInvoiceNote() {
		eventData.InvoiceNote = dataFields.InvoiceNote
	}
	if eventData.UpdateInvoiceEmail() {
		eventData.InvoiceEmail = dataFields.InvoiceEmail
	}
	if eventData.UpdateInvoiceEmailCC() {
		eventData.InvoiceEmailCC = dataFields.InvoiceEmailCC
	}
	if eventData.UpdateInvoiceEmailBCC() {
		eventData.InvoiceEmailBCC = dataFields.InvoiceEmailBCC
	}

	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContractUpdateEvent")
	}

	event := eventstore.NewBaseEvent(a, ContractUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContractUpdateEvent")
	}
	return event, nil
}

func (e ContractUpdateEvent) UpdateStatus() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskStatus)
}

func (e ContractUpdateEvent) UpdateContractUrl() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskContractURL)
}

func (e ContractUpdateEvent) UpdateName() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskName)
}

func (e ContractUpdateEvent) UpdateSignedAt() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskSignedAt)
}

func (e ContractUpdateEvent) UpdateEndedAt() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskEndedAt)
}

func (e ContractUpdateEvent) UpdateServiceStartedAt() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskServiceStartedAt)
}

func (e ContractUpdateEvent) UpdateInvoicingStartDate() bool {
	return utils.Contains(e.FieldsMask, FieldMaskInvoicingStartDate)
}

func (e ContractUpdateEvent) UpdateBillingCycle() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBillingCycle)
}

func (e ContractUpdateEvent) UpdateBillingCycleInMonths() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBillingCycleInMonths)
}

func (e ContractUpdateEvent) UpdateCurrency() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskCurrency)
}

func (e ContractUpdateEvent) UpdateAddressLine1() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskAddressLine1)
}

func (e ContractUpdateEvent) UpdateAddressLine2() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskAddressLine2)
}

func (e ContractUpdateEvent) UpdateZip() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskZip)
}

func (e ContractUpdateEvent) UpdateCountry() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskCountry)
}

func (e ContractUpdateEvent) UpdateRegion() bool {
	return utils.Contains(e.FieldsMask, FieldMaskRegion)
}

func (e ContractUpdateEvent) UpdateLocality() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskLocality)
}

func (e ContractUpdateEvent) UpdateOrganizationLegalName() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskOrganizationLegalName)
}

func (e ContractUpdateEvent) UpdateInvoiceEmail() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskInvoiceEmail)
}

func (e ContractUpdateEvent) UpdateInvoiceEmailCC() bool {
	return utils.Contains(e.FieldsMask, FieldMaskInvoiceEmailCC)
}

func (e ContractUpdateEvent) UpdateInvoiceEmailBCC() bool {
	return utils.Contains(e.FieldsMask, FieldMaskInvoiceEmailBCC)
}

func (e ContractUpdateEvent) UpdateInvoiceNote() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskInvoiceNote)
}

func (e ContractUpdateEvent) UpdateCanPayWithCard() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskCanPayWithCard)
}

func (e ContractUpdateEvent) UpdateCanPayWithDirectDebit() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskCanPayWithDirectDebit)
}

func (e ContractUpdateEvent) UpdateCanPayWithBankTransfer() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskCanPayWithBankTransfer)
}

func (e ContractUpdateEvent) UpdateNextInvoiceDate() bool {
	return utils.Contains(e.FieldsMask, FieldMaskNextInvoiceDate)
}

func (e ContractUpdateEvent) UpdateInvoicingEnabled() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskInvoicingEnabled)
}

func (e ContractUpdateEvent) UpdatePayOnline() bool {
	return utils.Contains(e.FieldsMask, FieldMaskPayOnline)
}

func (e ContractUpdateEvent) UpdatePayAutomatically() bool {
	return utils.Contains(e.FieldsMask, FieldMaskPayAutomatically)
}

func (e ContractUpdateEvent) UpdateAutoRenew() bool {
	return utils.Contains(e.FieldsMask, FieldMaskAutoRenew)
}

func (e ContractUpdateEvent) UpdateCheck() bool {
	return utils.Contains(e.FieldsMask, FieldMaskCheck)
}

func (e ContractUpdateEvent) UpdateDueDays() bool {
	return utils.Contains(e.FieldsMask, FieldMaskDueDays)
}

func (e ContractUpdateEvent) UpdateLengthInMonths() bool {
	return utils.Contains(e.FieldsMask, FieldMaskLengthInMonths)
}

func (e ContractUpdateEvent) UpdateApproved() bool {
	return utils.Contains(e.FieldsMask, FieldMaskApproved)
}
