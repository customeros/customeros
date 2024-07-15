package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type InvoiceForContractCreateEvent struct {
	Tenant               string        `json:"tenant" validate:"required"`
	ContractId           string        `json:"contractId" validate:"required"`
	CreatedAt            time.Time     `json:"createdAt"`
	SourceFields         common.Source `json:"sourceFields"`
	DryRun               bool          `json:"dryRun"`
	Currency             string        `json:"currency"`
	PeriodStartDate      time.Time     `json:"periodStartDate"`
	PeriodEndDate        time.Time     `json:"periodEndDate"`
	BillingCycle         string        `json:"billingCycle"` // Deprecated: Use BillingCycleInMonths instead
	BillingCycleInMonths int64         `json:"billingCycleInMonths" validate:"required_if=OffCycle false"`
	Note                 string        `json:"note"`
	OffCycle             bool          `json:"offCycle,omitempty"`
	Postpaid             bool          `json:"postpaid"`
	Preview              bool          `json:"preview"`
}

func NewInvoiceForContractCreateEvent(aggregate eventstore.Aggregate, sourceFields common.Source, contractId, currency, note string, billingCycleInMonths int64, dryRun, offCycle, postpaid, preview bool, createdAt, periodStartDate, periodEndDate time.Time) (eventstore.Event, error) {
	eventData := InvoiceForContractCreateEvent{
		Tenant:               aggregate.GetTenant(),
		ContractId:           contractId,
		CreatedAt:            createdAt,
		SourceFields:         sourceFields,
		Currency:             currency,
		DryRun:               dryRun,
		PeriodStartDate:      periodStartDate,
		PeriodEndDate:        periodEndDate,
		BillingCycleInMonths: billingCycleInMonths,
		Note:                 note,
		OffCycle:             offCycle,
		Postpaid:             postpaid,
		Preview:              preview,
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
