package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

type InvoiceForContractCreateEvent struct {
	Tenant          string             `json:"tenant" validate:"required"`
	ContractId      string             `json:"organizationId" validate:"required"`
	CreatedAt       time.Time          `json:"createdAt"`
	SourceFields    commonmodel.Source `json:"sourceFields"`
	DryRun          bool               `json:"dryRun"`
	Currency        string             `json:"currency"`
	PeriodStartDate time.Time          `json:"periodStartDate"`
	PeriodEndDate   time.Time          `json:"periodEndDate"`
	BillingCycle    string             `json:"billingCycle" validate:"required_if=OffCycle false"`
	Note            string             `json:"note"`
	OffCycle        bool               `json:"offCycle,omitempty"`
	Postpaid        bool               `json:"postpaid"`
}

func NewInvoiceForContractCreateEvent(aggregate eventstore.Aggregate, sourceFields commonmodel.Source, contractId, currency, billingCycle, note string, dryRun, offCycle, postpaid bool, createdAt, periodStartDate, periodEndDate time.Time) (eventstore.Event, error) {
	eventData := InvoiceForContractCreateEvent{
		Tenant:          aggregate.GetTenant(),
		ContractId:      contractId,
		CreatedAt:       createdAt,
		SourceFields:    sourceFields,
		Currency:        currency,
		DryRun:          dryRun,
		PeriodStartDate: periodStartDate,
		PeriodEndDate:   periodEndDate,
		BillingCycle:    billingCycle,
		Note:            note,
		OffCycle:        offCycle,
		Postpaid:        postpaid,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceCreateForContractV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceCreateEvent")
	}

	return event, nil
}
