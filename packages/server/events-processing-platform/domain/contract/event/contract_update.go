package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type ContractUpdateEvent struct {
	Tenant                string                     `json:"tenant" validate:"required"`
	Name                  string                     `json:"name"`
	ContractUrl           string                     `json:"contractUrl"`
	UpdatedAt             time.Time                  `json:"updatedAt"`
	ServiceStartedAt      *time.Time                 `json:"serviceStartedAt,omitempty"`
	SignedAt              *time.Time                 `json:"signedAt,omitempty"`
	EndedAt               *time.Time                 `json:"endedAt,omitempty"`
	RenewalCycle          string                     `json:"renewalCycle"`
	RenewalPeriods        *int64                     `json:"renewalPeriods,omitempty"`
	Status                string                     `json:"status"`
	ExternalSystem        commonmodel.ExternalSystem `json:"externalSystem,omitempty"`
	Source                string                     `json:"source"`
	InvoicingStartDate    *time.Time                 `json:"invoicingStartDate,omitempty"`
	Currency              string                     `json:"currency,omitempty"`
	BillingCycle          string                     `json:"billingCycle,omitempty"`
	AddressLine1          string                     `json:"addressLine1,omitempty"`
	AddressLine2          string                     `json:"addressLine2,omitempty"`
	Locality              string                     `json:"locality,omitempty"`
	Country               string                     `json:"country,omitempty"`
	Zip                   string                     `json:"zip,omitempty"`
	OrganizationLegalName string                     `json:"organizationLegalName,omitempty"`
	InvoiceEmail          string                     `json:"invoiceEmail,omitempty"`
	FieldsMask            []string                   `json:"fieldsMask,omitempty"`
}

func NewContractUpdateEvent(aggr eventstore.Aggregate, dataFields model.ContractDataFields, externalSystem commonmodel.ExternalSystem, source string, updatedAt time.Time, fieldsMask []string) (eventstore.Event, error) {
	eventData := ContractUpdateEvent{
		Tenant:                aggr.GetTenant(),
		Name:                  dataFields.Name,
		ContractUrl:           dataFields.ContractUrl,
		ServiceStartedAt:      dataFields.ServiceStartedAt,
		SignedAt:              dataFields.SignedAt,
		EndedAt:               dataFields.EndedAt,
		RenewalCycle:          dataFields.RenewalCycle,
		RenewalPeriods:        dataFields.RenewalPeriods,
		Status:                dataFields.Status.String(),
		Currency:              dataFields.Currency,
		BillingCycle:          dataFields.BillingCycle,
		AddressLine1:          dataFields.AddressLine1,
		AddressLine2:          dataFields.AddressLine2,
		Locality:              dataFields.Locality,
		Country:               dataFields.Country,
		Zip:                   dataFields.Zip,
		OrganizationLegalName: dataFields.OrganizationLegalName,
		InvoiceEmail:          dataFields.InvoiceEmail,
		InvoicingStartDate:    utils.ToDatePtr(dataFields.InvoicingStartDate),
		UpdatedAt:             updatedAt,
		Source:                source,
		FieldsMask:            fieldsMask,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContractUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggr, ContractUpdateV1)
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
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskInvoicingStartDate)
}

func (e ContractUpdateEvent) UpdateRenewalCycle() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskRenewalCycle)
}

func (e ContractUpdateEvent) UpdateRenewalPeriods() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskRenewalPeriods)
}

func (e ContractUpdateEvent) UpdateBillingCycle() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskBillingCycle)
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

func (e ContractUpdateEvent) UpdateLocality() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskLocality)
}

func (e ContractUpdateEvent) UpdateOrganizationLegalName() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskOrganizationLegalName)
}

func (e ContractUpdateEvent) UpdateInvoiceEmail() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskInvoiceEmail)
}
