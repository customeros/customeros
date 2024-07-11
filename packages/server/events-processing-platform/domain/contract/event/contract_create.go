package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
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
	Status                 string                     `json:"status"`
	CreatedAt              time.Time                  `json:"createdAt"`
	UpdatedAt              time.Time                  `json:"updatedAt"`
	Source                 events.Source              `json:"source"`
	ExternalSystem         commonmodel.ExternalSystem `json:"externalSystem,omitempty"`
	InvoicingStartDate     *time.Time                 `json:"invoicingStartDate,omitempty"`
	Currency               string                     `json:"currency"`
	BillingCycle           string                     `json:"billingCycle"` //Deprecated: Use BillingCycleInMonths instead
	InvoicingEnabled       bool                       `json:"invoicingEnabled"`
	PayOnline              bool                       `json:"payOnline,omitempty"`
	PayAutomatically       bool                       `json:"payAutomatically,omitempty"`
	CanPayWithCard         bool                       `json:"canPayWithCard,omitempty"`
	CanPayWithDirectDebit  bool                       `json:"canPayWithDirectDebit,omitempty"`
	CanPayWithBankTransfer bool                       `json:"canPayWithBankTransfer,omitempty"`
	AutoRenew              bool                       `json:"autoRenew,omitempty"`
	Check                  bool                       `json:"check,omitempty"`
	DueDays                int64                      `json:"dueDays,omitempty"`
	Country                string                     `json:"country"`
	LengthInMonths         int64                      `json:"lengthInMonths"`
	BillingCycleInMonths   int64                      `json:"billingCycleInMonths"`
	Approved               bool                       `json:"approved"`
}

func NewContractCreateEvent(aggregate eventstore.Aggregate, dataFields model.ContractDataFields, source events.Source, externalSystem commonmodel.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContractCreateEvent{
		Tenant:                 aggregate.GetTenant(),
		OrganizationId:         dataFields.OrganizationId,
		Name:                   dataFields.Name,
		ContractUrl:            dataFields.ContractUrl,
		CreatedByUserId:        dataFields.CreatedByUserId,
		ServiceStartedAt:       utils.ToDatePtr(dataFields.ServiceStartedAt),
		SignedAt:               utils.ToDatePtr(dataFields.SignedAt),
		Currency:               dataFields.Currency,
		BillingCycleInMonths:   dataFields.BillingCycleInMonths,
		InvoicingStartDate:     utils.ToDatePtr(dataFields.InvoicingStartDate),
		InvoicingEnabled:       dataFields.InvoicingEnabled,
		PayOnline:              dataFields.PayOnline,
		PayAutomatically:       dataFields.PayAutomatically,
		CanPayWithCard:         dataFields.CanPayWithCard,
		CanPayWithDirectDebit:  dataFields.CanPayWithDirectDebit,
		CanPayWithBankTransfer: dataFields.CanPayWithBankTransfer,
		AutoRenew:              dataFields.AutoRenew,
		Check:                  dataFields.Check,
		DueDays:                dataFields.DueDays,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
		Source:                 source,
		Country:                dataFields.Country,
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
