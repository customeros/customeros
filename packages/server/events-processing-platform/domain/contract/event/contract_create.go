package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ContractCreateEvent struct {
	Tenant                 string                     `json:"tenant" validate:"required"`
	OrganizationId         string                     `json:"organizationId" validate:"required"`
	Name                   string                     `json:"name"`
	ContractUrl            string                     `json:"contractUrl"`
	CreatedByUserId        string                     `json:"createdByUserId"`
	ServiceStartedAt       *time.Time                 `json:"serviceStartedAt,omitempty"`
	SignedAt               *time.Time                 `json:"signedAt,omitempty"`
	RenewalCycle           string                     `json:"renewalCycle"`
	RenewalPeriods         *int64                     `json:"renewalPeriods,omitempty"`
	Status                 string                     `json:"status"`
	CreatedAt              time.Time                  `json:"createdAt"`
	UpdatedAt              time.Time                  `json:"updatedAt"`
	Source                 commonmodel.Source         `json:"source"`
	ExternalSystem         commonmodel.ExternalSystem `json:"externalSystem,omitempty"`
	InvoicingStartDate     *time.Time                 `json:"invoicingStartDate,omitempty"`
	Currency               string                     `json:"currency"`
	BillingCycle           string                     `json:"billingCycle"`
	InvoicingEnabled       bool                       `json:"invoicingEnabled"`
	PayOnline              bool                       `json:"payOnline,omitempty"`
	PayAutomatically       bool                       `json:"payAutomatically,omitempty"`
	CanPayWithCard         bool                       `json:"canPayWithCard,omitempty"`
	CanPayWithDirectDebit  bool                       `json:"canPayWithDirectDebit,omitempty"`
	CanPayWithBankTransfer bool                       `json:"canPayWithBankTransfer,omitempty"`
	AutoRenew              bool                       `json:"autoRenew,omitempty"`
}

func NewContractCreateEvent(aggregate eventstore.Aggregate, dataFields model.ContractDataFields, source commonmodel.Source, externalSystem commonmodel.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContractCreateEvent{
		Tenant:                 aggregate.GetTenant(),
		OrganizationId:         dataFields.OrganizationId,
		Name:                   dataFields.Name,
		ContractUrl:            dataFields.ContractUrl,
		CreatedByUserId:        dataFields.CreatedByUserId,
		ServiceStartedAt:       dataFields.ServiceStartedAt,
		SignedAt:               dataFields.SignedAt,
		RenewalCycle:           dataFields.RenewalCycle,
		RenewalPeriods:         dataFields.RenewalPeriods,
		Status:                 dataFields.Status.String(),
		Currency:               dataFields.Currency,
		BillingCycle:           dataFields.BillingCycle,
		InvoicingStartDate:     utils.ToDatePtr(dataFields.InvoicingStartDate),
		InvoicingEnabled:       dataFields.InvoicingEnabled,
		PayOnline:              dataFields.PayOnline,
		PayAutomatically:       dataFields.PayAutomatically,
		CanPayWithCard:         dataFields.CanPayWithCard,
		CanPayWithDirectDebit:  dataFields.CanPayWithDirectDebit,
		CanPayWithBankTransfer: dataFields.CanPayWithBankTransfer,
		AutoRenew:              dataFields.AutoRenew,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
		Source:                 source,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContractCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContractCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContractCreateEvent")
	}
	return event, nil
}
